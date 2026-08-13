package git

import "testing"

func TestParsePorcelainV2(t *testing.T) {
	out := `# branch.oid abc
# branch.head main
1 M. N... 100644 100644 100644 abc123 def456 modified_staged.go
1 .M N... 100644 100644 100644 abc123 def456 modified_unstaged.go
1 A. N... 000000 100644 100644 000000 abc123 added_staged.go
1 MM N... 100644 100644 100644 abc123 def456 both.go
? untracked.go
! ignored.go
`

	files := parsePorcelainV2(out)

	got := map[string]FileStatus{}
	for _, f := range files {
		got[f.Path] = f
	}

	if len(files) != 6 {
		t.Fatalf("expected 6 files, got %d: %+v", len(files), files)
	}

	cases := []struct {
		path     string
		x, y     byte
		staged   bool
		unstaged bool
	}{
		{"modified_staged.go", 'M', '.', true, false},
		{"modified_unstaged.go", '.', 'M', false, true},
		{"added_staged.go", 'A', '.', true, false},
		{"both.go", 'M', 'M', true, true},
		{"untracked.go", '?', '?', false, true},
		{"ignored.go", '!', '!', false, true},
	}

	for _, c := range cases {
		fs, ok := got[c.path]
		if !ok {
			t.Errorf("missing file %q", c.path)
			continue
		}
		if fs.X != c.x || fs.Y != c.y {
			t.Errorf("%q: expected X=%c Y=%c, got X=%c Y=%c", c.path, c.x, c.y, fs.X, fs.Y)
		}
		if fs.HasStaged() != c.staged {
			t.Errorf("%q: HasStaged() = %v, want %v", c.path, fs.HasStaged(), c.staged)
		}
		if fs.HasUnstaged() != c.unstaged {
			t.Errorf("%q: HasUnstaged() = %v, want %v", c.path, fs.HasUnstaged(), c.unstaged)
		}
	}
}
