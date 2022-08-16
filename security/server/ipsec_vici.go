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
	Local        map[string]interface{} `vici:"local"`
	Remote       map[string]interface{} `vici:"remote"`
	Children     map[string]childSA     `vici:"children"`
	Version      int                    `vici:"version"`
	Proposals    []string               `vici:"proposals"`
	Sendcertreq  string                 `vici:"send_certreq"`
}

type unload_connection struct {
	//Cname string // This field will NOT be marshaled!
	Name string `vici:"name"`
}

/*
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
*/

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
		conn.Local = map[string]interface{} {
			"auth": strings.ToLower(c.GetLocalAuth().String()),
			"id": c.GetLocalAuth().GetId(),
		}
	}
	if c.GetRemoteAuth() != nil {
		conn.Remote = map[string]interface{} {
			"auth": strings.ToLower(c.GetRemoteAuth().String()),
			"id": c.GetRemoteAuth().GetId(),
		}
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
				local_ts = append(local_ts, c.Children[i].LocalTs.Ts[k].String())
			}
		}
		if c.Children[i].RemoteTs != nil {
			for k := 0; k < len(c.Children[i].RemoteTs.Ts); k++ {
				remote_ts = append(remote_ts, c.Children[i].RemoteTs.Ts[k].String())
			}
		}
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
