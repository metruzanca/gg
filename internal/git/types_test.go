package git

import (
	"reflect"
	"testing"
)

func TestBuildTree(t *testing.T) {
	files := []FileStatus{
		{Path: "internal/git/stage.go", X: 'M', Y: 'M'},
		{Path: "cmd/root.go", X: 'A', Y: '.'},
		{Path: "main.go", X: '.', Y: '?'},
		{Path: "go.mod", X: '.', Y: '?'},
	}

	staged, unstaged := BuildTree(files)

	if len(staged) != 2 {
		t.Fatalf("expected 2 staged roots, got %d", len(staged))
	}
	if len(unstaged) != 3 {
		t.Fatalf("expected 3 unstaged roots (stage.go, main.go, go.mod), got %d", len(unstaged))
	}

	internal := staged[0]
	if internal.Name != "internal" {
		t.Errorf("expected 'internal', got '%s'", internal.Name)
	}
	if !internal.IsDir {
		t.Error("expected internal to be a directory")
	}
	if internal.Path != "internal" {
		t.Errorf("expected path 'internal', got '%s'", internal.Path)
	}

	cmd := staged[1]
	if cmd.Name != "cmd" {
		t.Errorf("expected 'cmd', got '%s'", cmd.Name)
	}
	if !cmd.IsDir {
		t.Error("expected cmd to be a directory")
	}
	if cmd.Path != "cmd" {
		t.Errorf("expected path 'cmd', got '%s'", cmd.Path)
	}

	if len(internal.Children) != 1 {
		t.Fatalf("expected 1 child of internal, got %d", len(internal.Children))
	}
	gitDir := internal.Children[0]
	if gitDir.Name != "git" {
		t.Errorf("expected 'git', got '%s'", gitDir.Name)
	}
	if !gitDir.IsDir {
		t.Error("expected git to be a directory")
	}
	if gitDir.Path != "internal/git" {
		t.Errorf("expected path 'internal/git', got '%s'", gitDir.Path)
	}

	if len(gitDir.Children) != 1 {
		t.Fatalf("expected 1 child of git, got %d", len(gitDir.Children))
	}
	stageGo := gitDir.Children[0]
	if stageGo.Name != "stage.go" {
		t.Errorf("expected 'stage.go', got '%s'", stageGo.Name)
	}
	if stageGo.IsDir {
		t.Error("expected stage.go not to be a directory")
	}
	if stageGo.Path != "internal/git/stage.go" {
		t.Errorf("expected full path, got '%s'", stageGo.Path)
	}
	if stageGo.Status != 'M' {
		t.Errorf("expected status 'M', got '%c'", stageGo.Status)
	}
	if stageGo.Section != SectionStaged {
		t.Error("expected staged section")
	}

	mainGo := findEntry(unstaged, "main.go")
	if mainGo.Name != "main.go" {
		t.Errorf("expected 'main.go', got '%s'", mainGo.Name)
	}
	if mainGo.Path != "main.go" {
		t.Errorf("expected path 'main.go', got '%s'", mainGo.Path)
	}

	directoriesHaveAllPaths(t, staged)
	directoriesHaveAllPaths(t, unstaged)
}

func directoriesHaveAllPaths(t *testing.T, entries []*Entry) {
	for _, e := range entries {
		if e.IsDir && e.Path == "" {
			t.Errorf("directory '%s' has empty path, depth=%d", e.Name, e.Depth)
		}
		directoriesHaveAllPaths(t, e.Children)
	}
}

func TestFlattenVisibleCollapsed(t *testing.T) {
	files := []FileStatus{
		{Path: "internal/git/stage.go", X: 'M', Y: '.'},
		{Path: "cmd/root.go", X: 'A', Y: '.'},
	}
	staged, _ := BuildTree(files)

	expanded := map[string]bool{}
	visible := FlattenVisible(staged, nil, expanded, false)

	if len(visible) != 2 {
		t.Fatalf("expected 2 visible with all collapsed, got %d", len(visible))
	}
	if !visible[0].IsDir || visible[0].Name != "internal" {
		t.Error("expected internal directory")
	}
	if !visible[1].IsDir || visible[1].Name != "cmd" {
		t.Error("expected cmd directory")
	}
}

func TestFlattenVisibleExpanded(t *testing.T) {
	files := []FileStatus{
		{Path: "internal/git/stage.go", X: 'M', Y: '.'},
		{Path: "cmd/root.go", X: 'A', Y: '.'},
	}
	staged, _ := BuildTree(files)

	expanded := map[string]bool{"internal": true, "internal/git": true}
	visible := FlattenVisible(staged, nil, expanded, false)

	if len(visible) != 4 {
		t.Fatalf("expected 4 visible (internal, git, stage.go, cmd), got %d", len(visible))
	}

	names := make([]string, len(visible))
	for i, v := range visible {
		names[i] = v.Name
	}
	expected := []string{"internal", "git", "stage.go", "cmd"}
	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected %v, got %v", expected, names)
	}
}

func TestFlattenVisibleExpandAll(t *testing.T) {
	files := []FileStatus{
		{Path: "internal/git/stage.go", X: 'M', Y: '.'},
		{Path: "cmd/root.go", X: 'A', Y: '.'},
	}
	staged, _ := BuildTree(files)

	visible := FlattenVisible(staged, nil, nil, true)

	if len(visible) != 5 {
		t.Fatalf("expected 5 visible with all expanded, got %d", len(visible))
	}
}

func TestUnstagedFiles(t *testing.T) {
	files := []FileStatus{
		{Path: "a.go", X: '.', Y: '?'},
		{Path: "b.go", X: '.', Y: 'M'},
		{Path: "pkg/c.go", X: 'A', Y: '.'},
	}
	_, unstaged := BuildTree(files)

	paths := CollectFiles(unstaged)
	if len(paths) != 2 {
		t.Fatalf("expected 2 unstaged files, got %d", len(paths))
	}
}

func TestHasEntry(t *testing.T) {
	files := []FileStatus{
		{Path: "a.go", X: 'M', Y: '.'},
	}
	staged, _ := BuildTree(files)

	if !HasEntry(staged, "a.go") {
		t.Error("expected staged to have a.go")
	}
	if HasEntry(staged, "b.go") {
		t.Error("expected staged not to have b.go")
	}
}

func findEntry(entries []*Entry, name string) *Entry {
	for _, e := range entries {
		if e.Name == name {
			return e
		}
	}
	return nil
}

func TestAllDirPaths(t *testing.T) {
	files := []FileStatus{
		{Path: "a/b/c.go", X: 'M', Y: '.'},
	}
	staged, _ := BuildTree(files)

	paths := AllDirPaths(staged)
	expected := []string{"a", "a/b"}
	if !reflect.DeepEqual(paths, expected) {
		t.Errorf("expected %v, got %v", expected, paths)
	}
}
