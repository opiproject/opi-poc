package main

// Block Device Abstraction Layer

// Generated via https://mholt.github.io/json-to-go/

// bdev_get_bdevs
// bdev_get_iostat
// bdev_malloc_create
// bdev_malloc_delete
// bdev_null_create
// bdev_null_delete
// bdev_aio_create
// bdev_aio_delete
// bdev_nvme_attach_controller
// bdev_nvme_get_controllers
// bdev_nvme_detach_controller
// bdev_nvme_reset_controller
// bdev_nvme_get_transport_statistics
// bdev_nvme_get_controller_health_info
// bdev_iscsi_create
// bdev_iscsi_delete
// vhost_create_blk_controller
// vhost_delete_controller
// vhost_get_controllers

type BdevAioCreateParams struct {
	Name      string `json:"name"`
	Filename  string `json:"filename"`
	BlockSize int    `json:"block_size"`
}

type BdevAioCreateResult string

type BdevAioDeleteParams struct {
	Name string `json:"name"`
}

type BdevAioDeleteResult bool

type BdevMalloCreateParams struct {
	NumBlocks int    `json:"num_blocks"`
	BlockSize int    `json:"block_size"`
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
}

type BdevAMalloCreateResult string

type BdevMallocDeleteParams struct {
	Name string `json:"name"`
}

type BdevMallocDeleteResult bool

type BdevNvmeAttachControllerParams struct {
	Name        string `json:"name"`
	Type        string `json:"trtype"`
	Address     string `json:"traddr"`
	Family      string `json:"adrfam"`
	Port        string `json:"trsvcid"`
	Subsystem   string `json:"subnqn"`
}

type BdevNvmeAttachControllerResult string

type BdevNvmeDettachControllerParams struct {
	Name string `json:"name"`
}

type BdevNvmeDettachControllerResult bool

type BdevNvmeGetControllerParams struct {
	Name string `json:"name"`
}

type BdevNvmeGetControllerResult struct {
	Name   string `json:"name"`
	Ctrlrs []struct {
		State string `json:"state"`
		Trid  struct {
			Trtype  string `json:"trtype"`
			Adrfam  string `json:"adrfam"`
			Traddr  string `json:"traddr"`
			Trsvcid string `json:"trsvcid"`
			Subnqn  string `json:"subnqn"`
		} `json:"trid"`
		Cntlid int `json:"cntlid"`
		Host   struct {
			Nqn   string `json:"nqn"`
			Addr  string `json:"addr"`
			Svcid string `json:"svcid"`
		} `json:"host"`
	} `json:"ctrlrs"`
}

type BdevGetBdevsParams struct {
	Name string `json:"name"`
}

type BdevGetBdevsResult struct {
	Name        string `json:"name"`
	BlockSize   int64  `json:"block_size"`
	NumBlocks   int64  `json:"num_blocks"`
	Uuid        string `json:"uuid"`
}

type BdevGetIostatParams struct {
	Name string `json:"name"`
}

type BdevGetIostatResult struct {
	TickRate int   `json:"tick_rate"`
	Ticks    int64 `json:"ticks"`
	Bdevs    []struct {
		Name              string `json:"name"`
		BytesRead         int    `json:"bytes_read"`
		NumReadOps        int    `json:"num_read_ops"`
		BytesWritten      int    `json:"bytes_written"`
		NumWriteOps       int    `json:"num_write_ops"`
		BytesUnmapped     int    `json:"bytes_unmapped"`
		NumUnmapOps       int    `json:"num_unmap_ops"`
		ReadLatencyTicks  int    `json:"read_latency_ticks"`
		WriteLatencyTicks int    `json:"write_latency_ticks"`
		UnmapLatencyTicks int    `json:"unmap_latency_ticks"`
	} `json:"bdevs"`
}

type VhostCreateBlkConreollerParams struct {
	Ctrlr   string `json:"ctrlr"`
	DevName string `json:"dev_name"`
}

type VhostCreateBlkConreollerResult bool

type VhostDeleteControllerParams struct {
	Ctrlr   string `json:"ctrlr"`
}

type VhostDeleteControllerResult bool

type VhostGetControllersParams struct {
	Name string `json:"name"`
}

type VhostGetControllersResult struct {
	Ctrlr           string `json:"ctrlr"`
	Cpumask         string `json:"cpumask"`
	DelayBaseUs     int    `json:"delay_base_us"`
	IopsThreshold   int    `json:"iops_threshold"`
	Socket          string `json:"socket"`
	BackendSpecific struct {
		Block struct {
			Readonly bool   `json:"readonly"`
			Bdev     string `json:"bdev"`
		} `json:"block"`
	} `json:"backend_specific"`
}