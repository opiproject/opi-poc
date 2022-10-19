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

// BdevAioCreateParams holds the parameters required to create an AIO Block Device
type BdevAioCreateParams struct {
	Name      string `json:"name"`
	Filename  string `json:"filename"`
	BlockSize int    `json:"block_size"`
}

// BdevAioCreateResult is the result of creating an AIO Block Device
type BdevAioCreateResult string

// BdevAioDeleteParams holds the parameters required to delete an AIO Block Device
type BdevAioDeleteParams struct {
	Name string `json:"name"`
}

// BdevAioDeleteResult is the result of deleting an AIO Block Device
type BdevAioDeleteResult bool

// BdevMalloCreateParams holds the parameters required to create a Malloc Block Device
type BdevMalloCreateParams struct {
	NumBlocks int    `json:"num_blocks"`
	BlockSize int    `json:"block_size"`
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
}

// BdevAMalloCreateResult is the result of creating a Malloc Block Device
type BdevAMalloCreateResult string

// BdevMallocDeleteParams holds the parameters required to delete a Malloc Block Device
type BdevMallocDeleteParams struct {
	Name string `json:"name"`
}

// BdevMallocDeleteResult is the result of deleting a Malloc Block Device
type BdevMallocDeleteResult bool

// BdevNullCreateParams holds the parameters required to create a Null Block Device
// that discards all writes and returns undefined data for reads
type BdevNullCreateParams struct {
	BlockSize int    `json:"block_size"`
	NumBlocks int    `json:"num_blocks"`
	Name      string `json:"name"`
}

// BdevNullCreateResult is the result of creating a Null Block Device
type BdevNullCreateResult string

// BdevNullDeleteParams holds the parameters required to delete a Null Block Device
type BdevNullDeleteParams struct {
	Name string `json:"name"`
}

// BdevNullDeleteResult is the result of deleting a Null Block Device
type BdevNullDeleteResult bool

// BdevNvmeAttachControllerParams is the parameters required to create a block device based on an NVMe device
type BdevNvmeAttachControllerParams struct {
	Name      string `json:"name"`
	Type      string `json:"trtype"`
	Address   string `json:"traddr"`
	Family    string `json:"adrfam"`
	Port      string `json:"trsvcid"`
	Subsystem string `json:"subnqn"`
}

// BdevNvmeAttachControllerResult is the result of creating a block device based on an NVMe device
type BdevNvmeAttachControllerResult string

// BdevNvmeDetachControllerParams is the parameters required to detach a block device based on an NVMe device
type BdevNvmeDetachControllerParams struct {
	Name string `json:"name"`
}

// BdevNvmeDetachControllerResult is the result of detaching a block device based on an NVMe device
type BdevNvmeDetachControllerResult bool

// BdevNvmeGetControllerParams is the parameters required to get a block device based on an NVMe device
type BdevNvmeGetControllerParams struct {
	Name string `json:"name"`
}

// BdevNvmeGetControllerResult is the result of getting a block device based on an NVMe device
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

// BdevGetBdevsParams is the parameters required to get a block device
type BdevGetBdevsParams struct {
	Name string `json:"name"`
}

// BdevGetBdevsResult is the result of getting a block device
type BdevGetBdevsResult struct {
	Name      string `json:"name"`
	BlockSize int64  `json:"block_size"`
	NumBlocks int64  `json:"num_blocks"`
	UUID      string `json:"uuid"`
}

// BdevGetIostatParams hold the parameters required to get the IO stats of a block device
type BdevGetIostatParams struct {
	Name string `json:"name"`
}

// BdevGetIostatResult hold the results of getting the IO stats of a block device
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

// VhostCreateBlkControllerParams holds the parameters required to create a block device
// from a vhost controller
type VhostCreateBlkControllerParams struct {
	Ctrlr   string `json:"ctrlr"`
	DevName string `json:"dev_name"`
}

// VhostCreateBlkControllerResult is the result of creating a block device from a vhost controller
type VhostCreateBlkControllerResult bool

// VhostDeleteControllerParams holds the parameters required to delete a vhost controller
type VhostDeleteControllerParams struct {
	Ctrlr string `json:"ctrlr"`
}

// VhostDeleteControllerResult is the result of deleting a vhost controller
type VhostDeleteControllerResult bool

// VhostGetControllersParams holds the parameters required to get a vhost controller
type VhostGetControllersParams struct {
	Name string `json:"name"`
}

// VhostGetControllersResult is the result of getting a vhost controller
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

// VhostCreateScsiControllerParams holds the parameters required to create a SCSI controller
type VhostCreateScsiControllerParams struct {
	Ctrlr string `json:"ctrlr"`
}

// VhostCreateScsiControllerResult is the result of creating a SCSI controller
type VhostCreateScsiControllerResult bool

// NvmfSubsystemAddNsParams holds the parameters required to add a namespace to an existing subsystem
type NvmfSubsystemAddNsParams struct {
	Nqn       string `json:"nqn"`
	Namespace struct {
		BdevName string `json:"bdev_name"`
	} `json:"namespace"`
}

// NvmfSubsystemAddNsResult is the result NSID of attaching a namespace to an existing subsystem
type NvmfSubsystemAddNsResult int

// NvmfSubsystemRemoveNsParams holds the parameters required to Delete a NVMf subsystem
type NvmfSubsystemRemoveNsParams struct {
	Nqn  string `json:"nqn"`
	Nsid int    `json:"nsid"`
}

// NvmfSubsystemRemoveNsResult is the result of creating a NVMf subsystem
type NvmfSubsystemRemoveNsResult bool
