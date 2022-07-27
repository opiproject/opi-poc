package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"net"
)

func spdkCommunicate(buf []byte) ([]byte, error) {
	// TODO: use rpc_sock variable
	conn, err := net.Dial("unix", *rpc_sock)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(buf)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Fatal(err)
	}

	reply, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(reply))

	return reply, err
}