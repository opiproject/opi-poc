# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

all: proto/opiapi.pb.go api/storage/storage.pb.go

proto/opiapi.pb.go: proto/opiapi.proto
	protoc proto/opiapi.proto \
            --go_out=. \
            --go_opt=paths=source_relative \
            --go-grpc_out=. \
            --go-grpc_opt=paths=source_relative \
            --proto_path=.

api/storage/storage.pb.go: api/storage/storage.proto
	protoc api/storage/storage.proto \
            --go_out=. \
            --go_opt=paths=source_relative \
            --go-grpc_out=. \
            --go-grpc_opt=paths=source_relative \
            --proto_path=.

