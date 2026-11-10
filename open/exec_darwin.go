//go:build darwin

package open

import (
	"context"
	"os/exec"

	"github.com/skratchdot/open-golang/internal/cmdutil"
	"github.com/skratchdot/open-golang/opts"
)

func createCmd(ctx context.Context, target string, options ...opts.Option) (opts.Config, *exec.Cmd) {
	return cmdutil.CmdDarwin(ctx, target, options...)
}
