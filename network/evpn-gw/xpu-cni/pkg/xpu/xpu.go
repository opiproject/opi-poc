// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2023 Nordix Foundation.

package xpu

import (
	"context"
	"errors"
	"fmt"
	"time"

	xputypes "xpu-cni/pkg/types"

	xpuMgr "github.com/opiproject/opi-api/network/evpn-gw/v1alpha1/gen/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// initConnection initializes a connection to XPU Infra Manager
func initConnection(conf *xputypes.NetConf) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if conf.XpuInfraMgrConn == "" {
		return nil, errors.New("XpuInfraMgrConn netconf field cannot be empty")
	}

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(conf.XpuInfraMgrConn, opts...)
	if err != nil {
		return nil, fmt.Errorf("fail to dial: %q", err)
	}

	return conn, nil
}

// closeConnection closes a grpc connection
func closeConnection(conn *grpc.ClientConn) {
	conn.Close()
}

// getClient gets a BridgePort grpc client
func getClient(conn *grpc.ClientConn) xpuMgr.BridgePortServiceClient {
	client := xpuMgr.NewBridgePortServiceClient(conn)
	return client
}

// produceCreateBridgePortRequest produces a CreateBridgePortRequest object
func produceCreateBridgePortRequest(conf *xputypes.NetConf, mac string) *xpuMgr.CreateBridgePortRequest {
	var typeOfPort xpuMgr.BridgePortType
	var logicalBridges []string

	if conf.LogicalBridge != "" {
		typeOfPort = xpuMgr.BridgePortType_ACCESS
		logicalBridges = []string{conf.LogicalBridge}
	} else {
		typeOfPort = xpuMgr.BridgePortType_TRUNK
		if len(conf.LogicalBridges) > 0 {
			logicalBridges = conf.LogicalBridges
		}
	}

	bridgePortSpec := &xpuMgr.BridgePortSpec{
		MacAddress:     []byte(mac),
		Ptype:          typeOfPort,
		LogicalBridges: logicalBridges,
	}

	bridgePort := &xpuMgr.BridgePort{
		Spec: bridgePortSpec,
	}

	createBridgePortRequest := &xpuMgr.CreateBridgePortRequest{
		BridgePort: bridgePort,
	}

	return createBridgePortRequest
}

// produceDeleteBridgePortRequest produces a DeleteBridgePortRequest object
func produceDeleteBridgePortRequest(conf *xputypes.NetConf) *xpuMgr.DeleteBridgePortRequest {
	deleteBridgePortRequest := &xpuMgr.DeleteBridgePortRequest{
		Name: conf.BridgePortName,
	}

	return deleteBridgePortRequest
}

// CreateBridgePort creates a bridge port
func CreateBridgePort(conf *xputypes.NetConf, mac string) error {
	//Init Connection
	conn, err := initConnection(conf)
	if err != nil {
		return fmt.Errorf("CreateBridgePort: Error occurred while init connection:  %q", err)
	}

	defer closeConnection(conn)

	// Get a Client
	client := getClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// produce the createBridgePortRequest object
	createBridgePortRequest := produceCreateBridgePortRequest(conf, mac)

	// grpc call to create the bridge port
	bridgePort, err := client.CreateBridgePort(ctx, createBridgePortRequest)
	if err != nil {
		return fmt.Errorf("CreateBridgePort: Error occurred while creating Bridge Port: %q", err)
	}

	// storing the name of the created bridge port to the netconf object for caching purposes
	conf.BridgePortName = bridgePort.GetName()

	if bridgePort.GetStatus().GetOperStatus() != xpuMgr.BPOperStatus_BP_OPER_STATUS_UP {
		return errors.New("CreateBridgePort: The status of created BridgePort is not UP")
	}

	return nil
}

func DeleteBridgePort(conf *xputypes.NetConf) error {
	// Check if the BridgePortName exists in the NetConf object.
	// If it doesn't exist then we simply return nil as there is no point to continue
	// as we need the BridgePortName for the BridgePort delete process to execute.
	// The reason that we do not return error is because we want to give the chance
	// to the delete process to continue with the rest of the tasks
	// (e.g. ReleaseVFs, ResetVFs, etc...) so there is no leftovers in the system.
	if conf.BridgePortName == "" {
		return nil
	}

	//Init Connection
	conn, err := initConnection(conf)
	if err != nil {
		return fmt.Errorf("DeleteBridgePort: Error occurred while init connection:  %q", err)
	}

	defer closeConnection(conn)

	// Get a Client
	client := getClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// produce the deleteBridgePortRequest object
	deleteBridgePortRequest := produceDeleteBridgePortRequest(conf)

	// If error is BridgePort not found then return nil in order to serve idempotence.
	_, err = client.DeleteBridgePort(ctx, deleteBridgePortRequest)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}
		return fmt.Errorf("DeleteBridgePort: Error occurred while Deleting Bridge Port %s : %q", conf.BridgePortName, err)
	}

	return nil
}
