package rhel9

import (
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/subscription"
)

const gceKernelOptions = "net.ifnames=0 biosdevname=0 scsi_mod.use_blk_mq=Y console=ttyS0,38400n8d"

var (
	gceImgType = imageType{
		name:     "gce",
		filename: "image.tar.gz",
		mimeType: "application/gzip",
		packageSets: map[string]packageSetFunc{
			osPkgsKey: gcePackageSet,
		},
		kernelOptions:    gceKernelOptions,
		bootable:         true,
		defaultSize:      20 * common.GibiByte,
		image:            liveImage,
		buildPipelines:   []string{"build"},
		payloadPipelines: []string{"os", "image", "archive"},
		exports:          []string{"archive"},
		// TODO: the base partition table still contains the BIOS boot partition, but the image is UEFI-only
		basePartitionTables: defaultBasePartitionTables,
	}

	gceRhuiImgType = imageType{
		name:     "gce-rhui",
		filename: "image.tar.gz",
		mimeType: "application/gzip",
		packageSets: map[string]packageSetFunc{
			osPkgsKey: gceRhuiPackageSet,
		},
		kernelOptions:    gceKernelOptions,
		bootable:         true,
		defaultSize:      20 * common.GibiByte,
		image:            liveImage,
		buildPipelines:   []string{"build"},
		payloadPipelines: []string{"os", "image", "archive"},
		exports:          []string{"archive"},
		// TODO: the base partition table still contains the BIOS boot partition, but the image is UEFI-only
		basePartitionTables: defaultBasePartitionTables,
	}
)

func mkGCEImageType(rhsm bool) imageType {
	it := gceImgType
	it.defaultImageConfig = baseGCEImageConfig(rhsm)
	return it
}

func mkGCERHUIImageType(rhsm bool) imageType {
	it := gceRhuiImgType
	it.defaultImageConfig = defaultGceRhuiImageConfig(rhsm)
	return it
}

func baseGCEImageConfig(rhsm bool) *distro.ImageConfig {
	ic := &distro.ImageConfig{
		Timezone: common.ToPtr("UTC"),
		TimeSynchronization: &osbuild.ChronyStageOptions{
			Servers: []osbuild.ChronyConfigServer{{Hostname: "metadata.google.internal"}},
		},
		Firewall: &osbuild.FirewallStageOptions{
			DefaultZone: "trusted",
		},
		EnabledServices: []string{
			"sshd",
			"rngd",
			"dnf-automatic.timer",
		},
		DisabledServices: []string{
			"sshd-keygen@",
			"reboot.target",
		},
		DefaultTarget: common.ToPtr("multi-user.target"),
		Locale:        common.ToPtr("en_US.UTF-8"),
		Keyboard: &osbuild.KeymapStageOptions{
			Keymap: "us",
		},
		DNFConfig: []*osbuild.DNFConfigStageOptions{
			{
				Config: &osbuild.DNFConfig{
					Main: &osbuild.DNFConfigMain{
						IPResolve: "4",
					},
				},
			},
		},
		DNFAutomaticConfig: &osbuild.DNFAutomaticConfigStageOptions{
			Config: &osbuild.DNFAutomaticConfig{
				Commands: &osbuild.DNFAutomaticConfigCommands{
					ApplyUpdates: common.ToPtr(true),
					UpgradeType:  osbuild.DNFAutomaticUpgradeTypeSecurity,
				},
			},
		},
		YUMRepos: []*osbuild.YumReposStageOptions{
			{
				Filename: "google-cloud.repo",
				Repos: []osbuild.YumRepository{
					{
						Id:       "google-compute-engine",
						Name:     "Google Compute Engine",
						BaseURLs: []string{"https://packages.cloud.google.com/yum/repos/google-compute-engine-el9-x86_64-stable"},
						Enabled:  common.ToPtr(true),
						// TODO: enable GPG check once Google stops using SHA-1 in their keys
						// https://issuetracker.google.com/issues/223626963
						GPGCheck:     common.ToPtr(false),
						RepoGPGCheck: common.ToPtr(false),
						GPGKey: []string{
							"https://packages.cloud.google.com/yum/doc/yum-key.gpg",
							"https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg",
						},
					},
				},
			},
		},
		SshdConfig: &osbuild.SshdConfigStageOptions{
			Config: osbuild.SshdConfigConfig{
				PasswordAuthentication: common.ToPtr(false),
				ClientAliveInterval:    common.ToPtr(420),
				PermitRootLogin:        osbuild.PermitRootLoginValueNo,
			},
		},
		Sysconfig: []*osbuild.SysconfigStageOptions{
			{
				Kernel: &osbuild.SysconfigKernelOptions{
					DefaultKernel: "kernel-core",
					UpdateDefault: true,
				},
			},
		},
		Modprobe: []*osbuild.ModprobeStageOptions{
			{
				Filename: "blacklist-floppy.conf",
				Commands: osbuild.ModprobeConfigCmdList{
					osbuild.NewModprobeConfigCmdBlacklist("floppy"),
				},
			},
		},
		GCPGuestAgentConfig: &osbuild.GcpGuestAgentConfigOptions{
			ConfigScope: osbuild.GcpGuestAgentConfigScopeDistro,
			Config: &osbuild.GcpGuestAgentConfig{
				InstanceSetup: &osbuild.GcpGuestAgentConfigInstanceSetup{
					SetBotoConfig: common.ToPtr(false),
				},
			},
		},
	}

	if rhsm {
		ic.RHSMConfig = map[subscription.RHSMStatus]*osbuild.RHSMStageOptions{
			subscription.RHSMConfigNoSubscription: {
				SubMan: &osbuild.RHSMStageOptionsSubMan{
					Rhsmcertd: &osbuild.SubManConfigRHSMCERTDSection{
						AutoRegistration: common.ToPtr(true),
					},
					// Don't disable RHSM redhat.repo management on the GCE
					// image, which is BYOS and does not use RHUI for content.
					// Otherwise subscribing the system manually after booting
					// it would result in empty redhat.repo. Without RHUI, such
					// system would have no way to get Red Hat content, but
					// enable the repo management manually, which would be very
					// confusing.
				},
			},
			subscription.RHSMConfigWithSubscription: {
				SubMan: &osbuild.RHSMStageOptionsSubMan{
					Rhsmcertd: &osbuild.SubManConfigRHSMCERTDSection{
						AutoRegistration: common.ToPtr(true),
					},
					// do not disable the redhat.repo management if the user
					// explicitly request the system to be subscribed
				},
			},
		}
	}
	return ic
}

