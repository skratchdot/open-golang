package cmdutil

import (
	"context"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/skratchdot/open-golang/opts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TC struct {
	name string

	inputTarget string
	inputOpts   []opts.Option
	expectedCmd *exec.Cmd
}

func TestCmdWindows(t *testing.T) {
	testhelperWindows(t)
	windowsRunDll32 := filepath.Join(string([]rune{'N', ':', os.PathSeparator, os.PathSeparator}), "Windows", "System32", "rundll32.exe")

	testcases := []TC{
		{
			name:        "url",
			inputTarget: "https://www.github.com",
			inputOpts:   nil,
			expectedCmd: exec.CommandContext(context.TODO(), windowsRunDll32, windowsRunDll32Arg, "https://www.github.com"),
		},
		{
			name:        "url with app",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.WithApp("chrome"),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "cmd", "/C", "start", "", "chrome", "https://www.github.com"),
		},
		{
			name:        "url with wait",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.Wait(),
			},
			expectedCmd: exec.CommandContext(context.TODO(), windowsRunDll32, windowsRunDll32Arg, "https://www.github.com"),
		},
		{
			name:        "url with app and wait",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.WithApp("chrome"),
				opts.Wait(),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "cmd", "/C", "start", "", "chrome", "https://www.github.com"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, actualCmd := CmdWindows(context.TODO(), tc.inputTarget, tc.inputOpts...)
			require.NotNil(t, actualCmd, "cmd must not be nil")
			assert.Equal(t, tc.expectedCmd.Args, actualCmd.Args)
		})
	}
}

func TestCmdDarwin(t *testing.T) {
	testcases := []TC{
		{
			name:        "url",
			inputTarget: "https://www.github.com",
			inputOpts:   nil,
			expectedCmd: exec.CommandContext(context.TODO(), "open", "https://www.github.com"),
		},
		{
			name:        "url with app",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.WithApp("chrome"),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "open", "-a", "chrome", "https://www.github.com"),
		},
		{
			name:        "url with wait",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.Wait(),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "open", "https://www.github.com"),
		},
		{
			name:        "url with app and wait",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.WithApp("chrome"),
				opts.Wait(),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "open", "-a", "chrome", "https://www.github.com"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, actualCmd := CmdDarwin(context.TODO(), tc.inputTarget, tc.inputOpts...)
			require.NotNil(t, actualCmd, "cmd must not be nil")
			assert.Equal(t, tc.expectedCmd.Args, actualCmd.Args)
		})
	}
}

func TestCmdUnix(t *testing.T) {
	testcases := []TC{
		{
			name:        "url",
			inputTarget: "https://www.github.com",
			inputOpts:   []opts.Option{},
			expectedCmd: exec.CommandContext(context.TODO(), "xdg-open", "https://www.github.com"),
		},
		{
			name:        "url with app",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.WithApp("chrome"),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "chrome", "https://www.github.com"),
		},
		{
			name:        "url with wait",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.Wait(),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "xdg-open", "https://www.github.com"),
		},
		{
			name:        "url with app and wait",
			inputTarget: "https://www.github.com",
			inputOpts: []opts.Option{
				opts.WithApp("chrome"),
				opts.Wait(),
			},
			expectedCmd: exec.CommandContext(context.TODO(), "chrome", "https://www.github.com"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, actualCmd := CmdUnix(context.TODO(), tc.inputTarget, tc.inputOpts...)
			require.NotNil(t, actualCmd, "cmd must not be nil")
			assert.Equal(t, tc.expectedCmd.Args, actualCmd.Args)
		})
	}
}

func testhelperWindows(t *testing.T) {
	t.Helper()
	originalRoot, hasRoot := os.LookupEnv("SYSTEMROOT")
	_ = os.Setenv("SYSTEMROOT", path.Join(`N:`, "Windows"))

	t.Cleanup(func() {
		if !hasRoot {
			_ = os.Unsetenv("SYSTEMROOT")
		} else {
			_ = os.Setenv("SYSTEMROOT", originalRoot)
		}
	})
}
