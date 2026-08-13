package git

import "strings"

type FileStatus struct {
	Path string
	X    byte
	Y    byte
}

type Section int

const (
	SectionStaged Section = iota
	SectionUnstaged
)

type Entry struct {
	Name     string   // file or directory name (no path)
	Path     string   // full relative path
	IsDir    bool
	Depth    int
	Status   byte      // 0 for dirs
	Section  Section
	Children []*Entry
}

type CommitSpec struct {
	Type        string
	Title       string
	Description string
}

func (e *Entry) findChild(name string) *Entry {
	for _, c := range e.Children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func BuildTree(files []FileStatus) (staged, unstaged []*Entry) {
	sRoot := &Entry{Name: "", IsDir: true, Depth: -1}
	uRoot := &Entry{Name: "", IsDir: true, Depth: -1}

	for _, f := range files {
		parts := strings.Split(f.Path, "/")

		if f.HasStaged() {
			_ = newEntryFrom(parts, sRoot, SectionStaged, f.X, f.Path)
		}
		if f.HasUnstaged() {
			_ = newEntryFrom(parts, uRoot, SectionUnstaged, f.Y, f.Path)
		}
	}

	return sRoot.Children, uRoot.Children
}

func newEntryFrom(parts []string, parent *Entry, section Section, status byte, fullPath string) *Entry {
	current := parent
	for i, part := range parts {
		isLast := i == len(parts)-1
		child := current.findChild(part)
		if child == nil {
			child = &Entry{Name: part, Depth: current.Depth + 1}
			current.Children = append(current.Children, child)
		}
		if isLast {
			child.Path = fullPath
			child.Status = status
			child.Section = section
		} else {
			child.IsDir = true
			child.Path = strings.Join(parts[:i+1], "/")
			child.Section = section
		}
		current = child
	}
	return current
}

func FlattenVisible(staged, unstaged []*Entry, expandedDirs map[string]bool, allExpanded bool) []*Entry {
	var visible []*Entry
	for _, e := range staged {
		flattenOne(e, &visible, expandedDirs, allExpanded)
	}
	for _, e := range unstaged {
		flattenOne(e, &visible, expandedDirs, allExpanded)
	}
	return visible
}

func FlattenVisibleForPanel(entries []*Entry, expandedDirs map[string]bool, allExpanded bool) []*Entry {
	var visible []*Entry
	for _, e := range entries {
		flattenOne(e, &visible, expandedDirs, allExpanded)
	}
	return visible
}

func flattenOne(e *Entry, out *[]*Entry, expanded map[string]bool, all bool) {
	*out = append(*out, e)
	if e.IsDir && (all || expanded[e.Path]) {
		for _, c := range e.Children {
			flattenOne(c, out, expanded, all)
		}
	}
}

func CollectFiles(entries []*Entry) []string {
	var paths []string
	walkFiles(entries, &paths)
	return paths
}

func walkFiles(entries []*Entry, paths *[]string) {
	for _, e := range entries {
		if !e.IsDir {
			*paths = append(*paths, e.Path)
		}
		walkFiles(e.Children, paths)
	}
}

func HasEntry(entries []*Entry, path string) bool {
	for _, e := range entries {
		if !e.IsDir && e.Path == path {
			return true
		}
		if HasEntry(e.Children, path) {
			return true
		}
	}
	return false
}

func AllDirPaths(entries []*Entry) []string {
	var paths []string
	collectDirs(entries, &paths)
	return paths
}

func collectDirs(entries []*Entry, paths *[]string) {
	for _, e := range entries {
		if e.IsDir {
			*paths = append(*paths, e.Path)
			collectDirs(e.Children, paths)
		}
	}
}
