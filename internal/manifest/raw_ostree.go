package manifest

import (
	"fmt"

	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/platform"
)

// A RawOSTreeImage represents a raw ostree image file which can be booted in a
// hypervisor. It is created from an existing OSTreeDeployment.
type RawOSTreeImage struct {
	Base
	treePipeline *OSTreeDeployment
	Filename     string
	platform     platform.Platform
}

func NewRawOStreeImage(m *Manifest,
	buildPipeline *Build,
	platform platform.Platform,
	treePipeline *OSTreeDeployment) *RawOSTreeImage {
	p := &RawOSTreeImage{
		Base:         NewBase(m, "image", buildPipeline),
		treePipeline: treePipeline,
		Filename:     "disk.img",
		platform:     platform,
	}
	buildPipeline.addDependent(p)
	if treePipeline.Base.manifest != m {
		panic("tree pipeline from different manifest")
	}
	m.addPipeline(p)
	return p
}

func (p *RawOSTreeImage) getBuildPackages() []string {
	packages := p.platform.GetBuildPackages()
	packages = append(packages, p.platform.GetPackages()...)
	packages = append(packages, p.treePipeline.PartitionTable.GetBuildPackages()...)
	packages = append(packages,
		"rpm-ostree",

		// these should be defined on the platform
		"dracut-config-generic",
		"efibootmgr",
	)
	return packages
}

func (p *RawOSTreeImage) serialize() osbuild.Pipeline {
	pipeline := p.Base.serialize()

	pt := p.treePipeline.PartitionTable
	if pt == nil {
		panic("no partition table in live image")
	}

	for _, stage := range osbuild.GenImagePrepareStages(pt, p.Filename, osbuild.PTSfdisk) {
		pipeline.AddStage(stage)
	}

	inputName := "root-tree"
	treeCopyOptions, treeCopyDevices, treeCopyMounts := osbuild.GenCopyFSTreeOptions(inputName, p.treePipeline.Name(), p.Filename, pt)
	treeCopyInputs := osbuild.NewPipelineTreeInputs(inputName, p.treePipeline.Name())

	pipeline.AddStage(osbuild.NewCopyStage(treeCopyOptions, treeCopyInputs, treeCopyDevices, treeCopyMounts))

	bootFiles := p.platform.GetBootFiles()
	if len(bootFiles) > 0 {
		// we ignore the bootcopyoptions as they contain a full tree copy instead we make our own, we *do* still want all the other
		// information such as mountpoints and devices
		_, bootCopyDevices, bootCopyMounts := osbuild.GenCopyFSTreeOptions(inputName, p.treePipeline.Name(), p.Filename, pt)
		bootCopyOptions := &osbuild.CopyStageOptions{}

		commitChecksum := p.treePipeline.commit.Checksum

		bootCopyInputs := osbuild.OSTreeCheckoutInputs{
			"ostree-tree": *osbuild.NewOSTreeCheckoutInput("org.osbuild.source", commitChecksum),
		}

		for _, paths := range bootFiles {
			bootCopyOptions.Paths = append(bootCopyOptions.Paths, osbuild.CopyStagePath{
				From: fmt.Sprintf("input://ostree-tree/%s%s", commitChecksum, paths[0]),
				To:   fmt.Sprintf("mount://root%s", paths[1]),
			})
		}

		pipeline.AddStage(osbuild.NewCopyStage(bootCopyOptions, bootCopyInputs, bootCopyDevices, bootCopyMounts))
	}

	for _, stage := range osbuild.GenImageFinishStages(pt, p.Filename) {
		pipeline.AddStage(stage)
	}

	if grubLegacy := p.treePipeline.platform.GetBIOSPlatform(); grubLegacy != "" {
		pipeline.AddStage(osbuild.NewGrub2InstStage(osbuild.NewGrub2InstStageOption(p.Filename, pt, grubLegacy)))
	}

	return pipeline
}

func (p *RawOSTreeImage) Export() *artifact.Artifact {
	p.Base.export = true
	return artifact.New(p.Name(), p.Filename, nil)
}
