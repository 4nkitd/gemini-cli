# Changelog

All notable changes to this project are documented here.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versioning: [SemVer](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `CHANGELOG.md` and `CONTRIBUTING.md`

### Changed
- README naming unified as `gemini-cli` (previously mixed Gema CLI / gema / gema-cli)
- Added disambiguation note: not affiliated with Google's `gemini-cli`

### Known issues
- `main` does not currently build: `llm.go` references a `sapiens` API (`sapiens.Tool`, `agent.AddTools`, `RegisterToolImplementation`, `AddImageContent`, etc.) that does not exist in the pinned `sapiens v0.1.3` nor in any tagged sapiens version. Tracked in `PLAN.md`. Fix requires either migrating to the current sapiens API or pinning to a working commit.

---

## Pre-1.0 / pre-cleanup tags

Earlier releases were tagged inconsistently (`0.4.x`, `0.5.x`, `5`, `5.0.x`, plus `v0.0.x`). The next tagged release will be `v0.1.0` after the sapiens API mismatch above is resolved.

Highlights of prior work, from git history:

### Web mode
- `web` command starts a local HTTP server with chat UI + JSON API at `POST /answer`
- Web assets embedded into the binary

### CoPilot
- `assist` / `copilot` / `assistant` subcommands for screen-content explanation
- Optional text-to-speech of responses
- Conditional inclusion on macOS / Windows only

### Git
- `commit` subcommand generates AI commit messages from staged diff
- Optional repo path argument and custom `--prompt`

### CLI
- `cli` / `ask` for natural-language → command suggestion
- Confirmation before execution

### Writer
- `writer` / `revise` for text refinement
- `[length=N]` and `[type=email]` formatting modifiers

### Quality of life
- Colored output (`fatih/color`)
- Loading animation
- SQLite-backed command history
- System info captured via `sapiens` agent helper

[Unreleased]: https://github.com/4nkitd/gemini-cli/compare/HEAD...HEAD
