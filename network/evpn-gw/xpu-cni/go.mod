module xpu-cni

go 1.20

require (
	github.com/containernetworking/cni v1.1.2
	github.com/containernetworking/plugins v1.2.0
	github.com/imdario/mergo v0.3.15
	github.com/k8snetworkplumbingwg/sriovnet v1.2.0
	github.com/opiproject/opi-api v0.0.0-20231016162146-d81cc5ee60d4
	github.com/vishvananda/netlink v1.2.1-beta.2
	golang.org/x/net v0.8.0
	google.golang.org/grpc v1.55.0
)

require (
	github.com/coreos/go-iptables v0.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/safchain/ethtool v0.2.0 // indirect
	github.com/spf13/afero v1.9.4 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

//replace github.com/opiproject/opi-api => github.com/Nordix/opi-api v0.0.0-20230525125718-f99d0eaa5df5
