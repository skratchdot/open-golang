//go:build !windows && !darwin

package open

import (
	"context"
	"os/exec"

	"github.com/skratchdot/open-golang/internal/cmdutil"
	"github.com/skratchdot/open-golang/opts"
)

// http://sources.debian.net/src/xdg-utils/1.1.0~rc1%2Bgit20111210-7.1/scripts/xdg-open/
// http://sources.debian.net/src/xdg-utils/1.1.0~rc1%2Bgit20111210-7.1/scripts/xdg-mime/

func createCmd(ctx context.Context, target string, options ...opts.Option) (opts.Config, *exec.Cmd) {
	return cmdutil.CmdUnix(ctx, target, options...)
}
