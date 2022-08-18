// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"log"
	"strings"
	"strconv"

	"github.com/strongswan/govici/vici"
	pb "github.com/opiproject/opi-api/security/proto"
)

type localOpts struct {
	Auth    string   `vici:"auth"`
	Certs   []string `vici:"certs"`
	Id      string   `vici:"id"`
	EapId   string   `vici:"eap_id"`
	AaaId   string   `vici:"aaa_id"`
	XauthId string   `vici:"xauth_id"`
	PubKeys []string `vici:"pubkeys"`
}

type remoteOpts struct {
	Auth       string   `vici:"auth"`
	Id         string   `vici:"id"`
	EapId      string   `vici:"eap_id"`
	Groups     []string `vici:"groups"`
	CertPolicy []string `vici:"cert_policy"`
	Certs      []string `vici:"certs"`
	CaCerts    []string `vici:"certs"`
	PubKeys []string `vici:"pubkeys"`
}

type childSA struct {
	RemoteTrafficSelectors []string  `vici:"remote_ts"`
	LocalTrafficSelectors  []string   `vici:"local_ts"`
	Updown                 string     `vici:"updown"`
	ESPProposals           []string   `vici:"esp_proposals"`
	RekeyTime              string     `vici:"rekey_time"`
}

type connection struct {
	Name string // This field will NOT be marshaled!

	LocalAddrs   []string               `vici:"local_addrs"`
	RemoteAddrs  []string               `vici:"remote_addrs"`
	Local        localOpts              `vici:"local"`
	Remote       remoteOpts             `vici:"remote"`
	Children     map[string]childSA     `vici:"children"`
	Version      int                    `vici:"version"`
	Proposals    []string               `vici:"proposals"`
	Sendcertreq  string                 `vici:"send_certreq"`
	Vips         []string               `vici:"vips"`
}

type unload_connection struct {
	//Cname string // This field will NOT be marshaled!
	Name string `vici:"name"`
}

type init_connection struct {
	Child      string `vici:"child"`
	Ike        string `vici:"ike"`
	Timeout    int    `vici:"timeout"`
	InitLimits string `vici:"init-limits"`
	LogLevel   string `vici:"loglevel"`
}

type terminate_connection struct {
	Child      string `vici:"child"`
	Ike        string `vici:"ike"`
	ChildId    string `vici:"child"`
	IkeId      string `vici:"ike"`
	Force      string `vici:"force"`
	Timeout    int    `vici:"timeout"`
	LogLevel   string `vici:"loglevel"`
}

