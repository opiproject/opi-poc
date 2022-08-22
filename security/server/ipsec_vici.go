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
	AgProposals            []string   `vici:"ag_proposals"`
	RekeyTime              string     `vici:"rekey_time"`
	LifeTime               string     `vici:"life_time"`
	RandTime               string     `vici:"rand_time"`
	Inactivity             uint32     `vici:"inactivity"`
	MarkIn                 uint32     `vici:"mark_in"`
	MarkInSa               string     `vici:"mark_in_sa"`
	MarkOut                uint32     `vici:"mark_out"`
	SetMarkIn              uint32     `vici:"set_mark_in"`
	SetMarkOut             uint32     `vici:"set_mark_out"`
	HwOffload              string     `vici:"hw_offload"`
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
	Vips         []string               `vici:"vips"`
	LocalPort    uint32                 `vici:"local_port"`
	RemotePort   uint32                 `vici:"remote_port"`
	Dscp         uint64                 `vici:"dscp"`
	Encap        string                 `vici:"encap"`
	Mobike       string                 `vici:"mobike"`
	DpdDelay     uint32                 `vici:"dpd_delay"`
	DpdTimeout   uint32                 `vici:"dpd_timeout"`
	ReauthTime   uint32                 `vici:"reauth_time"`
	RekeyTime    string                 `vici:"rekey_time"`
	Pools        []string               `vici:"pools"`
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

func buildProposal(prop *pb.Proposals) (string, error) {
	var crypto strings.Builder
	var integ strings.Builder
	var prf strings.Builder
	var dh strings.Builder
	var compiled_proposal strings.Builder
	var tstr string

	for k := 0; k < len(prop.CryptoAlg); k++ {
		crypto.WriteString(strings.ToLower(prop.CryptoAlg[k].String()))
		if (k+1) < len(prop.CryptoAlg) {
			tstr = "-"
			crypto.WriteString(tstr)
		}
	}
	for k := 0; k < len(prop.IntegAlg); k++ {
		integ.WriteString(strings.ToLower(prop.IntegAlg[k].String()))
		if (k+1) < len(prop.IntegAlg) {
			tstr = "-"
			integ.WriteString(tstr)
		}
	}
	for k := 0; k < len(prop.Prf); k++ {
		prf.WriteString(strings.ToLower(prop.Prf[k].String()))
		if (k+1) < len(prop.Prf) {
			tstr = "-"
			prf.WriteString(tstr)
		}
	}
	for k := 0; k < len(prop.Dhgroups); k++ {
		dh.WriteString(strings.ToLower(prop.Dhgroups[k].String()))
		if (k+1) < len(prop.Dhgroups) {
			tstr = "-"
			dh.WriteString(tstr)
		}
	}

	if crypto.String() != "" {
		compiled_proposal.WriteString(crypto.String())
		if integ.String() != "" || prf.String() != "" || dh.String() != "" {
			tstr = "-"
			compiled_proposal.WriteString(tstr)
		}
	}
	if integ.String() != "" {
		compiled_proposal.WriteString(integ.String())
		if prf.String() != "" || dh.String() != "" {
			tstr = "-"
			compiled_proposal.WriteString(tstr)
		}
	}
	if prf.String() != "" {
		compiled_proposal.WriteString(prf.String())
		if dh.String() != "" {
			tstr = "-"
			compiled_proposal.WriteString(tstr)
		}
	}
	if dh.String() != "" {
		compiled_proposal.WriteString(dh.String())
	}

	return compiled_proposal.String(), nil
}

func ipsecVersion() (*pb.IPsecVersionResp, error) {
	s, err := vici.NewSession()
	if err != nil {
		log.Printf("Failed creating vici session")
		return nil, err
	}
	defer s.Close()

	m, err := s.CommandRequest("version", nil)
	if err != nil {
		log.Printf("Failed getting version")
		return nil, err
	}

	daemon  := m.Get("daemon").(string)
	version := m.Get("version").(string)
	sysname := m.Get("sysname").(string)
	release := m.Get("release").(string)
	machine := m.Get("machine").(string)

	// Assemble return value
	verresp := &pb.IPsecVersionResp {
		Daemon:  daemon,
		Version: version,
		Sysname: sysname,
		Release: release,
		Machine: machine,
	}

	return verresp, err
}

