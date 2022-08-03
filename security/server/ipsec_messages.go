// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

//import (
//	"fmt"

//	"github.com/strongswan/govici/vici"
//	pb "github.com/opiproject/opi-api/security/proto"
//)

// type childSA struct {
// 	LocalTrafficSelectors []string `vici:"local_ts"`
// 	Updown                string   `vici:"updown"`
// 	ESPProposals          []string `vici:"esp_proposals"`
// }

// type localOpts struct {
// 	Auth  string   `vici:"auth"`
// 	Certs []string `vici:"certs"`
// 	ID    string   `vici:"id"`
// }

// type remoteOpts struct {
// 	Auth string `vici:"auth"`
// }

// type connection struct {
// 	Name string // This field will NOT be marshaled!

// 	LocalAddrs []string            `vici:"local_addrs"`
// 	Local      *localOpts          `vici:"local"`
// 	Remote     *remoteOpts         `vici:"remote"`
// 	Children   map[string]*childSA `vici:"children"`
// 	Version    int                 `vici:"version"`
// 	Proposals  []string            `vici:"proposals"`
// }
