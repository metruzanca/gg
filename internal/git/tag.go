package git

import (
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