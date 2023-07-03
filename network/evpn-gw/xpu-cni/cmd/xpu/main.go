/*-
 * ============LICENSE_START=======================================================
 *  Copyright (C) 2023 Nordix Foundation.
 * ================================================================================
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 * ============LICENSE_END=========================================================
 */

package main

import (
	"errors"
	"fmt"
	"runtime"

	"xpu-cni/pkg/config"
	"xpu-cni/pkg/sriov"
	"xpu-cni/pkg/utils"
	"xpu-cni/pkg/xpu"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ipam"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

type envArgs struct {
	types.CommonArgs
	MAC types.UnmarshallableString `json:"mac,omitempty"`
}

func init() {
	// this ensures that main runs only on main thread (thread group leader).
	// since namespace ops (unshare, setns) are done for a single thread, we
	// must ensure that the goroutine does not jump from OS thread to thread
	runtime.LockOSThread()
}

func getEnvArgs(envArgsString string) (*envArgs, error) {
	if envArgsString != "" {
		e := envArgs{}
		err := types.LoadArgs(envArgsString, &e)
		if err != nil {
			return nil, err
		}
		return &e, nil
	}
	return nil, nil
}

func cmdAdd(args *skel.CmdArgs) error {
	var macAddr string
	netConf, err := config.LoadConf(args.StdinData)
	if err != nil {
		return fmt.Errorf("xpu-cni failed to load netconf: %v", err)
	}

	envArgs, err := getEnvArgs(args.Args)
	if err != nil {
		return fmt.Errorf("xpu-cni failed to parse args: %v", err)
	}

	if envArgs != nil {
		MAC := string(envArgs.MAC)
		if MAC != "" {
			netConf.MAC = MAC
		}
	}

	// RuntimeConfig takes preference than envArgs.
	// This maintains compatibility of using envArgs
	// for MAC config.
	if netConf.RuntimeConfig.Mac != "" {
		netConf.MAC = netConf.RuntimeConfig.Mac
	}

	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		return fmt.Errorf("failed to open netns %q: %v", netns, err)
	}
	defer netns.Close()

	sm := sriov.NewSriovManager()
	err = sm.FillOriginalVfInfo(netConf)
	if err != nil {
		return fmt.Errorf("failed to get original vf information: %v", err)
	}
	defer func() {
		// TODO: Check here if the cleanup is correct for all the cases.
		if err != nil {
			// Check here if BridgePortName exists in netconf. If yes that means that the bridgeport is there
			// and we need to delete it. If error occurs from delete operation then just let the flow going.
			if netConf.BridgePortName != "" {
				// Delete BridgePort
				_ = xpu.DeleteBridgePort(netConf)
			}
			err := netns.Do(func(_ ns.NetNS) error {
				_, err := netlink.LinkByName(args.IfName)
				return err
			})
			if err == nil {
				_ = sm.ReleaseVF(netConf, netns, args.Netns)
			}
			// Reset the VF if failure occurs before the netconf is cached
			_ = sm.ResetVFConfig(netConf)
		}
	}()
	// Here I have removed the ":" from the err := sm.ApplyVFConfig
	if err = sm.ApplyVFConfig(netConf); err != nil {
		return fmt.Errorf("xpu-cni failed to configure VF %q", err)
	}

	result := &current.Result{}
	result.Interfaces = []*current.Interface{{
		Name:    args.IfName,
		Sandbox: netns.Path(),
	}}

	if !netConf.DPDKMode {
		macAddr, err = sm.SetupVF(netConf, args.IfName, netns)
		if err != nil {
			return fmt.Errorf("failed to set up pod interface %q from the device %q: %v", args.IfName, netConf.Master, err)
		}

		result.Interfaces[0].Mac = macAddr
	} else {
		// Handle here the dpdk case for the VF
		if netConf.PciToMacPath == "" {
			return errors.New("PciToMacPath cannot be empty when the device is not NetDev")
		}
		macAddr, err = utils.RetrieveMacFromPci(netConf.DeviceID, netConf.PciToMacPath)
		if err != nil {
			return fmt.Errorf("error in retrieving the Mac from PciToMac %s file for pci %s : %q", netConf.PciToMacPath, netConf.DeviceID, err)
		}

		result.Interfaces[0].Mac = macAddr
	}

	// Create Bridge Port representor for xPU VF

	//Create BridgePort
	err = xpu.CreateBridgePort(netConf, result.Interfaces[0].Mac)
	if err != nil {
		return err
	}

	// run the IPAM plugin
	if netConf.IPAM.Type != "" {
		var r types.Result
		r, err = ipam.ExecAdd(netConf.IPAM.Type, args.StdinData)
		if err != nil {
			return fmt.Errorf("failed to set up IPAM plugin type %q from the device %q: %v", netConf.IPAM.Type, netConf.Master, err)
		}

		defer func() {
			if err != nil {
				_ = ipam.ExecDel(netConf.IPAM.Type, args.StdinData)
			}
		}()

		// Convert the IPAM result into the current Result type
		var newResult *current.Result
		newResult, err = current.NewResultFromResult(r)
		if err != nil {
			return err
		}

		if len(newResult.IPs) == 0 {
			err = errors.New("IPAM plugin returned missing IP config")
			return err
		}

		newResult.Interfaces = result.Interfaces

		for _, ipc := range newResult.IPs {
			// All addresses apply to the container interface (move from host)
			ipc.Interface = current.Int(0)
		}

		if !netConf.DPDKMode {
			err = netns.Do(func(_ ns.NetNS) error {
				err := ipam.ConfigureIface(args.IfName, newResult)
				if err != nil {
					return err
				}

				/* After IPAM configuration is done, the following needs to handle the case of an IP address being reused by a different pods.
				 * This is achieved by sending Gratuitous ARPs and/or Unsolicited Neighbor Advertisements unconditionally.
				 * Although we set arp_notify and ndisc_notify unconditionally on the interface (please see EnableArpAndNdiscNotify()), the kernel
				 * only sends GARPs/Unsolicited NA when the interface goes from down to up, or when the link-layer address changes on the interfaces.
				 * These scenarios are perfectly valid and recommended to be enabled for optimal network performance.
				 * However for our specific case, which the kernel is unaware of, is the reuse of IP addresses across pods where each pod has a different
				 * link-layer address for it's SRIOV interface. The ARP/Neighbor cache residing in neighbors would be invalid if an IP address is reused.
				 * In order to update the cache, the GARP/Unsolicited NA packets should be sent for performance reasons. Otherwise, the neighbors
				 * may be sending packets with the incorrect link-layer address. Eventually, most network stacks would send ARPs and/or Neighbor
				 * Solicitation packets when the connection is unreachable. This would correct the invalid cache; however this may take a significant
				 * amount of time to complete.
				 *
				 * The error is ignored here because enabling this feature is only a performance enhancement.
				 */
				_ = utils.AnnounceIPs(args.IfName, newResult.IPs)
				return nil
			})
			if err != nil {
				return err
			}
		}
		result = newResult
	}

	// Cache NetConf for CmdDel
	if err = utils.SaveNetConf(args.ContainerID, config.DefaultCNIDir, args.IfName, netConf); err != nil {
		return fmt.Errorf("error saving NetConf %q", err)
	}

	allocator := utils.NewPCIAllocator(config.DefaultCNIDir)
	// Mark the pci address as in used
	if err = allocator.SaveAllocatedPCI(netConf.DeviceID, args.Netns); err != nil {
		return fmt.Errorf("error saving the pci allocation for vf pci address %s: %v", netConf.DeviceID, err)
	}

	return types.PrintResult(result, netConf.CNIVersion)
	// I HAVE CHECKED THE CODE UP TO HERE. NO SIGNIFICANT CHANGES ARE NEEDED. Next steps to Check the cmdDel
	// I would like to check again the VLAN, add the bridgeport section
	// and update the caching netconf with vport_id (write also defer function for bridge port)
}

