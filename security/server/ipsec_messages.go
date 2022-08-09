// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

//	"github.com/strongswan/govici/vici"
	pb "github.com/opiproject/opi-api/security/proto"
)

//type childSA struct {
//	LocalTrafficSelectors []string `vici:"local_ts"`
//	Updown                string   `vici:"updown"`
//	ESPProposals          []string `vici:"esp_proposals"`
//}

//type localOpts struct {
//	Auth  string   `vici:"auth"`
//	Certs []string `vici:"certs"`
//	ID    string   `vici:"id"`
//}

//type remoteOpts struct {
//	Auth string `vici:"auth"`
//}

//type connection struct {
//	Name string // This field will NOT be marshaled!
//
//	LocalAddrs []string            `vici:"local_addrs"`
//	Local      *localOpts          `vici:"local"`
//	Remote     *remoteOpts         `vici:"remote"`
//	Children   map[string]*childSA `vici:"children"`
//	Version    int                 `vici:"version"`
//	Proposals  []string            `vici:"proposals"`
//}

func removeSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

func load_connection(command string, conn *pb.IPsecCreateRequest) error {
	var sb strings.Builder
	var tstr string

	// What has to happen here is:
	// 1. Write charon config file for connection
	// 2. Send vici command to reload-configuration
	// 3. Send vici command to connect with the new connection

	// Create file name, removing spaces
	no_space_sb := removeSpace(conn.Tunnel.Tunnels[0].GetName())
	saved_name := no_space_sb
	no_space_sb = fmt.Sprintf("/etc/swanctl/%s.conf", no_space_sb)
	log.Printf("Looking at file %s", no_space_sb)

	// Header
	tstr = fmt.Sprintf("%s {\n", saved_name)
	sb.WriteString(tstr)

	tstr = fmt.Sprintf("    remote_addrs = %s\n", conn.Tunnel.Tunnels[0].GetRemoteIp())
	sb.WriteString(tstr)

	tstr = "    vips = 0.0.0.0\n"
	sb.WriteString(tstr)

	tstr = "    local {\n"
	sb.WriteString(tstr)

	tstr = "        auth = pubkey\n"
	sb.WriteString(tstr)

	tstr = "        certs = clientCert.pem\n"
	sb.WriteString(tstr)

	tstr = "        id = client.strongswan.org\n"
	sb.WriteString(tstr)

	tstr = "    }\n"
	sb.WriteString(tstr)

	tstr = "    remote {\n"
	sb.WriteString(tstr)

	tstr = "        auth = pubkey\n"
	sb.WriteString(tstr)

	tstr = "        id = server.strongswan.org\n"
	sb.WriteString(tstr)

	tstr = "    }\n"
	sb.WriteString(tstr)

	tstr = "    children {\n"
	sb.WriteString(tstr)

	tstr = "        net {\n"
	sb.WriteString(tstr)

	tstr = "            remote_ts = 10.1.0.0/16\n"
	sb.WriteString(tstr)

	tstr = "            esp_proposals = "
	sb.WriteString(tstr)

	for i := 0; i < len(conn.Tunnel.Tunnels); i++ {
		sb.WriteString(strings.ToLower(conn.Tunnel.Tunnels[i].GetCryptoAlg().String()))
		tstr = "-"
		sb.WriteString(tstr)
		sb.WriteString(strings.ToLower(conn.Tunnel.Tunnels[i].GetIntegAlg().String()))
		if (i+1) < len(conn.Tunnel.Tunnels) {
			tstr = "-"
			sb.WriteString(tstr)
		}
	}

	tstr = "\n"
	sb.WriteString(tstr)

	tstr = "            dpd_action = trap\n"
	sb.WriteString(tstr)

	tstr = "        }\n"
	sb.WriteString(tstr)

	tstr = "        host {\n"
	sb.WriteString(tstr)

	tstr = "            esp_proposals = "
	sb.WriteString(tstr)

	for i := 0; i < len(conn.Tunnel.Tunnels); i++ {
		sb.WriteString(strings.ToLower(conn.Tunnel.Tunnels[i].GetCryptoAlg().String()))
		tstr = "-"
		sb.WriteString(tstr)
		sb.WriteString(strings.ToLower(conn.Tunnel.Tunnels[i].GetIntegAlg().String()))
		if (i+1) < len(conn.Tunnel.Tunnels) {
			tstr = "-"
			sb.WriteString(tstr)
		}
	}

	tstr = "\n"
	sb.WriteString(tstr)

	tstr = "            dpd_action = trap\n"
	sb.WriteString(tstr)

	tstr = "        }\n"
	sb.WriteString(tstr)

	tstr = "    }\n"
	sb.WriteString(tstr)

	tstr = "    version = 2\n"
	sb.WriteString(tstr)

	tstr = "    proposals = "
	sb.WriteString(tstr)

	for i := 0; i < len(conn.Sa.Sas); i++ {
		sb.WriteString(strings.ReplaceAll(strings.ToLower(conn.Sa.Sas[i].GetCryptoAlg().String()), "_", ""))
		tstr = "-"
		sb.WriteString(tstr)
		sb.WriteString(strings.ReplaceAll(strings.ToLower(conn.Sa.Sas[i].GetIntegAlg().String()), "_", ""))
		if (i+1) < len(conn.Tunnel.Tunnels) {
			tstr = "-"
			sb.WriteString(tstr)
		}
	}

	tstr = "\n"
	sb.WriteString(tstr)

	tstr = "    dpd_delay = 60s\n"
	sb.WriteString(tstr)


	tstr = "}\n"
	sb.WriteString(tstr)

	log.Printf("----- Config file -----")
	log.Printf("\n%v\n", sb.String())
	log.Printf("-----------------------")

	// Check if file exists
	if _, err := os.Stat(no_space_sb); errors.Is(err, os.ErrExist) {
		return errors.New("Config file already exists")
	}

	// Create file
	f, err := os.Create(no_space_sb)
	if err != nil {
		return errors.New("Failed creating file")
	}
	defer f.Close()

	// Write config to file
	n3, err := f.WriteString(sb.String())
	if err != nil {
		return errors.New("Failed writing configuration file")
	}
	log.Printf("Wrote %d bytes to file", n3)

	err = f.Sync()
	if err != nil {
		log.Printf("Failed syncing file, soft eror")
	}

	return nil
}

func delete_connection(command string, conn *pb.IPsecCreateRequest) error {
	no_space_sb := removeSpace(conn.Tunnel.Tunnels[0].GetName())
	no_space_sb = fmt.Sprintf("/etc/swanctl/%s.conf", no_space_sb)
	log.Printf("Looking at file %s", no_space_sb)

	err := os.Remove(no_space_sb)
	if err != nil {
		return errors.New("Failed removing file")
	}

	return nil
}
