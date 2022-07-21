# OPI API Conglomerate

This implements the OPI API Conglomerate, which is a pluggable server to
realize the OPI APIs.

## Organization

For now, this sample program is taken from the wonderful golang plugin library
from Hashicorp called [plugin-go](https://github.com/hashicorp/go-plugin).
The example `OPI API Conglomerate` is based on the [gRPC example](https://github.com/hashicorp/go-plugin/tree/master/examples/grpc).

Currently, we show the ability to do the following:

* Load plugins into the main `opi-api` application.
* Use gRPC between `opi-api` and the plugin
* Implement two simple APIs:
  * CreateNetwork()
  * GetNetwork()
* The plugin implements these and saves the Network name and data to a
  local file on CreateNetwork(). On GetNetwork(), it will retrieve the
  data which was stored.

This is meant to show the concept of what is discussed [here](https://github.com/opiproject/opi-api/pull/41).

** TODO

Lots! Seriously, off the top of my head:

* Make this into a daemon
* Provide gRPC northbound from the daemon
* Flesh out the actual APIs in the protobuf file
* Flesh out the sample plugin
