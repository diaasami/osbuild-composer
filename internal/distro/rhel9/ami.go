package rhel9

import (
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/subscription"
)

const amiKernelOptions = "console=ttyS0,115200n8 console=tty0 net.ifnames=0 rd.blacklist=nouveau nvme_core.io_timeout=4294967295"

var (
	amiImgTypeX86_64 = imageType{
		name:     "ami",
		filename: "image.raw",
		mimeType: "application/octet-stream",
		packageSets: map[string]packageSetFunc{
			osPkgsKey: ec2CommonPackageSet,
		},
		kernelOptions:       amiKernelOptions,
		bootable:            true,
		defaultSize:         10 * common.GibiByte,
		image:               liveImage,
		buildPipelines:      []string{"build"},
		payloadPipelines:    []string{"os", "image"},
		exports:             []string{"image"},
		basePartitionTables: defaultBasePartitionTables,
	}

	ec2ImgTypeX86_64 = imageType{
		name:        "ec2",
		filename:    "image.raw.xz",
		mimeType:    "application/xz",
		compression: "xz",
		packageSets: map[string]packageSetFunc{
			osPkgsKey: rhelEc2PackageSet,
		},
		kernelOptions:       amiKernelOptions,
		bootable:            true,
		defaultSize:         10 * common.GibiByte,
		image:               liveImage,
		buildPipelines:      []string{"build"},
		payloadPipelines:    []string{"os", "image", "xz"},
		exports:             []string{"xz"},
		basePartitionTables: defaultBasePartitionTables,
	}

	ec2HaImgTypeX86_64 = imageType{
		name:        "ec2-ha",
		filename:    "image.raw.xz",
		mimeType:    "application/xz",
		compression: "xz",
		packageSets: map[string]packageSetFunc{
			buildPkgsKey: ec2BuildPackageSet,
			osPkgsKey:    rhelEc2HaPackageSet,
		},
		kernelOptions:       amiKernelOptions,
		bootable:            true,
		defaultSize:         10 * common.GibiByte,
		image:               liveImage,
		buildPipelines:      []string{"build"},
		payloadPipelines:    []string{"os", "image", "xz"},
		exports:             []string{"xz"},
		basePartitionTables: defaultBasePartitionTables,
	}

	amiImgTypeAarch64 = imageType{
		name:     "ami",
		filename: "image.raw",
		mimeType: "application/octet-stream",
		packageSets: map[string]packageSetFunc{
			buildPkgsKey: ec2BuildPackageSet,
			osPkgsKey:    ec2CommonPackageSet,
		},
		kernelOptions:       "console=ttyS0,115200n8 console=tty0 net.ifnames=0 rd.blacklist=nouveau nvme_core.io_timeout=4294967295 iommu.strict=0",
		bootable:            true,
		defaultSize:         10 * common.GibiByte,
		image:               liveImage,
		buildPipelines:      []string{"build"},
		payloadPipelines:    []string{"os", "image"},
		exports:             []string{"image"},
		basePartitionTables: defaultBasePartitionTables,
	}

	ec2ImgTypeAarch64 = imageType{
		name:        "ec2",
		filename:    "image.raw.xz",
		mimeType:    "application/xz",
		compression: "xz",
		packageSets: map[string]packageSetFunc{
			buildPkgsKey: ec2BuildPackageSet,
			osPkgsKey:    rhelEc2PackageSet,
		},
		kernelOptions:       "console=ttyS0,115200n8 console=tty0 net.ifnames=0 rd.blacklist=nouveau nvme_core.io_timeout=4294967295 iommu.strict=0",
		bootable:            true,
		defaultSize:         10 * common.GibiByte,
		image:               liveImage,
		buildPipelines:      []string{"build"},
		payloadPipelines:    []string{"os", "image", "xz"},
		exports:             []string{"xz"},
		basePartitionTables: defaultBasePartitionTables,
	}

	ec2SapImgTypeX86_64 = imageType{
		name:        "ec2-sap",
		filename:    "image.raw.xz",
		mimeType:    "application/xz",
		compression: "xz",
		packageSets: map[string]packageSetFunc{
			buildPkgsKey: ec2BuildPackageSet,
			osPkgsKey:    rhelEc2SapPackageSet,
		},
		kernelOptions:       "console=ttyS0,115200n8 console=tty0 net.ifnames=0 rd.blacklist=nouveau nvme_core.io_timeout=4294967295 processor.max_cstate=1 intel_idle.max_cstate=1",
		bootable:            true,
		defaultSize:         10 * common.GibiByte,
		image:               liveImage,
		buildPipelines:      []string{"build"},
		payloadPipelines:    []string{"os", "image", "xz"},
		exports:             []string{"xz"},
		basePartitionTables: defaultBasePartitionTables,
	}
)