func ipsecStats() (*pb.IPsecStatsResp, error) {
	s, err := vici.NewSession()
	if err != nil {
		log.Printf("Failed creating vici session")
		return nil, err
	}
	defer s.Close()

	m, err := s.CommandRequest("stats", nil)
	if err != nil {
		log.Printf("Failed getting stats")
		return nil, err
	}

	var stat_string strings.Builder

	value := m.Get("uptime")
	field, ok := value.(*vici.Message)
	if !ok {
		log.Printf("Embedded map key was not marshaled as a sub-message")
	} else {
		running := field.Get("running").(string)
		since   := field.Get("since").(string)
		stat_string.WriteString("Running time: " + running + "\n")
		stat_string.WriteString("Absolute startup time: " + since + "\n")
	}

	value = m.Get("workers")
	field, ok = value.(*vici.Message)
	if ! ok {
		log.Printf("Cannot find workers in map")
	} else {
		total := field.Get("total").(string)
		idle  := field.Get("idle").(string)
		stat_string.WriteString("Total # of worker threads: " + total + "\n")
		stat_string.WriteString("Worker threads currently idle: " + idle + "\n")

		value = field.Get("active")
		subfield, subok := value.(*vici.Message)
		if !subok {
			log.Printf("Cannot find active in map")
		} else {
			critical := subfield.Get("critical").(string)
			high     := subfield.Get("high").(string)
			medium   := subfield.Get("medium").(string)
			low      := subfield.Get("low").(string)
			stat_string.WriteString("Threads processing critical priority jobs: " + critical + "\n")
			stat_string.WriteString("Threads processing high priority jobs: " + high + "\n")
			stat_string.WriteString("Threads processing medium priority jobs: " + medium + "\n")
			stat_string.WriteString("Threads processing low priority jobs: " + low + "\n")
		}
	}

	value = m.Get("queues")
	field, ok = value.(*vici.Message)
	if !ok {
		log.Printf("Cannot find queues in map")
	} else {
		critical := field.Get("critical").(string)
		high     := field.Get("high").(string)
		medium   := field.Get("medium").(string)
		low      := field.Get("low").(string)
		stat_string.WriteString("Jobs queued with critical priority: " + critical + "\n")
		stat_string.WriteString("Jobs queued with high priority: " + high + "\n")
		stat_string.WriteString("Jobs queued with medium priority: " + medium + "\n")
		stat_string.WriteString("Jobs queued with low priority: " + low + "\n")
	}

	scheduled := m.Get("scheduled").(string)
	stat_string.WriteString("# of jobs scheduled for timed execution: " + scheduled + "\n")

	value = m.Get("ikesas")
	field, ok = value.(*vici.Message)
	if !ok {
		log.Printf("Cannot find ikesas")
	} else {
		total     := field.Get("total").(string)
		half_open := field.Get("half-open").(string)
		stat_string.WriteString("Total number of IKE_SAs active: " + total + "\n")
		stat_string.WriteString("Number of IKE_SAs in half-open state: " + half_open + "\n")
	}

	stat_string.WriteString("Plugins: ")
	plugins := m.Get("plugins").([]string)
	for c := 0; c < len(plugins); c++ {
		stat_string.WriteString(plugins[c] + " ")
	}
	stat_string.WriteString("\n")

	value = m.Get("mem")
	field, ok = value.(*vici.Message)
	if !ok {
		log.Printf("Cannot find mem")
	} else {
		total  := field.Get("total").(string)
		allocs := field.Get("allocs").(string)
		stat_string.WriteString("Total heap memory usage in bytes: " + total + "\n")
		stat_string.WriteString("Total heap allocation in blocks: " + allocs + "\n")

		// NOTE: Skipping heap-name (Windows only) for now since OPI
		//       does not run on Windows.
	}

	value = m.Get("mallinfo")
	field, ok = value.(*vici.Message)
	if !ok {
		log.Printf("Cannot find mallinfo")
	} else {
		sbrk := field.Get("sbrk").(string)
		mmap := field.Get("mmap").(string)
		used := field.Get("used").(string)
		free := field.Get("free").(string)
		stat_string.WriteString("Non-mmap'd space available: " + sbrk + "\n")
		stat_string.WriteString("Mmap'd space available: " + mmap + "\n")
		stat_string.WriteString("Total number of bytes used: " + used + "\n")
		stat_string.WriteString("Available but unsued bytes: " + free + "\n")
	}

	statsresp := &pb.IPsecStatsResp {
		Status: stat_string.String(),
	}

	return statsresp, nil
}

