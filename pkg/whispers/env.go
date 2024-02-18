//go:build amd64

package whispers

import (
	"context"

	"github.com/cilium/ebpf/rlimit"

	"github.com/peter-mcconnell/whispers/pkg/config"
)

func EnvSetup(_ context.Context, _ *config.Config) error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return err
	}
	return nil
}