func cmdDel(args *skel.CmdArgs) error {
	netConf, cRefPath, err := config.LoadConfFromCache(args)
	if err != nil {
		return fmt.Errorf("cmdDel() error loading the cached Netconf: %q", err)
	}

	sm := sriov.NewSriovManager()

	defer func() {
		if err == nil && cRefPath != "" {
			_ = utils.CleanCachedNetConf(cRefPath)
		}
	}()

	defer func() {
		// Try to cleanup as much as possible.
		// The functions are idempotent so should run correctly
		// when you get an error or not from the cmdDel.
		_ = xpu.DeleteBridgePort(netConf)
		_ = sm.ResetVFConfig(netConf)
		_ = sm.ResetVF(netConf)

	}()

	if netConf.IPAM.Type != "" {
		err = ipam.ExecDel(netConf.IPAM.Type, args.StdinData)
		if err != nil {
			return err
		}
	}

	// Delete Bridge Port.
	err = xpu.DeleteBridgePort(netConf)
	if err != nil {
		return fmt.Errorf("cmdDel() error deleting Bridge Port: %q", err)
	}

	/* ResetVFConfig resets a VF administratively. We must run ResetVFConfig
	   before ReleaseVF because some drivers will error out if we try to
	   reset netdev VF with trust off. So, reset VF MAC address via PF first.
	*/
	err = sm.ResetVFConfig(netConf)
	if err != nil {
		return fmt.Errorf("cmdDel() error resetting VF administratively: %q", err)
	}

	if !netConf.DPDKMode {
		if args.Netns == "" {
			// Reset netdev VF to its original state
			err = sm.ResetVF(netConf)
			if err != nil {
				return fmt.Errorf("cmdDel() error Resetting VF to original state: %q", err)
			}
		} else {
			var netns ns.NetNS
			netns, err = ns.GetNS(args.Netns)
			if err != nil {
				_, ok := err.(ns.NSPathNotExistErr)
				if ok {
					// Reset netdev VF to its original state if NS Path not exist
					err = sm.ResetVF(netConf)
					if err != nil {
						return fmt.Errorf("cmdDel() error Resetting VF to original state: %q", err)
					}
				} else {
					return fmt.Errorf("cmdDel() failed to open netns %s: %q", netns, err)
				}
			} else {

				defer netns.Close()

				// Release VF form Pods namespace and rename it to the original name
				err = sm.ReleaseVF(netConf, netns, args.Netns)
				if err != nil {
					return fmt.Errorf("cmdDel() error releasing VF: %q", err)
				}
			}
		}
	}

	// Mark the pci address as released
	allocator := utils.NewPCIAllocator(config.DefaultCNIDir)
	if err = allocator.DeleteAllocatedPCI(netConf.DeviceID); err != nil {
		return fmt.Errorf("cmdDel() error cleaning the pci allocation for vf pci address %s: %v", netConf.DeviceID, err)
	}

	return nil

}

func cmdCheck(_ *skel.CmdArgs) error {
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, "")
}
