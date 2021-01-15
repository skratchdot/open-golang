// +build windows
//go:generate mkwinsyscall -output zexec_windows.go exec_windows.go
//sys ShellExecute(hwnd int, verb string, file string, args string, cwd string, showCmd int) (err error) = shell32.ShellExecuteW

package open

import (
	"os/exec"
	"strings"
)

const SW_SHOWNORMAL = 1

func cleaninput(input string) string {
	r := strings.NewReplacer("&", "^&")
	return r.Replace(input)
}

func open(input string) *exec.Cmd {
   ShellExecute(0, "", input, "", "", SW_SHOWNORMAL)
   return nil
}

func openWith(input string, appName string) *exec.Cmd {
   ShellExecute(0, "", appName, input, "", SW_SHOWNORMAL)
   return nil
}
