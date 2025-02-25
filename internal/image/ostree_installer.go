package image

import (
	"fmt"
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/ostree"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/runner"
	"github.com/osbuild/osbuild-composer/internal/users"
)

type OSTreeInstaller struct {
	Base
	Platform          platform.Platform
	ExtraBasePackages rpmmd.PackageSet
	Users             []users.User
	Groups            []users.Group

	SquashfsCompression string

	ISOLabelTempl string
	Product       string
	Variant       string
	OSName        string
	OSVersion     string
	Release       string

	Commit ostree.CommitSpec

	Filename string

	AdditionalDracutModules   []string
	AdditionalAnacondaModules []string
	AdditionalDrivers         []string
}

func NewOSTreeInstaller(commit ostree.CommitSpec) *OSTreeInstaller {
	return &OSTreeInstaller{
		Base:   NewBase("ostree-installer"),
		Commit: commit,
	}
}

func (img *OSTreeInstaller) InstantiateManifest(m *manifest.Manifest,
	repos []rpmmd.RepoConfig,
	runner runner.Runner,
	rng *rand.Rand) (*artifact.Artifact, error) {
	buildPipeline := manifest.NewBuild(m, runner, repos)
	buildPipeline.Checkpoint()

	anacondaPipeline := manifest.NewAnaconda(m,
		buildPipeline,
		img.Platform,
		repos,
		"kernel",
		img.Product,
		img.OSVersion)
	anacondaPipeline.ExtraPackages = img.ExtraBasePackages.Include
	anacondaPipeline.ExtraRepos = img.ExtraBasePackages.Repositories
	anacondaPipeline.Users = img.Users
	anacondaPipeline.Groups = img.Groups
	anacondaPipeline.Variant = img.Variant
	anacondaPipeline.Biosdevname = (img.Platform.GetArch() == platform.ARCH_X86_64)
	anacondaPipeline.Checkpoint()
	anacondaPipeline.AdditionalDracutModules = img.AdditionalDracutModules
	anacondaPipeline.AdditionalAnacondaModules = img.AdditionalAnacondaModules
	anacondaPipeline.AdditionalDrivers = img.AdditionalDrivers

	rootfsPartitionTable := &disk.PartitionTable{
		Size: 20 * common.MebiByte,
		Partitions: []disk.Partition{
			{
				Start: 0,
				Size:  20 * common.MebiByte,
				Payload: &disk.Filesystem{
					Type:       "vfat",
					Mountpoint: "/",
					UUID:       disk.NewVolIDFromRand(rng),
				},
			},
		},
	}

	// TODO: replace isoLabelTmpl with more high-level properties
	isoLabel := fmt.Sprintf(img.ISOLabelTempl, img.Platform.GetArch())

	rootfsImagePipeline := manifest.NewISORootfsImg(m, buildPipeline, anacondaPipeline)
	rootfsImagePipeline.Size = 4 * common.GibiByte

	bootTreePipeline := manifest.NewEFIBootTree(m, buildPipeline, img.Product, img.OSVersion)
	bootTreePipeline.Platform = img.Platform
	bootTreePipeline.UEFIVendor = img.Platform.GetUEFIVendor()
	bootTreePipeline.ISOLabel = isoLabel
	bootTreePipeline.KernelOpts = []string{fmt.Sprintf("inst.stage2=hd:LABEL=%s", isoLabel), fmt.Sprintf("inst.ks=hd:LABEL=%s:%s", isoLabel, kspath)}

	// enable ISOLinux on x86_64 only
	isoLinuxEnabled := img.Platform.GetArch() == platform.ARCH_X86_64

	isoTreePipeline := manifest.NewAnacondaISOTree(m,
		buildPipeline,
		anacondaPipeline,
		rootfsImagePipeline,
		bootTreePipeline,
		isoLabel)
	isoTreePipeline.PartitionTable = rootfsPartitionTable
	isoTreePipeline.Release = img.Release
	isoTreePipeline.OSName = img.OSName
	isoTreePipeline.Users = img.Users
	isoTreePipeline.Groups = img.Groups

	isoTreePipeline.SquashfsCompression = img.SquashfsCompression

	// For ostree installers, always put the kickstart file in the root of the ISO
	isoTreePipeline.KSPath = kspath
	isoTreePipeline.PayloadPath = "/ostree/repo"

	isoTreePipeline.OSTree = &img.Commit
	isoTreePipeline.ISOLinux = isoLinuxEnabled

	isoPipeline := manifest.NewISO(m, buildPipeline, isoTreePipeline, isoLabel)
	isoPipeline.Filename = img.Filename
	isoPipeline.ISOLinux = isoLinuxEnabled
	artifact := isoPipeline.Export()

	return artifact, nil
}
