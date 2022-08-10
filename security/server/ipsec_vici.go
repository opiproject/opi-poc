// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	//"errors"
	//"fmt"
	"log"
	//"os"
	"strings"
	//"unicode"

	"github.com/strongswan/govici/vici"
	pb "github.com/opiproject/opi-api/security/proto"
)

type childSA struct {
	RemoteTrafficSelectors []string  `vici:"remote_ts"`
	LocalTrafficSelectors []string   `vici:"local_ts"`
	Updown                string     `vici:"updown"`
	ESPProposals          []string   `vici:"esp_proposals"`
}

type connection struct {
	Name string // This field will NOT be marshaled!

	LocalAddrs   []string               `vici:"local_addrs"`
	Local        map[string]interface{} `vici:"local"`
	Remote       map[string]interface{} `vici:"remote"`
	Children     map[string]childSA     `vici:"children"`
	Version      int                    `vici:"version"`
	Proposals    []string               `vici:"proposals"`
	Sendcertreq  string                 `vici:"send_certreq"`
}

func loadConn(connreq *pb.IPsecCreateRequest) error {
	var ike_prop strings.Builder
	var esp_prop strings.Builder
	var tstr string

	// Setup IKE proposals
	for i := 0; i < len(connreq.Sa.Sas); i++ {
		ike_prop.WriteString(strings.ReplaceAll(strings.ToLower(connreq.Sa.Sas[i].GetCryptoAlg().String()), "_", ""))
		tstr = "-"
		ike_prop.WriteString(tstr)
		ike_prop.WriteString(strings.ReplaceAll(strings.ToLower(connreq.Sa.Sas[i].GetIntegAlg().String()), "_", ""))
		tstr = "-"
		ike_prop.WriteString(tstr)
		ike_prop.WriteString(strings.ToLower(connreq.Sa.Sas[i].GetDhgroups().String()))
		if (i+1) < len(connreq.Sa.Sas) {
			tstr = "-"
			ike_prop.WriteString(tstr)
		}
	}

	// Setup child proposals
	for i := 0; i < len(connreq.Tunnel.Tunnels); i++ {
		esp_prop.WriteString(strings.ToLower(connreq.Tunnel.Tunnels[i].GetCryptoAlg().String()))
		tstr = "-"
		esp_prop.WriteString(tstr)
		esp_prop.WriteString(strings.ToLower(connreq.Tunnel.Tunnels[i].GetIntegAlg().String()))
		tstr = "-"
		esp_prop.WriteString(tstr)
		esp_prop.WriteString(strings.ToLower(connreq.Tunnel.Tunnels[i].GetDhgroups().String()))
		if (i+1) < len(connreq.Tunnel.Tunnels) {
			tstr = "-"
			esp_prop.WriteString(tstr)
		}
        }

	// Build the connection object
	conn := &connection {
		Name: connreq.GetName(),
		Local: map[string]interface{} {
			"auth": "psk",
			"id": "hacker@strongswan.org",
		},
		Remote: map[string]interface{} {
			"auth": "psk",
			"id": "server.strongswan.org",
		},
		Children: map[string]childSA{
			"net": childSA {
				RemoteTrafficSelectors: []string {"10.1.0.0/16"},
				ESPProposals: []string {esp_prop.String()},
			},
			"host": childSA {
				ESPProposals: []string {esp_prop.String()},
			},
		},
		Version: 2,
		Proposals: []string {ike_prop.String()},
		Sendcertreq: "no",
	}

	log.Printf("Built connection object: %v", conn)
	for k, v := range conn.Children {
		log.Printf("key[%s] value[%s]\n", k, v)
	}
	log.Printf("Local: %v", conn.Local)
	log.Printf("Remote: %v", conn.Remote)

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

	m := vici.NewMessage()
	if err := m.Set(conn.Name, c); err != nil {
		log.Printf("Failed setting command")
		return err
	}

	_, err = s.CommandRequest("load-conn", m)

	log.Printf("command error return [%v]", err)

	return err
}

/*
func initiate_connection(ike, child string) error {
	s, err := vici.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	m := vici.NewMessage()

	if err := m.Set("child", child); err != nil {
		return err
	}

	if err := m.Set("ike", ike); err != nil {
		return err
	}

	ms, err := s.StreamedCommandRequest("initiate", "control-log", m)
	if err != nil {
		return err
	}

	for _, msg := range ms.Messages() {
		if err := msg.Err(); err != nil {
			return err
		}
	}

	return nil
}
*/