func defaultGceRhuiImageConfig(rhsm bool) *distro.ImageConfig {
	ic := &distro.ImageConfig{
		RHSMConfig: map[subscription.RHSMStatus]*osbuild.RHSMStageOptions{
			subscription.RHSMConfigNoSubscription: {
				SubMan: &osbuild.RHSMStageOptionsSubMan{
					Rhsmcertd: &osbuild.SubManConfigRHSMCERTDSection{
						AutoRegistration: common.ToPtr(true),
					},
					Rhsm: &osbuild.SubManConfigRHSMSection{
						ManageRepos: common.ToPtr(false),
					},
				},
			},
			subscription.RHSMConfigWithSubscription: {
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
	return ic.InheritFrom(baseGCEImageConfig(rhsm))
}

func gceCommonPackageSet(t *imageType) rpmmd.PackageSet {
	ps := rpmmd.PackageSet{
		Include: []string{
			"langpacks-en", // not in Google's KS
			"acpid",
			"dhcp-client",
			"dnf-automatic",
			"net-tools",
			//"openssh-server", included in core
			"python3",
			"rng-tools",
			"tar",
			"vim",

			// GCE guest tools
			"google-compute-engine",
			"google-osconfig-agent",
			"gce-disk-expand",

			// Not explicitly included in GCP kickstart, but present on the image
			// for time synchronization
			"chrony",
			"timedatex",
			// EFI
			"grub2-tools-efi",
			"firewalld", // not pulled in any more as on RHEL-8
		},
		Exclude: []string{
			"alsa-utils",
			"b43-fwcutter",
			"dmraid",
			"eject",
			"gpm",
			"irqbalance",
			"microcode_ctl",
			"smartmontools",
			"aic94xx-firmware",
			"atmel-firmware",
			"b43-openfwwf",
			"bfa-firmware",
			"ipw2100-firmware",
			"ipw2200-firmware",
			"ivtv-firmware",
			"iwl100-firmware",
			"iwl1000-firmware",
			"iwl3945-firmware",
			"iwl4965-firmware",
			"iwl5000-firmware",
			"iwl5150-firmware",
			"iwl6000-firmware",
			"iwl6000g2a-firmware",
			"iwl6050-firmware",
			"kernel-firmware",
			"libertas-usb8388-firmware",
			"ql2100-firmware",
			"ql2200-firmware",
			"ql23xx-firmware",
			"ql2400-firmware",
			"ql2500-firmware",
			"rt61pci-firmware",
			"rt73usb-firmware",
			"xorg-x11-drv-ati-firmware",
			"zd1211-firmware",
			// RHBZ#2075815
			"qemu-guest-agent",
		},
	}.Append(coreOsCommonPackageSet(t)).Append(distroSpecificPackageSet(t))

	// Some excluded packages are part of the @core group package set returned
	// by coreOsCommonPackageSet(). Ensure that the conflicting packages are
	// returned from the list of `Include` packages.
	return ps.ResolveConflictsExclude()
}

// GCE BYOS image
func gcePackageSet(t *imageType) rpmmd.PackageSet {
	return gceCommonPackageSet(t)
}

// GCE RHUI image
func gceRhuiPackageSet(t *imageType) rpmmd.PackageSet {
	return rpmmd.PackageSet{
		Include: []string{
			"google-rhui-client-rhel9",
		},
	}.Append(gceCommonPackageSet(t))
}
