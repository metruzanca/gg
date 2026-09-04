package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func Tags() ([]string, error) {
	cmd := exec.Command("git", "tag", "--list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			tags = append(tags, line)
		}
	}
	return tags, nil
}

func CreateTag(name string) error {
	cmd := exec.Command("git", "tag", name)
	return cmd.Run()
}

func PushTag(name string) error {
	cmd := exec.Command("git", "push", "origin", name)
	return cmd.Run()
}

// CurrentBranch returns the current branch name, or "HEAD" when the working
// tree is on a detached HEAD.
func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// PushBranch pushes the named branch to origin, creating the remote branch if
// it does not exist yet.
func PushBranch(branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	return cmd.Run()
}

// CommitOnRemote reports whether rev is reachable from any remote-tracking
// branch, i.e. whether its commits have been pushed.
func CommitOnRemote(rev string) (bool, error) {
	cmd := exec.Command("git", "branch", "-r", "--contains", rev)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) != "", nil
}

// PushTagWithCommits pushes the commits a tag references before pushing the
// tag itself, so the remote never holds a tag for an unpushed commit. It
// pushes the current branch first; on a detached HEAD it verifies the tagged
// commit is already on a remote branch and refuses to push otherwise.
func PushTagWithCommits(tag string) error {
	branch, err := CurrentBranch()
	if err != nil {
		return fmt.Errorf("resolve current branch: %w", err)
	}
	if branch != "HEAD" {
		if err := PushBranch(branch); err != nil {
			return fmt.Errorf("push branch %q: %w", branch, err)
		}
	} else {
		onRemote, err := CommitOnRemote(tag)
		if err != nil {
			return fmt.Errorf("check remote for %q: %w", tag, err)
		}
		if !onRemote {
			return errors.New("detached HEAD: commit referenced by tag is not on any remote branch, push it first")
		}
	}
	if err := PushTag(tag); err != nil {
		return fmt.Errorf("push tag %q: %w", tag, err)
	}
	return nil
}
