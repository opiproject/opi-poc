// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
        "log"

        pb "github.com/opiproject/opi-api/security/proto"
)

func parse_child_list_sas(childsa list_child_sa, name string) (*pb.ListChildSa, error) {
	log.Printf("Found key %v", childsa )

	child := &pb.ListChildSa {}

	child.Name = name
	if childsa.Protocol != "" {
		child.Protocol = childsa.Protocol
	}
	if childsa.Encap != "" {
		child.Encap = childsa.Encap
	}
	if childsa.SpiIn != "" {
		child.SpiIn = childsa.SpiIn
	}
	if childsa.SpiOut != "" {
		child.SpiOut = childsa.SpiOut
	}
	if childsa.CpiIn != "" {
		child.CpiIn = childsa.CpiIn
	}
	if childsa.CpiOut != "" {
		child.CpiOut = childsa.CpiOut
	}
	if childsa.MarkIn != "" {
		child.MarkIn = childsa.MarkIn
	}
	if childsa.MarkMaskIn != "" {
		child.MarkMaskIn = childsa.MarkMaskIn
	}
	if childsa.MarkOut != "" {
		child.MarkOut = childsa.MarkOut
	}
	if childsa.MarkMaskOut != "" {
		child.MarkMaskOut = childsa.MarkMaskOut
	}
	if childsa.IfIdIn != "" {
		child.IfIdIn = childsa.IfIdIn
	}
	if childsa.IfIdOut != "" {
		child.IfIdOut = childsa.IfIdOut
	}
	if childsa.EncrAlg != "" {
		child.EncrAlg = childsa.EncrAlg
	}
	if childsa.EncKeysize != "" {
		child.EncrKeysize = childsa.EncKeysize
	}
	if childsa.IntegAlg != "" {
		child.IntegAlg = childsa.IntegAlg
	}
	if childsa.IntegKeysize != "" {
		child.IntegKeysize = childsa.IntegKeysize
	}
	if childsa.DhGroup != "" {
		child.DhGroup = childsa.DhGroup
	}
	if childsa.Esn != "" {
		child.Esn = childsa.Esn
	}

	return child, nil
}

func parse_ike_list_sas(ikesa *list_ike_sa, km string) (*pb.ListIkeSa, error) {
	list_ret := &pb.ListIkeSa {}

	list_ret.Name = km
	if ikesa.UniqueId != "" {
		list_ret.Uniqueid = ikesa.UniqueId
	}
	if ikesa.Version != "" {
		list_ret.Version = ikesa.Version
	}
	if ikesa.State != "" {
		list_ret.Ikestate = pb.IkeSaState(pb.IkeSaState_value[ikesa.State])
	}
	if ikesa.LocalHost != "" {
		list_ret.LocalHost = ikesa.LocalHost
	}
	if ikesa.LocalPort != "" {
		list_ret.LocalPort = ikesa.LocalPort
	}
	if ikesa.LocalId != "" {
		list_ret.LocalId = ikesa.LocalId
	}
	if ikesa.RemoteHost != "" {
		list_ret.RemoteHost = ikesa.RemoteHost
	}
	if ikesa.RemotePort != "" {
		list_ret.RemotePort = ikesa.RemotePort
	}
	if ikesa.RemoteId != "" {
		list_ret.RemoteId = ikesa.RemoteId
	}
	if ikesa.RemoteXauthId != "" {
		list_ret.RemoteXauthId = ikesa.RemoteXauthId
	}
	if ikesa.RemoteEapId != "" {
		list_ret.RemoteEapId = ikesa.RemoteEapId
	}
	if ikesa.Initiator != "" {
		list_ret.Initiator = ikesa.Initiator
	}
	if ikesa.InitiatorSpi != "" {
		list_ret.InitiatorSpi = ikesa.InitiatorSpi
	}
	if ikesa.ResponderSpi != "" {
		list_ret.ResponderSpi = ikesa.ResponderSpi
	}
	if ikesa.NatLocal != "" {
		list_ret.NatLocal = ikesa.NatLocal
	}
	if ikesa.NatRemote != "" {
		list_ret.NatRemote = ikesa.NatRemote
	}
	if ikesa.NatFake != "" {
		list_ret.NatFake = ikesa.NatFake
	}
	if ikesa.NatAny != "" {
		list_ret.NatAny = ikesa.NatAny
	}
	if ikesa.IfIdIn != "" {
		list_ret.IfIdIn = ikesa.IfIdIn
	}
	if ikesa.IfIdOut != "" {
		list_ret.IfIdOut = ikesa.IfIdOut
	}
	if ikesa.EncrAlg != "" {
		list_ret.EncrAlg = ikesa.EncrAlg
	}
	if ikesa.EncrKeysize != "" {
		list_ret.EncrKeysize = ikesa.EncrKeysize
	}
	if ikesa.IntegAlg != "" {
		list_ret.IntegAlg = ikesa.IntegAlg
	}
	if ikesa.IntegKeysize != "" {
		list_ret.IntegKeysize = ikesa.IntegKeysize
	}
	if ikesa.PrfAlg != "" {
		list_ret.PrfAlg = ikesa.PrfAlg
	}
	if ikesa.DhGroup != "" {
		list_ret.DhGroup = ikesa.DhGroup
	}
	if ikesa.Ppk != "" {
		list_ret.Ppk = ikesa.Ppk
	}
	if ikesa.Established != "" {
		list_ret.Established = ikesa.Established
	}
	if ikesa.RekeyTime != "" {
		list_ret.RekeyTime = ikesa.RekeyTime
	}
	if ikesa.ReauthTime != "" {
		list_ret.ReauthTime = ikesa.ReauthTime
	}
	if ikesa.LocalVips != nil {
		list_ret.LocalVips = ikesa.LocalVips
	}
	if ikesa.RemoteVips != nil {
		list_ret.RemoteVips = ikesa.RemoteVips
	}
	if ikesa.TasksQueued != nil {
		list_ret.TasksQueued = ikesa.TasksQueued
	}
	if ikesa.TasksActive != nil {
		list_ret.TasksActive = ikesa.TasksActive
	}
	if ikesa.TasksPassive != nil {
		list_ret.TasksPassive = ikesa.TasksPassive
	}
	if ikesa.ChildSas != nil {
		log.Printf("Looking at ChildSas %v", ikesa.ChildSas)
		for k, mess := range ikesa.ChildSas {
			childsa, err := parse_child_list_sas(mess, k)
			if err != nil {
				log.Printf("Error parsing CHILD_SA")
				return nil, nil
			}
			list_ret.Childsas = append(list_ret.Childsas, childsa)
		}
	}

	return list_ret, nil
}
