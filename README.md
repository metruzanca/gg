# gg

A minimal git TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Inspect staged and unstaged changes, stage and unstage files, write
conventional commits, and manage tags with instant semantic-version bumps.

<p align="center">
  <img src="demos/tag-bump.gif" alt="Creating a tag and bumping its version in gg" width="1200">
</p>

## Features

- **Status** panel showing staged and unstaged changes
- **Staging** files and directories with `s`, expand collapsible directories
- **Commits** with conventional-commit types, titles, and descriptions
- **Tags** TUI (`gg tag`) that prefills the latest version and bumps it with the
  arrow keys

### Semantic version bumping

In the tag editor, pressing `↑`/`↓` bumps the number under the cursor and resets
the lower-order components to zero, so bumping the major clears the minor and
patch, and bumping the minor clears the patch:

```text
v0.1.0  ↑   v0.1.1   v0.1.2   v0.1.3      (patch)
        ←   v0.2.0   v0.3.0              (minor, patch resets)
        ←   v1.0.0   v2.0.0              (major, minor + patch reset)
        → + -beta    v2.0.0-beta
        ←   v2.0.1-beta                  (patch)
```

## Install

With Go installed:

```sh
go install github.com/metruzanca/gg@latest
```

Prebuilt binaries for Linux, macOS, and Windows are published with
[GoReleaser](https://goreleaser.com).

## Usage

```sh
gg        # open the git TUI
gg tag    # manage tags
```

## Demos

The demos are recorded with [VHS](https://github.com/charmbracelet/vhs), using a
patched fork build that fixes a startup race (see `demos/investigate/README.md`).
Replay them from the `demos/` directory:

```sh
demos/render.sh
```