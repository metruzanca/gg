package git

import (
	"os/exec"
)

func DiffCached() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--stat")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func StagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	raw := string(out)
	if raw == "" {
		return nil, nil
	}
	lines := make([]string, 0)
	for _, l := range splitLines(raw) {
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
