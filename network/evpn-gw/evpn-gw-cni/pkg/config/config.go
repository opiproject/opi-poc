// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2023 Network Plumping Working Group
// Copyright (C) 2023 Nordix Foundation.

package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	evpngwtypes "evpn-gw-cni/pkg/types"
	"evpn-gw-cni/pkg/utils"
	"github.com/containernetworking/cni/pkg/skel"

	"github.com/imdario/mergo"
)

const UnsupportedFields = "MinTxRate, MaxTxRate, SpoofChk, Trust, LinkState"

var (
	// DefaultCNIDir used for caching NetConf
	DefaultCNIDir = "/var/lib/cni/evpn-gw"
)

// LoadConf parses and validates stdin netconf and returns NetConf object
func LoadConf(bytes []byte) (*evpngwtypes.NetConf, error) {
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
		// Here the SRIOV CNI was returning the NetConf object instead of nil.
		// I do not see the point of that so I have changed that to return nil
		return nil, err
	}

	if isAllocated {
		// Here the SRIOV CNI was returning the NetConf object instead of nil.
		// I do not see the point of that so I have changed that to return nil
		return nil, fmt.Errorf("pci address %s is already allocated", n.DeviceID)
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

	if n.LogicalBridge != "" && len(n.LogicalBridges) > 0 {
		return nil, fmt.Errorf("LoadConf(): can not define both LogicalBridge and LogicalBridges. Have to pick one of those")
	}

	// validate that link state is one of supported values
	/*if n.LinkState != "" && n.LinkState != "auto" && n.LinkState != "enable" && n.LinkState != "disable" {
		return nil, fmt.Errorf("LoadConf(): invalid link_state value: %s", n.LinkState)
	}*/

	// This block of code should be removed when the parameters are supported by XPU.
	// Also uncomment the related code blocks for the supported fields here and in the sriov pkg
	if n.MinTxRate != nil || n.MaxTxRate != nil || n.SpoofChk != "" || n.Trust != "" || n.LinkState != "" {
		fmt.Printf("LoadConf(): The %s configuration fields are not supported currently", UnsupportedFields)
		fmt.Printf("LoadConf(): The %s configuration fields will be ignored", UnsupportedFields)
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
func LoadConfFromCache(args *skel.CmdArgs) (*evpngwtypes.NetConf, string, error) {
	netConf := &evpngwtypes.NetConf{}

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

func loadNetConf(bytes []byte) (*evpngwtypes.NetConf, error) {
	netconf := &evpngwtypes.NetConf{}
	if err := json.Unmarshal(bytes, netconf); err != nil {
		return nil, fmt.Errorf("failed to load netconf: %v", err)
	}

	return netconf, nil
}

func loadFlatNetConf(configPath string) (*evpngwtypes.NetConf, error) {
	confFiles := getEvpnGwConfFiles()
	if configPath != "" {
		confFiles = append([]string{configPath}, confFiles...)
	}

	// loop through the path and parse the JSON config
	flatNetConf := &evpngwtypes.NetConf{}
	for _, confFile := range confFiles {
		confExists, err := utils.PathExists(confFile)
		if err != nil {
			return nil, fmt.Errorf("error checking evpn-gw-cni config file: error: %v", err)
		}
		if confExists {
			jsonFile, err := os.Open(confFile)
			if err != nil {
				return nil, fmt.Errorf("open evpn-gw-cni config file %s error: %v", confFile, err)
			}
			defer jsonFile.Close()
			jsonBytes, err := io.ReadAll(jsonFile)
			if err != nil {
				return nil, fmt.Errorf("load evpn-gw-cni config file %s: error: %v", confFile, err)
			}
			if err := json.Unmarshal(jsonBytes, flatNetConf); err != nil {
				return nil, fmt.Errorf("parse evpn-gw-cni config file %s: error: %v", confFile, err)
			}
			break
		}
	}

	return flatNetConf, nil
}

func mergeConf(netconf, flatNetConf *evpngwtypes.NetConf) (*evpngwtypes.NetConf, error) {
	if err := mergo.Merge(netconf, flatNetConf); err != nil {
		return nil, fmt.Errorf("merge with evpn-gw-cni config file: error: %v", err)
	}
	return netconf, nil
}

func getEvpnGwConfFiles() []string {
	return []string{"/etc/kubernetes/cni/net.d/evpn-gw.d/evpn-gw.conf", "/etc/cni/net.d/evpn-gw.d/evpn-gw.conf"}
}