func loadConn(connreq *pb.IPsecLoadConnReq) error {
	// Declare the connection variable, as we have to conditionally load it
	var conn = &connection {}
	c := connreq.GetConnection()

	if c.GetName() != ""{
		conn.Name = c.GetName()
	 }
	if c.GetVersion() != "" {
		ver, _ := strconv.Atoi(c.GetVersion())
		conn.Version = ver
	} else {
		conn.Version = 2
	}
	if c.GetLocalAuth() != nil {
		conn.Local = localOpts {
			Auth: strings.ToLower(c.GetLocalAuth().GetAuth().String()),
			Id: c.GetLocalAuth().GetId(),
		}

		log.Printf("DUMPING conn.Local: %v", conn.Local)
	}
	if c.GetRemoteAuth() != nil {
		conn.Remote = remoteOpts {
			Auth: strings.ToLower(c.GetRemoteAuth().GetAuth().String()),
			Id: c.GetRemoteAuth().GetId(),
		}

		log.Printf("DUMPING conn.Remove: %v", conn.Remote)
	}

	for i := 0; i < len(c.Vips.Vip); i++ {
		conn.Vips = append(conn.Vips, c.Vips.Vip[i])
	}
	for i := 0; i < len(c.LocalAddrs); i++ {
		conn.LocalAddrs = append(conn.LocalAddrs, c.LocalAddrs[i].GetAddr())
	}
	for i := 0; i < len(c.RemoteAddrs); i++ {
		conn.RemoteAddrs = append(conn.RemoteAddrs, c.RemoteAddrs[i].GetAddr())
	}

	for i := 0; i < len(c.Children); i++ {
		var esp_crypto strings.Builder
		var esp_integ strings.Builder
		var esp_prf strings.Builder
		var esp_dh strings.Builder
		var local_ts []string
		var remote_ts []string
		var tstr string

		for k := 0; k < len(c.Children[i].EspProposals.CryptoAlg); k++ {
			esp_crypto.WriteString(strings.ToLower(c.Children[i].EspProposals.CryptoAlg[k].String()))
			if (k+1) < len(c.Children[i].EspProposals.CryptoAlg) {
				tstr = "-"
				esp_crypto.WriteString(tstr)
			}
		}
		for k := 0; k < len(c.Children[i].EspProposals.IntegAlg); k++ {
			esp_integ.WriteString(strings.ToLower(c.Children[i].EspProposals.IntegAlg[k].String()))
			if (k+1) < len(c.Children[i].EspProposals.IntegAlg) {
				tstr = "-"
				esp_integ.WriteString(tstr)
			}
		}
		for k := 0; k < len(c.Children[i].EspProposals.Prf); k++ {
			esp_prf.WriteString(strings.ToLower(c.Children[i].EspProposals.Prf[k].String()))
			if (k+1) < len(c.Children[i].EspProposals.Prf) {
				tstr = "-"
				esp_prf.WriteString(tstr)
			}
		}
		for k := 0; k < len(c.Children[i].EspProposals.Dhgroups); k++ {
			esp_dh.WriteString(strings.ToLower(c.Children[i].EspProposals.Dhgroups[k].String()))
			if (k+1) < len(c.Children[i].EspProposals.Dhgroups) {
				tstr = "-"
				esp_dh.WriteString(tstr)
			}
		}
		var compiled_proposal strings.Builder
		if esp_crypto.String() != "" {
			compiled_proposal.WriteString(esp_crypto.String())
			if esp_integ.String() != "" || esp_prf.String() != "" || esp_dh.String() != "" {
				tstr = "-"
				compiled_proposal.WriteString(tstr)
			}
		}
		if esp_integ.String() != "" {
			compiled_proposal.WriteString(esp_integ.String())
			if esp_prf.String() != "" || esp_dh.String() != "" {
				tstr = "-"
				compiled_proposal.WriteString(tstr)
			}
		}
		if esp_prf.String() != "" {
			compiled_proposal.WriteString(esp_prf.String())
			if esp_dh.String() != "" {
				tstr = "-"
				compiled_proposal.WriteString(tstr)
			}
		}
		if esp_dh.String() != "" {
			compiled_proposal.WriteString(esp_dh.String())
		}

		if c.Children[i].LocalTs != nil {
			for k := 0; k < len(c.Children[i].LocalTs.Ts); k++ {
				local_ts = append(local_ts, c.Children[i].LocalTs.Ts[k].GetCidr())
			}
		}
		if c.Children[i].RemoteTs != nil {
			for k := 0; k < len(c.Children[i].RemoteTs.Ts); k++ {
				remote_ts = append(remote_ts, c.Children[i].RemoteTs.Ts[k].GetCidr())
			}
		}
		log.Printf("Dumping local_ts [%v] remote_ts [%v]", local_ts, remote_ts)
		csa := childSA {
				LocalTrafficSelectors: local_ts,
				RemoteTrafficSelectors: remote_ts,
				ESPProposals: []string { compiled_proposal.String() },
			}
		if c.Children[i].RekeyTime != 0 {
			var s = strconv.FormatUint(uint64(c.Children[i].RekeyTime), 10)
			csa.RekeyTime = s + "s"
		}

		conn.Children = make(map[string]childSA)
		conn.Children[c.Children[i].GetName()] = csa

		log.Printf("Dumping child object: %v", conn.Children[c.Children[i].GetName()])
	}

	log.Printf("Dumping connection object: %v", conn)

	s, err := vici.NewSession()
	if err != nil {
		log.Printf("Failed creating vici session")
		return err
	}
	defer s.Close()

	mm, err := vici.MarshalMessage(conn)
	if err != nil {
		log.Printf("Failed marshalling message")
		return err
	}

	log.Printf("Marshaled connection request: %v", mm)

	m := vici.NewMessage()
	if err := m.Set(conn.Name, mm); err != nil {
		log.Printf("Failed setting command")
		return err
	}

	_, err = s.CommandRequest("load-conn", m)

	log.Printf("command error return [%v]", err)

	return err
}