// default EC2 images config (common for all architectures)
func baseEc2ImageConfig() *distro.ImageConfig {
	return &distro.ImageConfig{
		Locale:   common.ToPtr("en_US.UTF-8"),
		Timezone: common.ToPtr("UTC"),
		TimeSynchronization: &osbuild.ChronyStageOptions{
			Servers: []osbuild.ChronyConfigServer{
				{
					Hostname: "169.254.169.123",
					Prefer:   common.ToPtr(true),
					Iburst:   common.ToPtr(true),
					Minpoll:  common.ToPtr(4),
					Maxpoll:  common.ToPtr(4),
				},
			},
			// empty string will remove any occurrences of the option from the configuration
			LeapsecTz: common.ToPtr(""),
		},
		Keyboard: &osbuild.KeymapStageOptions{
			Keymap: "us",
			X11Keymap: &osbuild.X11KeymapOptions{
				Layouts: []string{"us"},
			},
		},
		EnabledServices: []string{
			"sshd",
			"NetworkManager",
			"nm-cloud-setup.service",
			"nm-cloud-setup.timer",
			"cloud-init",
			"cloud-init-local",
			"cloud-config",
			"cloud-final",
			"reboot.target",
			"tuned",
		},
		DefaultTarget: common.ToPtr("multi-user.target"),
		Sysconfig: []*osbuild.SysconfigStageOptions{
			{
				Kernel: &osbuild.SysconfigKernelOptions{
					UpdateDefault: true,
					DefaultKernel: "kernel",
				},
				Network: &osbuild.SysconfigNetworkOptions{
					Networking: true,
					NoZeroConf: true,
				},
				NetworkScripts: &osbuild.NetworkScriptsOptions{
					IfcfgFiles: map[string]osbuild.IfcfgFile{
						"eth0": {
							Device:    "eth0",
							Bootproto: osbuild.IfcfgBootprotoDHCP,
							OnBoot:    common.ToPtr(true),
							Type:      osbuild.IfcfgTypeEthernet,
							UserCtl:   common.ToPtr(true),
							PeerDNS:   common.ToPtr(true),
							IPv6Init:  common.ToPtr(false),
						},
					},
				},
			},
		},
		SystemdLogind: []*osbuild.SystemdLogindStageOptions{
			{
				Filename: "00-getty-fixes.conf",
				Config: osbuild.SystemdLogindConfigDropin{

					Login: osbuild.SystemdLogindConfigLoginSection{
						NAutoVTs: common.ToPtr(0),
					},
				},
			},
		},
		CloudInit: []*osbuild.CloudInitStageOptions{
			{
				Filename: "00-rhel-default-user.cfg",
				Config: osbuild.CloudInitConfigFile{
					SystemInfo: &osbuild.CloudInitConfigSystemInfo{
						DefaultUser: &osbuild.CloudInitConfigDefaultUser{
							Name: "ec2-user",
						},
					},
				},
			},
		},
		Modprobe: []*osbuild.ModprobeStageOptions{
			{
				Filename: "blacklist-nouveau.conf",
				Commands: osbuild.ModprobeConfigCmdList{
					osbuild.NewModprobeConfigCmdBlacklist("nouveau"),
				},
			},
			{
				Filename: "blacklist-amdgpu.conf",
				Commands: osbuild.ModprobeConfigCmdList{
					osbuild.NewModprobeConfigCmdBlacklist("amdgpu"),
				},
			},
		},
		// COMPOSER-1807
		DracutConf: []*osbuild.DracutConfStageOptions{
			{
				Filename: "sgdisk.conf",
				Config: osbuild.DracutConfigFile{
					Install: []string{"sgdisk"},
				},
			},
		},
		SystemdUnit: []*osbuild.SystemdUnitStageOptions{
			// RHBZ#1822863
			{
				Unit:   "nm-cloud-setup.service",
				Dropin: "10-rh-enable-for-ec2.conf",
				Config: osbuild.SystemdServiceUnitDropin{
					Service: &osbuild.SystemdUnitServiceSection{
						Environment: "NM_CLOUD_SETUP_EC2=yes",
					},
				},
			},
		},
		Authselect: &osbuild.AuthselectStageOptions{
			Profile: "sssd",
		},
		SshdConfig: &osbuild.SshdConfigStageOptions{
			Config: osbuild.SshdConfigConfig{
				PasswordAuthentication: common.ToPtr(false),
			},
		},
	}
}

