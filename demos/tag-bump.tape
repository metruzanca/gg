# Demo: bump a tag version with semantic-version resets.
# Repo pre-seeded with tag v0.1.0. Steam: patch -> minor -> major,
# then append -beta and bump the patch once more.

Output demos/tag-bump.gif

Set Width 860
Set Height 520
Set FontSize 20
Set Padding 60
Set Theme "Dracula"
Set BorderRadius 2

Sleep 300ms

Type "cd /tmp/opencode/vhs-demo"
Enter
Sleep 400ms
Type "./gg tag"
Enter
Sleep 1500ms

# ---------- patch: v0.1.0 -> v0.1.3 ----------
Type "n"
Sleep 800ms
Up
Sleep 500ms
Up
Sleep 500ms
Up
Sleep 700ms

# ---------- minor: v0.1.3 -> v0.3.0 ----------
Left
Sleep 150ms
Left
Sleep 150ms
Up
Sleep 500ms
Up
Sleep 700ms

# ---------- major: v0.3.0 -> v2.0.0 ----------
Left
Sleep 150ms
Left
Sleep 150ms
Up
Sleep 500ms
Up
Sleep 700ms

# ---------- append -beta ----------
Right
Sleep 150ms
Right
Sleep 150ms
Right
Sleep 150ms
Right
Sleep 150ms
Type "-beta"
Sleep 700ms

# ---------- bump patch: v2.0.0-beta -> v2.0.1-beta ----------
Left
Sleep 150ms
Left
Sleep 150ms
Left
Sleep 150ms
Left
Sleep 150ms
Left
Sleep 150ms
Up
Sleep 900ms

# ---------- create the tag, skip push ----------
Enter
Sleep 800ms
Type "n"
Sleep 1500ms