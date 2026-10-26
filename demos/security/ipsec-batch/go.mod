module github.com/opiproject/opi-poc/demos/security/ipsec-batch

go 1.26.0

require (
	github.com/opiproject/opi-api/network/opinetcommon/v1alpha1/gen/go v0.0.0-20261026183418-48a1ffe43abe
	github.com/opiproject/opi-api/v1/gen/go/security v0.0.0-20261026183418-48a1ffe43abe
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.32.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20261026165753-2ffab32bce7e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261019180522-ce908d29299b // indirect
)

// Drop these replace lines once opiproject/opi-api#496 merges.
replace (
	github.com/opiproject/opi-api/network/opinetcommon/v1alpha1/gen/go => github.com/mgheorghe/opi-api/network/opinetcommon/v1alpha1/gen/go v0.0.0-20261026183418-48a1ffe43abe
	github.com/opiproject/opi-api/v1/gen/go/security => github.com/mgheorghe/opi-api/v1/gen/go/security v0.0.0-20261026183418-48a1ffe43abe
)
