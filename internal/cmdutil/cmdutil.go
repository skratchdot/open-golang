package cmdutil

import (
	"context"
	"os/exec"
	"strings"

	"github.com/skratchdot/open-golang/opts"
)

// CmdWindows returns the configuration and command used to open target on
// Windows, optionally using an application from options.
func CmdWindows(ctx context.Context, target string, options ...opts.Option) (opts.Config, *exec.Cmd) {
	cfg := opts.Init(
		opts.ConfigBuilder{},
		options...,
	)

	args := make([]string, 0)
	args = append(args, "/C", "start")
	if cfg.Wait() {
		args = append(args, "/wait")
	}
	if appName := cfg.AppName(); appName != "" {
		args = append(args, "", appName)
	}
	args = append(args, cleanWindowsInput(target))
	return cfg, exec.CommandContext(ctx, "cmd", args...)
}

func cleanWindowsInput(input string) string {
	r := strings.NewReplacer("&", "^&")
	return r.Replace(input)
}

// CmdDarwin returns the configuration and command used to open target on
// macOS, optionally using an application from options.
func CmdDarwin(ctx context.Context, target string, options ...opts.Option) (opts.Config, *exec.Cmd) {
	cfg := opts.Init(
		opts.ConfigBuilder{},
		options...,
	)

	args := make([]string, 0)
	if appName := cfg.AppName(); appName != "" {
		args = append(args, "-a", appName)
	}
	args = append(args, target)
	return cfg, exec.CommandContext(ctx, "open", args...)
}

// CmdUnix returns the configuration and command used to open target on
// Unix-like systems, optionally using an application from options.
func CmdUnix(ctx context.Context, target string, options ...opts.Option) (opts.Config, *exec.Cmd) {
	cfg := opts.Init(
		opts.ConfigBuilder{}.AppName("xdg-open"),
		options...,
	)
	return cfg, exec.CommandContext(ctx, cfg.AppName(), target)
}
