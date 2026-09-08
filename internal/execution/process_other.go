//go:build !windows

package execution

import "os/exec"

func configureProcess(_ *exec.Cmd) {}