func defaultEc2ImageConfig(osVersion string, rhsm bool) *distro.ImageConfig {
	ic := baseEc2ImageConfig()
	if rhsm && common.VersionLessThan(osVersion, "9.1") {
		ic = appendRHSM(ic)
		// Disable RHSM redhat.repo management
		rhsmConf := ic.RHSMConfig[subscription.RHSMConfigNoSubscription]
		rhsmConf.SubMan.Rhsm = &osbuild.SubManConfigRHSMSection{ManageRepos: common.ToPtr(false)}
		ic.RHSMConfig[subscription.RHSMConfigNoSubscription] = rhsmConf
	}
	return ic
}

// default AMI (EC2 BYOS) images config
func defaultAMIImageConfig(osVersion string, rhsm bool) *distro.ImageConfig {
	ic := defaultEc2ImageConfig(osVersion, rhsm)
	if rhsm {
		// defaultAMIImageConfig() adds the rhsm options only for RHEL < 9.1
		// Add it unconditionally for AMI
		ic = appendRHSM(ic)
	}
	return ic
}

func defaultEc2ImageConfigX86_64(osVersion string, rhsm bool) *distro.ImageConfig {
	ic := defaultEc2ImageConfig(osVersion, rhsm)
	return appendEC2DracutX86_64(ic)
}

func defaultAMIImageConfigX86_64(osVersion string, rhsm bool) *distro.ImageConfig {
	ic := defaultAMIImageConfig(osVersion, rhsm).InheritFrom(defaultEc2ImageConfigX86_64(osVersion, rhsm))
	return appendEC2DracutX86_64(ic)
}

// common ec2 image build package set
func ec2BuildPackageSet(t *imageType) rpmmd.PackageSet {
	return distroBuildPackageSet(t).Append(
		rpmmd.PackageSet{
			Include: []string{
				"python3-pyyaml",
			},
		})
}

func ec2CommonPackageSet(t *imageType) rpmmd.PackageSet {
	return rpmmd.PackageSet{
		Include: []string{
			"authselect-compat",
			"chrony",
			"cloud-init",
			"cloud-utils-growpart",
			"dhcp-client",
			"yum-utils",
			"dracut-config-generic",
			"gdisk",
			"grub2",
			"langpacks-en",
			"NetworkManager-cloud-setup",
			"redhat-release",
			"redhat-release-eula",
			"rsync",
			"tar",
		},
		Exclude: []string{
			"aic94xx-firmware",
			"alsa-firmware",
			"alsa-tools-firmware",
			"biosdevname",
			"iprutils",
			"ivtv-firmware",
			"libertas-sd8787-firmware",
			"plymouth",
			// RHBZ#2064087
			"dracut-config-rescue",
			// RHBZ#2075815
			"qemu-guest-agent",
		},
	}.Append(coreOsCommonPackageSet(t)).Append(distroSpecificPackageSet(t))
}

// common rhel ec2 RHUI image package set
func rhelEc2CommonPackageSet(t *imageType) rpmmd.PackageSet {
	ps := ec2CommonPackageSet(t)
	// Include "redhat-cloud-client-configuration" on 9.1+ (COMPOSER-1805)
	if !common.VersionLessThan(t.arch.distro.osVersion, "9.1") {
		ps.Include = append(ps.Include, "redhat-cloud-client-configuration")
	}
	return ps
}

// rhel-ec2 image package set
func rhelEc2PackageSet(t *imageType) rpmmd.PackageSet {
	ec2PackageSet := rhelEc2CommonPackageSet(t)
	ec2PackageSet = ec2PackageSet.Append(rpmmd.PackageSet{
		Include: []string{
			"rh-amazon-rhui-client",
		},
		Exclude: []string{
			"alsa-lib",
		},
	})
	return ec2PackageSet
}

