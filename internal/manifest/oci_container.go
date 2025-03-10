package manifest

import (
	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
)

// An OCIContainer represents an OCI container, containing a filesystem
// tree created by another Pipeline.
type OCIContainer struct {
	Base
	Filename     string
	Cmd          []string
	ExposedPorts []string

	treePipeline Tree
}

func NewOCIContainer(m *Manifest,
	buildPipeline *Build,
	treePipeline Tree) *OCIContainer {
	p := &OCIContainer{
		Base:         NewBase(m, "container", buildPipeline),
		treePipeline: treePipeline,
		Filename:     "oci-archive.tar",
	}
	if treePipeline.GetManifest() != m {
		panic("tree pipeline from different manifest")
	}
	buildPipeline.addDependent(p)
	m.addPipeline(p)
	return p
}

func (p *OCIContainer) serialize() osbuild.Pipeline {
	pipeline := p.Base.serialize()

	options := &osbuild.OCIArchiveStageOptions{
		Architecture: p.treePipeline.GetPlatform().GetArch().String(),
		Filename:     p.Filename,
		Config: &osbuild.OCIArchiveConfig{
			Cmd:          p.Cmd,
			ExposedPorts: p.ExposedPorts,
		},
	}
	baseInput := osbuild.NewTreeInput("name:" + p.treePipeline.Name())
	inputs := &osbuild.OCIArchiveStageInputs{Base: baseInput}
	pipeline.AddStage(osbuild.NewOCIArchiveStage(options, inputs))

	return pipeline
}

func (p *OCIContainer) getBuildPackages() []string {
	return []string{"tar"}
}

func (p *OCIContainer) Export() *artifact.Artifact {
	p.Base.export = true
	mimeType := "application/x-tar"
	return artifact.New(p.Name(), p.Filename, &mimeType)
}
