# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

proto/opiapi.pb.go: opi-api/opiproto.proto
	protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative proto/opiapi.proto

api/storage.pb.go: api/storage.proto
	protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative api/storage.proto