// rhel-ha-ec2 image package set
func rhelEc2HaPackageSet(t *imageType) rpmmd.PackageSet {
	ec2HaPackageSet := rhelEc2CommonPackageSet(t)
	ec2HaPackageSet = ec2HaPackageSet.Append(rpmmd.PackageSet{
		Include: []string{
			"fence-agents-all",
			"pacemaker",
			"pcs",
			"rh-amazon-rhui-client-ha",
		},
		Exclude: []string{
			"alsa-lib",
		},
	})
	return ec2HaPackageSet
}

// rhel-sap-ec2 image package set
// Includes the common ec2 package set, the common SAP packages, and
// the amazon rhui sap package
func rhelEc2SapPackageSet(t *imageType) rpmmd.PackageSet {
	return rpmmd.PackageSet{
		Include: []string{
			"rh-amazon-rhui-client-sap-bundle-e4s",
		},
	}.Append(rhelEc2CommonPackageSet(t)).Append(SapPackageSet(t))
}

func mkEc2ImgTypeX86_64(osVersion string, rhsm bool) imageType {
	it := ec2ImgTypeX86_64
	ic := defaultEc2ImageConfigX86_64(osVersion, rhsm)
	it.defaultImageConfig = ic
	return it
}

func mkAMIImgTypeX86_64(osVersion string, rhsm bool) imageType {
	it := amiImgTypeX86_64
	ic := defaultAMIImageConfigX86_64(osVersion, rhsm)
	it.defaultImageConfig = ic
	return it
}

func mkEC2SapImgTypeX86_64(osVersion string, rhsm bool) imageType {
	it := ec2SapImgTypeX86_64
	it.defaultImageConfig = sapImageConfig(osVersion).InheritFrom(defaultEc2ImageConfigX86_64(osVersion, rhsm))
	return it
}

func mkEc2HaImgTypeX86_64(osVersion string, rhsm bool) imageType {
	it := ec2HaImgTypeX86_64
	ic := defaultEc2ImageConfigX86_64(osVersion, rhsm)
	it.defaultImageConfig = ic
	return it
}

func mkAMIImgTypeAarch64(osVersion string, rhsm bool) imageType {
	it := amiImgTypeAarch64
	ic := defaultAMIImageConfig(osVersion, rhsm)
	it.defaultImageConfig = ic
	return it
}

func mkEC2ImgTypeAarch64(osVersion string, rhsm bool) imageType {
	it := ec2ImgTypeAarch64
	ic := defaultEc2ImageConfig(osVersion, rhsm)
	it.defaultImageConfig = ic
	return it
}

// Add RHSM config options to ImageConfig.
// Used for RHEL distros.
func appendRHSM(ic *distro.ImageConfig) *distro.ImageConfig {
	rhsm := &distro.ImageConfig{
		RHSMConfig: map[subscription.RHSMStatus]*osbuild.RHSMStageOptions{
			subscription.RHSMConfigNoSubscription: {
				// RHBZ#1932802
				SubMan: &osbuild.RHSMStageOptionsSubMan{
					Rhsmcertd: &osbuild.SubManConfigRHSMCERTDSection{
						AutoRegistration: common.ToPtr(true),
					},
					// Don't disable RHSM redhat.repo management on the AMI
					// image, which is BYOS and does not use RHUI for content.
					// Otherwise subscribing the system manually after booting
					// it would result in empty redhat.repo. Without RHUI, such
					// system would have no way to get Red Hat content, but
					// enable the repo management manually, which would be very
					// confusing.
				},
			},
			subscription.RHSMConfigWithSubscription: {
				// RHBZ#1932802
				SubMan: &osbuild.RHSMStageOptionsSubMan{
					Rhsmcertd: &osbuild.SubManConfigRHSMCERTDSection{
						AutoRegistration: common.ToPtr(true),
					},
					// do not disable the redhat.repo management if the user
					// explicitly request the system to be subscribed
				},
			},
		},
	}
	return rhsm.InheritFrom(ic)
}

func appendEC2DracutX86_64(ic *distro.ImageConfig) *distro.ImageConfig {
	ic.DracutConf = append(ic.DracutConf,
		&osbuild.DracutConfStageOptions{
			Filename: "ec2.conf",
			Config: osbuild.DracutConfigFile{
				AddDrivers: []string{
					"nvme",
					"xen-blkfront",
				},
			},
		})
	return ic
}
