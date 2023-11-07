// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2023 Nordix Foundation.

package types

import (
	"github.com/containernetworking/cni/pkg/types"
	"github.com/vishvananda/netlink"
)

// VfState represents the state of the VF
type VfState struct {
	HostIFName   string
	SpoofChk     bool   // Not Used
	AdminMAC     string // Not Used
	EffectiveMAC string // Not Used
	MinTxRate    int    //Not Used
	MaxTxRate    int    //Not Used
	LinkState    uint32 //Not Used
}

// FillFromVfInfo - Fill attributes according to the provided netlink.VfInfo struct - Not Used
func (vs *VfState) FillFromVfInfo(info *netlink.VfInfo) {
	vs.AdminMAC = info.Mac.String()
	vs.LinkState = info.LinkState
	vs.MaxTxRate = int(info.MaxTxRate)
	vs.MinTxRate = int(info.MinTxRate)
	vs.SpoofChk = info.Spoofchk
}

// NetConf extends types.NetConf for xpu-cni
type NetConf struct {
	types.NetConf
	OrigVfState VfState // Stores the original VF state as it was prior to any operations done during cmdAdd flow
	DPDKMode    bool
	Master      string
	MAC         string
	LogicalBridge string `json:"logical_bridge,omitempty"`
	LogicalBridges []string `json:"logical_bridges,omitempty"`
	PortType          string
	BridgePortName    string // Stores the "name" of the created Bridge Port
	DeviceID          string `json:"deviceID"` // PCI address of a VF in valid sysfs format
	VFID              int
	MinTxRate         *int   `json:"min_tx_rate"`          // Mbps, 0 = disable rate limiting (XPU Not supported)
	MaxTxRate         *int   `json:"max_tx_rate"`          // Mbps, 0 = disable rate limiting (XPU Not supported)
	SpoofChk          string `json:"spoofchk,omitempty"`   // on|off (XPU Not supported)
	Trust             string `json:"trust,omitempty"`      // on|off (XPU Not supported)
	LinkState         string `json:"link_state,omitempty"` // auto|enable|disable (XPU Not supported)
	XpuInfraMgrConn   string `json:"xpu_infra_mgr_conn"`   // the IP and port where the xpu_infra_manager listens. Format "IP:Port"
	ConfigurationPath string `json:"configuration_path"`   // Configuration path for xpu-cni conf files
	PciToMacPath      string `json:"pci_to_mac_path"`      // The Path where we keep the mapping of PCI addresses to MAC addresses for the xPU VFs
	RuntimeConfig     struct {
		Mac string `json:"mac,omitempty"`
	} `json:"runtimeConfig,omitempty"`
}
