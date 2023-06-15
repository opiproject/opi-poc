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

package types

import (
	"github.com/containernetworking/cni/pkg/types"
	"github.com/vishvananda/netlink"
)

// VfState represents the state of the VF
type VfState struct {
	HostIFName   string
	SpoofChk     bool
	AdminMAC     string
	EffectiveMAC string
	MinTxRate    int
	MaxTxRate    int
	LinkState    uint32
}

// FillFromVfInfo - Fill attributes according to the provided netlink.VfInfo struct
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
	//Vlan              *uint    `json:"vlan"`
	LogicalBridge string `json:"logical_bridge,omitempty"`
	//Trunk             []*Trunk `json:"trunk,omitempty"`
	LogicalBridges []string `json:"logical_bridges,omitempty"`
	//TrunkIDs          []uint
	PortType          string
	BridgePortName    string // Stores the "name" of the created Bridge Port
	DeviceID          string `json:"deviceID"` // PCI address of a VF in valid sysfs format
	VFID              int
	ContIFNames       string // VF names after in the container; used during deletion
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

// Trunk containing selective vlan IDs
/*type Trunk struct {
	MinID *uint `json:"minID,omitempty"`
	MaxID *uint `json:"maxID,omitempty"`
	ID    *uint `json:"id,omitempty"`
}*/
