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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containernetworking/cni/pkg/skel"
	//TODO: Change this.
	xputypes "xpu-cni/pkg/types"
	"xpu-cni/pkg/utils"

	"github.com/imdario/mergo"
)

var (
	// DefaultCNIDir used for caching NetConf
	DefaultCNIDir = "/var/lib/cni/xpu"
)

// LoadConf parses and validates stdin netconf and returns NetConf object
func LoadConf(bytes []byte) (*xputypes.NetConf, error) {
	n, err := loadNetConf(bytes)
	if err != nil {
		return nil, fmt.Errorf("LoadConf(): failed to load netconf: %v", err)
	}
	flatNetConf, err := loadFlatNetConf(n.ConfigurationPath)
	if err != nil {
		return nil, fmt.Errorf("LoadConf(): failed to load flat netconf: %v", err)
	}
	n, err = mergeConf(n, flatNetConf)
	if err != nil {
		return nil, fmt.Errorf("LoadConf(): failed to merge netconf and flat netconf: %v", err)
	}

	// DeviceID takes precedence; if we are given a VF pciaddr then work from there
	if n.DeviceID != "" {
		// Get rest of the VF information
		pfName, vfID, err := getVfInfo(n.DeviceID)
		if err != nil {
			return nil, fmt.Errorf("LoadConf(): failed to get VF information: %q", err)
		}
		n.VFID = vfID
		n.Master = pfName
	} else {
		return nil, fmt.Errorf("LoadConf(): VF pci addr is required")
	}

	allocator := utils.NewPCIAllocator(DefaultCNIDir)
	// Check if the device is already allocated.
	// This is to prevent issues where kubelet request to delete a pod and in the same time a new pod using the same
	// vf is started. we can have an issue where the cmdDel of the old pod is called AFTER the cmdAdd of the new one
	// This will block the new pod creation until the cmdDel is done.
	isAllocated, err := allocator.IsAllocated(n.DeviceID)
	if err != nil {
		return n, err
	}

	if isAllocated {
		return n, fmt.Errorf("pci address %s is already allocated", n.DeviceID)
	}

	// Assuming VF is netdev interface; Get interface name(s)
	hostIFNames, err := utils.GetVFLinkNames(n.DeviceID)
	if err != nil || hostIFNames == "" {
		// VF interface not found; check if VF has dpdk driver
		hasDpdkDriver, err := utils.HasDpdkDriver(n.DeviceID)
		if err != nil {
			return nil, fmt.Errorf("LoadConf(): failed to detect if VF %s has dpdk driver %q", n.DeviceID, err)
		}
		n.DPDKMode = hasDpdkDriver
	}

	if hostIFNames != "" {
		n.OrigVfState.HostIFName = hostIFNames
	}

	if hostIFNames == "" && !n.DPDKMode {
		return nil, fmt.Errorf("LoadConf(): the VF %s does not have a interface name or a dpdk driver", n.DeviceID)
	}

	/*if n.Vlan != nil && len(n.Trunk) > 0 {
		return nil, fmt.Errorf("LoadConf(): can not define both vlan id and trunk id. Have to pick one of those")
	}*/
	if n.LogicalBridge != "" && len(n.LogicalBridges) > 0 {
		return nil, fmt.Errorf("LoadConf(): can not define both LogicalBridge and LogicalBridges. Have to pick one of those")
	}

	/*if n.Vlan != nil {
		// validate vlan id range
		if *n.Vlan < 0 || *n.Vlan > 4094 {
			return nil, fmt.Errorf("LoadConf(): vlan id %d invalid: value must be in the range 0-4094", *n.Vlan)
		}
		n.PortType = "access"
	} else {
		n.PortType = "trunk"
		if len(n.Trunk) > 0 {
			n.TrunkIDs, err = utils.SplitVlanIds(n.Trunk)
			if err != nil {
				return nil, fmt.Errorf("LoadConf():Error %q occurred while parsing the Trunk list %d", err, n.Trunk)
			}
		}
	}*/

	// The below has been moved to the xpu.go
	/*if n.LogicalBridge != "" {
		// Access case
		n.PortType = "access"
	} else {
		// Trunk case
		n.PortType = "trunk"

	}*/

	// validate that link state is one of supported values
	if n.LinkState != "" && n.LinkState != "auto" && n.LinkState != "enable" && n.LinkState != "disable" {
		return nil, fmt.Errorf("LoadConf(): invalid link_state value: %s", n.LinkState)
	}

	return n, nil
}

func getVfInfo(vfPci string) (string, int, error) {
	var vfID int

	pf, err := utils.GetPfName(vfPci)
	if err != nil {
		return "", vfID, err
	}

	vfID, err = utils.GetVfid(vfPci, pf)
	if err != nil {
		return "", vfID, err
	}

	return pf, vfID, nil
}

// LoadConfFromCache retrieves cached NetConf returns it along with a handle for removal
func LoadConfFromCache(args *skel.CmdArgs) (*xputypes.NetConf, string, error) {
	netConf := &xputypes.NetConf{}

	s := []string{args.ContainerID, args.IfName}
	cRef := strings.Join(s, "-")
	cRefPath := filepath.Join(DefaultCNIDir, cRef)

	netConfBytes, err := utils.ReadScratchNetConf(cRefPath)
	if err != nil {
		return nil, "", fmt.Errorf("error reading cached NetConf in %s with name %s", DefaultCNIDir, cRef)
	}

	if err = json.Unmarshal(netConfBytes, netConf); err != nil {
		return nil, "", fmt.Errorf("failed to parse NetConf: %q", err)
	}

	return netConf, cRefPath, nil
}

func loadNetConf(bytes []byte) (*xputypes.NetConf, error) {
	netconf := &xputypes.NetConf{}
	if err := json.Unmarshal(bytes, netconf); err != nil {
		return nil, fmt.Errorf("failed to load netconf: %v", err)
	}

	return netconf, nil
}

func loadFlatNetConf(configPath string) (*xputypes.NetConf, error) {
	confFiles := getXpuConfFiles()
	if configPath != "" {
		confFiles = append([]string{configPath}, confFiles...)
	}

	// loop through the path and parse the JSON config
	flatNetConf := &xputypes.NetConf{}
	for _, confFile := range confFiles {
		confExists, err := utils.PathExists(confFile)
		if err != nil {
			return nil, fmt.Errorf("error checking xpu-cni config file: error: %v", err)
		}
		if confExists {
			jsonFile, err := os.Open(confFile)
			if err != nil {
				return nil, fmt.Errorf("open xpu-cni config file %s error: %v", confFile, err)
			}
			defer jsonFile.Close()
			jsonBytes, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				return nil, fmt.Errorf("load xpu-cni config file %s: error: %v", confFile, err)
			}
			if err := json.Unmarshal(jsonBytes, flatNetConf); err != nil {
				return nil, fmt.Errorf("parse xpu-cni config file %s: error: %v", confFile, err)
			}
			break
		}
	}

	return flatNetConf, nil
}

func mergeConf(netconf, flatNetConf *xputypes.NetConf) (*xputypes.NetConf, error) {
	if err := mergo.Merge(netconf, flatNetConf); err != nil {
		return nil, fmt.Errorf("merge with xpu-cni config file: error: %v", err)
	}
	return netconf, nil
}

func getXpuConfFiles() []string {
	return []string{"/etc/kubernetes/cni/net.d/xpu.d/xpu.conf", "/etc/cni/net.d/xpu.d/xpu.conf"}
}
