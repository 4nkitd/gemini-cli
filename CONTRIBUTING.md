# Contributing to gemini-cli

A terminal AI assistant powered by Google Gemini.

> ⚠️ **Heads up: `main` does not currently build.** `llm.go` uses a `sapiens` API that does not exist in the pinned `sapiens v0.1.3`. See CHANGELOG → Known issues. Fixing this is the highest-priority contribution right now.

## Dev setup

Prereqs:

- Go 1.23+

```bash
git clone https://github.com/4nkitd/gemini-cli.git
cd gemini-cli
go mod download
go build -o gemini-cli .
```

(Once the sapiens API mismatch is fixed, the build above will succeed.)

## Environment

```bash
export GENAI_API_KEY=<your-google-ai-studio-key>
export GENAI_DEFAULT_MODEL=gemini-2.0-flash-exp   # or any current Gemini model
```

Get a key at [Google AI Studio](https://aistudio.google.com/).

## Running

```bash
./gemini-cli writer "Hey, can we meet to discuss the project?"
./gemini-cli cli "find files larger than 100MB"
./gemini-cli commit
./gemini-cli web
./gemini-cli assist "explain this error"
```

## Layout

```
gemini-cli/
├── main.go        # cobra root command
├── cli.go         # `cli` / `ask` subcommand
├── writer.go      # `writer` / `revise` subcommand
├── git.go         # `commit` subcommand
├── copilot.go     # `assist` / `copilot` subcommand (macOS/Windows only)
├── web.go         # `web` subcommand (local HTTP UI + /answer JSON API)
├── llm.go         # sapiens agent wrapper
├── storage.go     # SQLite-backed command history
├── system.go      # OS / hardware info collection
├── utils.go       # shared helpers
├── writer.go      # text-refinement subcommand
└── web/           # embedded web UI assets
```

## Branches and commits

- Branch from `main`: `feat/<name>`, `fix/<name>`, `docs/<name>`.
- Conventional Commits encouraged.

## PR checklist

- [ ] `go build .` succeeds
- [ ] `go vet ./...` clean
- [ ] Subcommand still works end-to-end with a real API key
- [ ] README updated if user-facing
- [ ] `CHANGELOG.md` entry under `## [Unreleased]`

## Releases

Maintainers only. Tag with `vX.Y.Z`, push the tag. goreleaser config exists; release runs are local (Actions disabled at account level for now).

## Reporting issues

[GitHub issues](https://github.com/4nkitd/gemini-cli/issues).