func loadConn(connreq *pb.IPsecLoadConnReq) error {
	// Declare the connection variable, as we have to conditionally load it
	var conn = &connection {
		LocalPort: 500,
		RemotePort: 500,
		RekeyTime: "4h",
	}
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
	if c.GetLocalPort() != 0 {
		conn.LocalPort = c.GetLocalPort()
	}
	if c.GetRemotePort() != 0 {
		conn.RemotePort = c.GetRemotePort()
	}
	log.Printf("Local Port [%d] Remote Port [%d]", conn.LocalPort, conn.RemotePort)

	if c.GetDscp() != 0 {
		conn.Dscp = c.GetDscp()
	}
	if c.GetEncap() != "" {
		conn.Encap = c.GetEncap()
	}
	if c.GetMobike() != "" {
		conn.Mobike = c.GetMobike()
	}
	if c.GetDpdDelay() != 0 {
		conn.DpdDelay = c.GetDpdDelay()
	}
	if c.GetDpdTimeout() != 0 {
		conn.DpdTimeout = c.GetDpdTimeout()
	}
	if c.GetReauthTime() != 0 {
		conn.ReauthTime = c.GetReauthTime()
	}
	if c.GetRekeyTime() != 0 {
		var s = strconv.FormatUint(uint64(c.GetRekeyTime()), 10)
		conn.RekeyTime = s + "s"
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
	if c.Pools != nil {
		for i := 0; i < len(c.Pools.Pool); i++ {
			conn.Pools = append(conn.Pools, c.Pools.Pool[i])
		}
	}

	if c.Proposals != nil {
		ike_proposal, _ := buildProposal(c.Proposals)
		conn.Proposals = []string { ike_proposal }
		log.Printf("IKE proposal: %v", conn.Proposals)
	}

	for i := 0; i < len(c.Children); i++ {
		var local_ts []string
		var remote_ts []string

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
				RekeyTime: "1h",
				LifeTime: "66m",
				Inactivity: c.Children[i].GetInactivity(),
				MarkIn: c.Children[i].GetMarkIn(),
				MarkInSa: "no",
				MarkOut: c.Children[i].GetMarkOut(),
				SetMarkIn: c.Children[i].GetSetMarkIn(),
				SetMarkOut: c.Children[i].GetSetMarkOut(),
				HwOffload: "no",
		}
		if c.Children[i].EspProposals != nil {
			proposal, _ := buildProposal(c.Children[i].EspProposals)
			csa.ESPProposals = []string { proposal }
		}
		if c.Children[i].AgProposals != nil {
			ag_proposal, _ := buildProposal(c.Children[i].AgProposals)
			csa.AgProposals = []string { ag_proposal }
		}

		if c.Children[i].RekeyTime != 0 {
			var s = strconv.FormatUint(uint64(c.Children[i].GetRekeyTime()), 10)
			csa.RekeyTime = s + "s"
		}
		if c.Children[i].LifeTime != 0 {
			var s = strconv.FormatUint(uint64(c.Children[i].GetLifeTime()), 10)
			csa.LifeTime = s + "s"
		}
		if c.Children[i].RandTime != 0 {
			var s = strconv.FormatUint(uint64(c.Children[i].GetRandTime()), 10)
			csa.RandTime = s + "s"
		}
		if c.Children[i].GetMarkInSa() != "" {
			csa.MarkInSa = c.Children[i].GetMarkInSa()
		}
		if c.Children[i].GetHwOffload() != "" {
			csa.HwOffload = c.Children[i].GetHwOffload()
		}

		conn.Children = make(map[string]childSA)
		conn.Children[c.Children[i].GetName()] = csa

		log.Printf("Dumping child object: %v", conn.Children[c.Children[i].GetName()])
		log.Printf("Dumping Proposals: %v", csa.ESPProposals)
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
