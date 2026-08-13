package git

import (
	"os/exec"
)

func Stage(path string) error {
	cmd := exec.Command("git", "add", "--", path)
	return cmd.Run()
}

func Unstage(path string) error {
	cmd := exec.Command("git", "reset", "HEAD", "--", path)
	return cmd.Run()
}
