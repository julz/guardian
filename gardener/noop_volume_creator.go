package gardener

import (
	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
)

type NoopVolumeCreator struct{}

var ErrGraphDisabled = errors.New("volume graph is disabled")

func (NoopVolumeCreator) Create(lager.Logger, string, DesiredRootFSSpec) (string, []string, error) {
	return "", nil, ErrGraphDisabled
}

func (NoopVolumeCreator) Destroy(lager.Logger, string, string) error {
	return nil
}

func (NoopVolumeCreator) Metrics(lager.Logger, string, string) (garden.ContainerDiskStat, error) {
	return garden.ContainerDiskStat{}, nil
}

func (NoopVolumeCreator) GC(lager.Logger) error {
	return nil
}
