package git

import (
	"os/exec"
	"strings"
)

func Status() ([]FileStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain=v2", "--branch")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePorcelainV2(string(out)), nil
}

func parsePorcelainV2(out string) []FileStatus {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	var files []FileStatus
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "? ") {
			path := strings.TrimPrefix(line, "? ")
			files = append(files, FileStatus{Path: path, X: '?', Y: '?'})
			continue
		}
		if strings.HasPrefix(line, "! ") {
			path := strings.TrimPrefix(line, "! ")
			files = append(files, FileStatus{Path: path, X: '!', Y: '!'})
			continue
		}
		if strings.HasPrefix(line, "1 ") || strings.HasPrefix(line, "2 ") {
			parts := strings.SplitN(line, " ", 9)
			if len(parts) < 9 {
				continue
			}
			xy := parts[1]
			if len(xy) < 2 {
				continue
			}
			x := xy[0]
			y := xy[1]
			path := parts[8]
			files = append(files, FileStatus{Path: path, X: x, Y: y})
		}
	}
	return files
}

func (fs FileStatus) HasStaged() bool {
	return fs.X != '.' && fs.X != '?' && fs.X != '!'
}

func (fs FileStatus) HasUnstaged() bool {
	return fs.Y != '.'
}

func StatusLabel(b byte) string {
	switch b {
	case 'M':
		return "M  modified"
	case 'A':
		return "A  added"
	case 'D':
		return "D  deleted"
	case 'R':
		return "R  renamed"
	case 'C':
		return "C  copied"
	case '?':
		return "?  untracked"
	case '!':
		return "!  ignored"
	case 'T':
		return "T  typechange"
	case 'U':
		return "U  unmerged"
	default:
		return string(b)
	}
}
