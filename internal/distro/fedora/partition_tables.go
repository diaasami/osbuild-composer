package fedora

import (
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/platform"
)

var defaultBasePartitionTables = distro.BasePartitionTableMap{
	platform.ARCH_X86_64.String(): disk.PartitionTable{
		UUID: "D209C89E-EA5E-4FBD-B161-B461CCE297E0",
		Type: "gpt",
		Partitions: []disk.Partition{
			{
				Size:     1 * common.MebiByte, // 1MB
				Bootable: true,
				Type:     disk.BIOSBootPartitionGUID,
				UUID:     disk.BIOSBootPartitionUUID,
			},
			{
				Size: 200 * common.MebiByte, // 200 MB
				Type: disk.EFISystemPartitionGUID,
				UUID: disk.EFISystemPartitionUUID,
				Payload: &disk.Filesystem{
					Type:         "vfat",
					UUID:         disk.EFIFilesystemUUID,
					Mountpoint:   "/boot/efi",
					Label:        "EFI-SYSTEM",
					FSTabOptions: "defaults,uid=0,gid=0,umask=077,shortname=winnt",
					FSTabFreq:    0,
					FSTabPassNo:  2,
				},
			},
			{
				Size: 500 * common.MebiByte, // 500 MB
				Type: disk.FilesystemDataGUID,
				UUID: disk.FilesystemDataUUID,
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Mountpoint:   "/boot",
					Label:        "boot",
					FSTabOptions: "defaults",
					FSTabFreq:    0,
					FSTabPassNo:  0,
				},
			},
			{
				Size: 2 * common.GibiByte, // 2GiB
				Type: disk.FilesystemDataGUID,
				UUID: disk.RootPartitionUUID,
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Label:        "root",
					Mountpoint:   "/",
					FSTabOptions: "defaults",
					FSTabFreq:    0,
					FSTabPassNo:  0,
				},
			},
		},
	},
	platform.ARCH_AARCH64.String(): disk.PartitionTable{
		UUID: "D209C89E-EA5E-4FBD-B161-B461CCE297E0",
		Type: "gpt",
		Partitions: []disk.Partition{
			{
				Size: 200 * common.MebiByte, // 200 MB
				Type: disk.EFISystemPartitionGUID,
				UUID: disk.EFISystemPartitionUUID,
				Payload: &disk.Filesystem{
					Type:         "vfat",
					UUID:         disk.EFIFilesystemUUID,
					Mountpoint:   "/boot/efi",
					Label:        "EFI-SYSTEM",
					FSTabOptions: "defaults,uid=0,gid=0,umask=077,shortname=winnt",
					FSTabFreq:    0,
					FSTabPassNo:  2,
				},
			},
			{
				Size: 500 * common.MebiByte, // 500 MB
				Type: disk.FilesystemDataGUID,
				UUID: disk.FilesystemDataUUID,
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Mountpoint:   "/boot",
					Label:        "boot",
					FSTabOptions: "defaults",
					FSTabFreq:    0,
					FSTabPassNo:  0,
				},
			},
			{
				Size: 2 * common.GibiByte, // 2GiB
				Type: disk.FilesystemDataGUID,
				UUID: disk.RootPartitionUUID,
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Label:        "root",
					Mountpoint:   "/",
					FSTabOptions: "defaults",
					FSTabFreq:    0,
					FSTabPassNo:  0,
				},
			},
		},
	},
}

var iotBasePartitionTables = distro.BasePartitionTableMap{
	platform.ARCH_X86_64.String(): disk.PartitionTable{
		UUID: "D209C89E-EA5E-4FBD-B161-B461CCE297E0",
		Type: "gpt",
		Partitions: []disk.Partition{
			{
				Size: 501 * common.MebiByte, // 501 MiB
				Type: disk.EFISystemPartitionGUID,
				UUID: disk.EFISystemPartitionUUID,
				Payload: &disk.Filesystem{
					Type:         "vfat",
					UUID:         disk.EFIFilesystemUUID,
					Mountpoint:   "/boot/efi",
					Label:        "EFI-SYSTEM",
					FSTabOptions: "umask=0077,shortname=winnt",
					FSTabFreq:    0,
					FSTabPassNo:  2,
				},
			},
			{
				Size: 1 * common.GibiByte, // 1 GiB
				Type: disk.FilesystemDataGUID,
				UUID: disk.FilesystemDataUUID,
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Mountpoint:   "/boot",
					Label:        "boot",
					FSTabOptions: "defaults",
					FSTabFreq:    1,
					FSTabPassNo:  2,
				},
			},
			{
				Size: 2569 * common.MebiByte, // 2.5 GiB
				Type: disk.FilesystemDataGUID,
				UUID: disk.RootPartitionUUID,
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Label:        "root",
					Mountpoint:   "/",
					FSTabOptions: "defaults",
					FSTabFreq:    1,
					FSTabPassNo:  1,
				},
			},
		},
	},
	platform.ARCH_AARCH64.String(): disk.PartitionTable{
		UUID: "0xc1748067",
		Type: "dos",
		Partitions: []disk.Partition{
			{
				Size:     501 * common.MebiByte, // 501 MiB
				Type:     "06",
				Bootable: true,
				Payload: &disk.Filesystem{
					Type:         "vfat",
					UUID:         disk.EFIFilesystemUUID,
					Mountpoint:   "/boot/efi",
					Label:        "EFI-SYSTEM",
					FSTabOptions: "umask=0077,shortname=winnt",
					FSTabFreq:    0,
					FSTabPassNo:  2,
				},
			},
			{
				Size: 1 * common.GibiByte, // 1 GiB
				Type: "83",
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Mountpoint:   "/boot",
					Label:        "boot",
					FSTabOptions: "defaults",
					FSTabFreq:    1,
					FSTabPassNo:  2,
				},
			},
			{
				Size: 2569 * common.MebiByte, // 2.5 GiB
				Type: "83",
				Payload: &disk.Filesystem{
					Type:         "ext4",
					Label:        "root",
					Mountpoint:   "/",
					FSTabOptions: "defaults",
					FSTabFreq:    1,
					FSTabPassNo:  1,
				},
			},
		},
	},
}
