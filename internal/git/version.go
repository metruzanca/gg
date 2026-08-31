package git

import (
	"strconv"
	"strings"
)

func IncrementVersionAt(value string, pos, delta int) (string, int, bool) {
	if pos < 0 || pos > len(value) {
		return value, pos, false
	}

	digitIdx := -1
	if pos > 0 && isDigit(value[pos-1]) {
		digitIdx = pos - 1
	} else if pos < len(value) && isDigit(value[pos]) {
		digitIdx = pos
	}
	if digitIdx < 0 {
		return value, pos, false
	}

	start := digitIdx
	for start > 0 && isDigit(value[start-1]) {
		start--
	}
	end := digitIdx + 1
	for end < len(value) && isDigit(value[end]) {
		end++
	}

	n, err := strconv.Atoi(value[start:end])
	if err != nil {
		return value, pos, false
	}

	n += delta
	if n < 0 {
		n = 0
	}

	width := end - start
	newNum := strconv.Itoa(n)
	if width > 1 && value[start] == '0' {
		for len(newNum) < width {
			newNum = "0" + newNum
		}
	}

	value = value[:start] + newNum + value[end:]
	return value, start + len(newNum), true
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func LatestVersion(tags []string) string {
	best := ""
	var bestParts []int
	for _, t := range tags {
		parts, ok := parseVersion(t)
		if !ok {
			continue
		}
		if best == "" || compareVersionParts(parts, bestParts) > 0 {
			best = t
			bestParts = parts
		}
	}
	return best
}

func parseVersion(t string) ([]int, bool) {
	s := t
	if strings.HasPrefix(s, "v") {
		s = s[1:]
	}
	if s == "" {
		return nil, false
	}
	segs := strings.Split(s, ".")
	if len(segs) < 2 {
		return nil, false
	}
	var parts []int
	for _, seg := range segs {
		if seg == "" {
			return nil, false
		}
		for _, c := range seg {
			if c < '0' || c > '9' {
				return nil, false
			}
		}
		n, err := strconv.Atoi(seg)
		if err != nil {
			return nil, false
		}
		parts = append(parts, n)
	}
	return parts, true
}

func compareVersionParts(a, b []int) int {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		var av, bv int
		if i < len(a) {
			av = a[i]
		}
		if i < len(b) {
			bv = b[i]
		}
		if av != bv {
			if av > bv {
				return 1
			}
			return -1
		}
	}
	return 0
}