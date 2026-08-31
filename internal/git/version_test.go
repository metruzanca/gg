package git

import "testing"

func TestIncrementVersionAt(t *testing.T) {
	cases := []struct {
		value string
		pos   int
		delta int
		want  string
	}{
		{"v0.1.0", 3, 1, "v0.2.0"},
		{"v0.1.0", 4, 1, "v0.2.0"},
		{"v0.9.0", 3, 1, "v0.10.0"},
		{"v0.1.0", 2, 1, "v1.1.0"},
		{"v0.1.0", 1, 1, "v1.1.0"},
		{"v0.1.0", 5, 1, "v0.1.1"},
		{"v0.9.0", 5, 1, "v0.9.1"},
		{"v0.09.0", 3, 1, "v0.10.0"},
		{"v0.01.0", 3, 1, "v0.02.0"},
		{"v0.1.0", 4, -1, "v0.0.0"},
		{"v0.1.0", 3, -1, "v0.0.0"},
		{"v0.2.0", 3, -1, "v0.1.0"},
		{"v0.10.0", 3, -1, "v0.9.0"},
	}
	for _, c := range cases {
		got, _, ok := IncrementVersionAt(c.value, c.pos, c.delta)
		if !ok {
			t.Errorf("IncrementVersionAt(%q, %d, %d): unexpectedly failed", c.value, c.pos, c.delta)
			continue
		}
		if got != c.want {
			t.Errorf("IncrementVersionAt(%q, %d, %d) = %q, want %q", c.value, c.pos, c.delta, got, c.want)
		}
	}
}

func TestIncrementVersionAtFailsSilently(t *testing.T) {
	const value = "v0.1.0"
	for _, p := range []int{0} {
		got, _, ok := IncrementVersionAt(value, p, 1)
		if ok {
			t.Errorf("IncrementVersionAt(%q, %d, 1): expected failure", value, p)
		}
		if got != value {
			t.Errorf("IncrementVersionAt(%q, %d, 1) = %q, want value unchanged", value, p, got)
		}
	}
	cursorAtEnd, _, ok := IncrementVersionAt(value, len(value), 1)
	if !ok || cursorAtEnd != "v0.1.1" {
		t.Errorf("cursor at end: got %q ok=%v, want v0.1.1", cursorAtEnd, ok)
	}

	nonNumeric := []string{"", "abc", "v..0"}
	badPos := []int{1, 1, 0}
	for i, v := range nonNumeric {
		p := badPos[i]
		got, _, ok := IncrementVersionAt(v, p, 1)
		if ok {
			t.Errorf("IncrementVersionAt(%q, %d, 1): expected failure", v, p)
		}
		if got != v {
			t.Errorf("IncrementVersionAt(%q, %d, 1): value changed to %q", v, p, got)
		}
	}
}

func TestIncrementVersionAtCursorSticks(t *testing.T) {
	got, pos, ok := IncrementVersionAt("v0.9.0", 4, 1)
	if !ok || got != "v0.10.0" {
		t.Fatalf("first increment: got %q ok=%v", got, ok)
	}
	got, pos, ok = IncrementVersionAt(got, pos, 1)
	if !ok || got != "v0.11.0" {
		t.Fatalf("second increment from pos=%d: got %q ok=%v", pos, got, ok)
	}
}

func TestLatestVersion(t *testing.T) {
	cases := []struct {
		tags []string
		want string
	}{
		{[]string{"v1.9.0", "v1.10.0", "v1.2.3"}, "v1.10.0"},
		{[]string{"1.9.0", "1.10.0"}, "1.10.0"},
		{[]string{"v0.1.0", "v0.1.9", "v0.1.10"}, "v0.1.10"},
		{[]string{"v1.2.3", "release", "wip"}, "v1.2.3"},
		{[]string{}, ""},
		{[]string{"release", "foo"}, ""},
	}
	for _, c := range cases {
		if got := LatestVersion(c.tags); got != c.want {
			t.Errorf("LatestVersion(%v) = %q, want %q", c.tags, got, c.want)
		}
	}
}