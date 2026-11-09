package cmdutil

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/skratchdot/open-golang/opts"
)

const (
	windowsRunDll32Arg = "url.dll,FileProtocolHandler"
)

// CmdWindows returns the configuration and command used to open target on
// Windows, optionally using an application from options.
func CmdWindows(ctx context.Context, target string, options ...opts.Option) (opts.Config, *exec.Cmd) {
	windowsRunDll32 := getWindowsRunDll32()
	cfg := opts.Init(
		opts.ConfigBuilder{}.AppName(windowsRunDll32),
		options...,
	)

	args := make([]string, 0)
	appName := cfg.AppName()
	if appName == windowsRunDll32 {
		// If we're using the runDll32 executable, we need to follow up
		// with the file protocol handler.
		args = append(args, windowsRunDll32Arg)
	} else if appName != "" {
		// If we're invoking a specific application, we need
		// to change the strategy to invoke the app via cmd
		args = append(args, "/C", "start", "", appName)
		appName = "cmd"
	}
	args = append(args, cleanWindowsInput(target))
	return cfg, exec.CommandContext(ctx, appName, args...)
}

func getWindowsRunDll32() string {
	return filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
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
