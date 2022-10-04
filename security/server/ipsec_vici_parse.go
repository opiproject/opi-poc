// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"encoding/base64"
	"log"
	"strings"

	pb "github.com/opiproject/opi-api/security/v1/gen/go"
)

func parse_child_list_sas(childsa list_child_sa, name string) (*pb.ListChildSa, error) {
	log.Printf("Found key %v", childsa)

	child := &pb.ListChildSa{}

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
	list_ret := &pb.ListIkeSa{}

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

func parse_auth(conn_auth list_auth, name string) (*pb.ListConnAuth, error) {
	log.Printf("Found key %v", conn_auth)

	auth := &pb.ListConnAuth{}

	if conn_auth.Class != "" {
		auth.Class = conn_auth.Class
	}
	if conn_auth.EapType != "" {
		auth.Eaptype = conn_auth.EapType
	}
	if conn_auth.EapVendor != "" {
		auth.Eapvendor = conn_auth.EapVendor
	}
	if conn_auth.Xauth != "" {
		auth.Xauth = conn_auth.Xauth
	}
	if conn_auth.Revocation != "" {
		auth.Revocation = conn_auth.Revocation
	}
	if conn_auth.Id != "" {
		auth.Id = conn_auth.Id
	}
	if conn_auth.CaId != "" {
		auth.CaId = conn_auth.CaId
	}
	if conn_auth.AaaId != "" {
		auth.AaaId = conn_auth.AaaId
	}
	if conn_auth.EapId != "" {
		auth.EapId = conn_auth.EapId
	}
	if conn_auth.XauthId != "" {
		auth.XauthId = conn_auth.XauthId
	}
	if conn_auth.Groups != nil {
		for k := 0; k < len(conn_auth.Groups); k++ {
			auth.Group.Group = append(auth.Group.Group, conn_auth.Groups[k])
		}
	}
	if conn_auth.CertPolicy != nil {
		for k := 0; k < len(conn_auth.CertPolicy); k++ {
			auth.CertPolicy.CertPolicy = append(auth.CertPolicy.CertPolicy, conn_auth.CertPolicy[k])
		}
	}
	if conn_auth.Certs != nil {
		for k := 0; k < len(conn_auth.Certs); k++ {
			auth.Certs.Cert = append(auth.Certs.Cert, conn_auth.Certs[k])
		}
	}
	if conn_auth.CaCerts != nil {
		for k := 0; k < len(conn_auth.CaCerts); k++ {
			auth.Cacerts.Cacert = append(auth.Cacerts.Cacert, conn_auth.CaCerts[k])
		}
	}

	return auth, nil
}

func parse_connection_child(list_child list_child, name string) (*pb.ListChild, error) {
	log.Printf("Found key %v", list_child)

	child := &pb.ListChild{}

	child.Name = name
	if list_child.Mode != "" {
		child.Mode = list_child.Mode
	}
	if list_child.Label != "" {
		child.Label = list_child.Label
	}
	if list_child.RekeyTime != 0 {
		child.RekeyTime = list_child.RekeyTime
	}
	if list_child.RekeyBytes != 0 {
		child.RekeyBytes = list_child.RekeyBytes
	}
	if list_child.RekeyPackets != 0 {
		child.RekeyPackets = list_child.RekeyPackets
	}
	if list_child.DpdAction != "" {
		child.DpdAction = list_child.DpdAction
	}
	if list_child.CloseAction != "" {
		child.CloseAction = list_child.CloseAction
	}
	if list_child.RemoteTs != nil {
		for k := 0; k < len(list_child.RemoteTs); k++ {
			s := strings.Split(list_child.RemoteTs[k], ":")
			ts := &pb.TrafficSelectors_TrafficSelector{}
			if len(s) >= 1 && s[0] != "" {
				ts.Cidr = s[0]
			}
			if len(s) >= 2 && s[1] != "" {
				ts.Proto = s[1]
			}
			if len(s) >= 3 && s[2] != "" {
				ts.Port = s[2]
			}

			child.RemoteTs = &pb.TrafficSelectors{}
			child.RemoteTs.Ts = []*pb.TrafficSelectors_TrafficSelector{}
			child.RemoteTs.Ts = append(child.RemoteTs.Ts, ts)
		}
	}
	if list_child.LocalTs != nil {
		for k := 0; k < len(list_child.LocalTs); k++ {
			s := strings.Split(list_child.LocalTs[k], ":")
			ts := &pb.TrafficSelectors_TrafficSelector{}
			if len(s) >= 1 && s[0] != "" {
				ts.Cidr = s[0]
			}
			if len(s) >= 2 && s[1] != "" {
				ts.Proto = s[1]
			}
			if len(s) >= 3 && s[2] != "" {
				ts.Port = s[2]
			}

			child.LocalTs = &pb.TrafficSelectors{}
			child.LocalTs.Ts = []*pb.TrafficSelectors_TrafficSelector{}
			child.LocalTs.Ts = append(child.LocalTs.Ts, ts)
		}
	}
	if list_child.Interface != "" {
		child.Interface = list_child.Interface
	}
	if list_child.Priority != "" {
		child.Priority = list_child.Priority
	}

	return child, nil
}

func parse_connection(conn *list_ike, km string) (*pb.ListConnResp, error) {
	list_ret := &pb.ListConnResp{}

	list_ret.Name = km
	for i := 0; i < len(conn.LocalAddrs); i++ {
		addr := &pb.Addrs{Addr: conn.LocalAddrs[i]}
		list_ret.LocalAddrs = append(list_ret.LocalAddrs, addr)
	}
	for i := 0; i < len(conn.RemoteAddrs); i++ {
		addr := &pb.Addrs{Addr: conn.RemoteAddrs[i]}
		list_ret.RemoteAddrs = append(list_ret.RemoteAddrs, addr)
	}
	if conn.Version != "" {
		list_ret.Version = conn.Version
	}
	if conn.ReauthTime != 0 {
		list_ret.ReauthTime = conn.ReauthTime
	}
	if conn.RekeyTime != 0 {
		list_ret.RekeyTime = conn.RekeyTime
	}
	if conn.Unique != "" {
		list_ret.Unique = conn.Unique
	}
	if conn.DpdDelay != 0 {
		list_ret.DpdDelay = conn.DpdDelay
	}
	if conn.DpdTimeout != 0 {
		list_ret.DpdTimeout = conn.DpdTimeout
	}
	if conn.Ppk != "" {
		list_ret.Ppk = conn.Ppk
	}
	if conn.PpkRequired != "" {
		list_ret.PpkRequired = conn.PpkRequired
	}
	if conn.Local != nil {
		log.Printf("Looking at Local Auth%v", conn.Local)
		for k, mess := range conn.Local {
			local, err := parse_auth(mess, k)
			if err != nil {
				log.Printf("Error parsing local auth")
				return nil, nil
			}
			list_ret.LocalAuth = append(list_ret.LocalAuth, local)
		}
	}
	if conn.Remote != nil {
		log.Printf("Looking at Remote Auth%v", conn.Remote)
		for k, mess := range conn.Remote {
			remote, err := parse_auth(mess, k)
			if err != nil {
				log.Printf("Error parsing remote auth")
				return nil, nil
			}
			list_ret.RemoteAuth = append(list_ret.RemoteAuth, remote)
		}
	}
	if conn.Children != nil {
		log.Printf("Looking at Children%v", conn.Children)
		for k, mess := range conn.Children {
			child, err := parse_connection_child(mess, k)
			if err != nil {
				log.Printf("Error parsing child")
				return nil, nil
			}
			list_ret.Children = append(list_ret.Children, child)
		}
	}

	return list_ret, nil
}

func parse_certificate(cert *list_cert) (*pb.ListCert, error) {
	list_ret := &pb.ListCert{}

	if cert.Type != "" {
		s1 := strings.ToUpper(cert.Type)

		if strings.Contains(s1, "X509_AC") {
			list_ret.Type = pb.CertificateType_CERT_X509
		} else if strings.Contains(s1, "X509_CRL") {
			list_ret.Type = pb.CertificateType_CERT_X509_CRL
		} else if strings.Contains(s1, "X509") {
			list_ret.Type = pb.CertificateType_CERT_X509
		} else if strings.Contains(s1, "OCSP_RESPONSE") {
			list_ret.Type = pb.CertificateType_CERT_OCSP_RESPONSE
		} else if strings.Contains(s1, "PUBKEY") {
			list_ret.Type = pb.CertificateType_CERT_PUBKEY
		}
	}
	if cert.Flag != "" {
		s1 := strings.ToUpper(cert.Flag)

		if strings.Contains(s1, "OCSP") {
			list_ret.Flag = pb.X509CertificateFlag_X509_CERT_FLAG_OCSP
		} else if strings.Contains(s1, "AA") {
			list_ret.Flag = pb.X509CertificateFlag_X509_CERT_FLAG_AA
		} else if strings.Contains(s1, "CA") {
			list_ret.Flag = pb.X509CertificateFlag_X509_CERT_FLAG_CA
		} else if strings.Contains(s1, "NONE") {
			list_ret.Flag = pb.X509CertificateFlag_X509_CERT_FLAG_NONE
		}
	}
	if cert.HasPrivKey != "" {
		list_ret.Hasprivkey = cert.HasPrivKey
	}
	if cert.Data != "" {
		encodedText := base64.StdEncoding.EncodeToString([]byte(cert.Data))
		list_ret.Data = encodedText
	}
	if cert.Subject != "" {
		list_ret.Subject = cert.Subject
	}
	if cert.NotBefore != "" {
		list_ret.Notbefore = cert.NotBefore
	}
	if cert.NotAfter != "" {
		list_ret.Notafter = cert.NotAfter
	}

	return list_ret, nil
}
