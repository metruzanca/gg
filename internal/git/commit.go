package git

import (
	"os/exec"
)

func Commit(spec CommitSpec) error {
	message := spec.Type + ": " + spec.Title
	args := []string{"commit", "-m", message}
	if spec.Description != "" {
		args = append(args, "-m", spec.Description)
	}
	cmd := exec.Command("git", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
