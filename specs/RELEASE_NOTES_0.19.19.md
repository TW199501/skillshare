# skillshare v0.19.19 Release Notes

Release date: 2026-05-23

## TL;DR

1. **Tracked root-level `SKILL.md` repos now get integrity hashes** — `doctor` can verify repos installed with `--track` when the skill lives at the repository root
2. **Existing v0.19.18 installs can be backfilled with `update`** — run `skillshare update _repo` to populate missing hashes for already-installed tracked root-skill repos

This is a patch release — bug fixes only, no new commands or flags.

---

## Integrity Hashes for Root-Level Tracked Skills

v0.19.18 fixed tracked repos whose `SKILL.md` lives at the repository root so they install, sync, and show up in status correctly. One metadata gap remained: those installs could still miss the `file_hashes` block that `skillshare doctor` uses for integrity verification.

v0.19.19 closes that gap. New tracked installs now write integrity hashes for root-level skills:

```bash
skillshare install github.com/team/my-skill --track
skillshare doctor
# Skill integrity verifies normally instead of warning about missing file hashes
```

If you already installed a root-level tracked repo with v0.19.18 and `doctor` reports that it is missing file hashes, update that repo once:

```bash
skillshare update _my-skill
skillshare doctor
```

The update command backfills the missing hashes even when the repository is already up to date. No reinstall is required. Refs: #165