func unloadConn(connreq *pb.IPsecUnloadConnReq) error {
	// Build the connection object
	conn := &unload_connection {
		Name: connreq.GetName(),
	}

	log.Printf("Dumping connection to unload: %v", conn)

	s, err := vici.NewSession()
	if err != nil {
		log.Printf("Failed creating vici session")
		return err
	}
	defer s.Close()

	c, err := vici.MarshalMessage(conn)
	if err != nil {
		log.Printf("Failed marshalling message")
		return err
	}

	log.Printf("Marshaled vici message: %v", c)

	_, err = s.CommandRequest("unload-conn", c)

	log.Printf("command error return [%v]", err)

	return err
}

func initiateConn(initreq *pb.IPsecInitiateReq) error {
	init_conn := &init_connection {}

	if initreq.GetChild() != "" {
		init_conn.Child = initreq.GetChild()
	}
	if initreq.GetIke() != "" {
		init_conn.Ike = initreq.GetIke()
	}
	if initreq.GetTimeout() != "" {
		timeout, _ := strconv.Atoi(initreq.GetTimeout())
		init_conn.Timeout = timeout
	}
	if initreq.GetLoglevel() != "" {
		init_conn.LogLevel = initreq.GetLoglevel()
	}

	log.Printf("Dumping connection to initiate: %v", init_conn)

	s, err := vici.NewSession()
	if err != nil {
		log.Printf("Failed creating vici session")
		return err
	}
	defer s.Close()

	c, err := vici.MarshalMessage(init_conn)
	if err != nil {
		log.Printf("Failed marshalling message")
		return err
	}

	log.Printf("Marshaled vici message: %v", c)

	_, err = s.CommandRequest("initiate", c)

	log.Printf("command error return [%v]", err)

	return err
}

func terminateConn(termreq *pb.IPsecTerminateReq) (uint32, error) {
	term_conn := &terminate_connection {}

	if termreq.GetChild() != "" {
		term_conn.Child = termreq.GetChild()
	}
	if termreq.GetIke() != "" {
		term_conn.Ike = termreq.GetIke()
	}
	if termreq.GetChildId() != "" {
		term_conn.ChildId = termreq.GetChildId()
	}
	if termreq.GetIkeId() != "" {
		term_conn.IkeId = termreq.GetIkeId()
	}
	if termreq.GetTimeout() != "" {
		timeout, _ := strconv.Atoi(termreq.GetTimeout())
		term_conn.Timeout = timeout
	}
	if termreq.GetForce() != "" {
		term_conn.Force = termreq.GetForce()
	}
	if termreq.GetLoglevel() != "" {
		term_conn.LogLevel = termreq.GetLoglevel()
	}

	log.Printf("Dumping connection to terminate: %v", term_conn)

	s, err := vici.NewSession()
	if err != nil {
		log.Printf("Failed creating vici session")
		return 0, err
	}
	defer s.Close()

	c, err := vici.MarshalMessage(term_conn)
	if err != nil {
		log.Printf("Failed marshalling message")
		return 0, err
	}

	log.Printf("Marshaled vici message: %v", c)

	_, err = s.CommandRequest("terminate", c)

	log.Printf("command error return [%v]", err)

	return 1, err
}
