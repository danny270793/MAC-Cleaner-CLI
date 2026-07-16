# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project

MAC Cleaner CLI — a Go command-line tool.

## Toolchain

- Go version is pinned via [.tool-versions](.tool-versions) and managed with `asdf`. See [docs/asdf.md](docs/asdf.md) for setup.
- Run `asdf install` from the repo root to install the pinned Go version.

## Common commands

- Run locally: `./scripts/start.sh`
- Build binary (outputs to `./build`): `./scripts/build.sh`

## Conventions

- Module path: `danny270793/maccleaner`.
- Keep `./build` out of version control (already in `.gitignore`).
- Commit related changes file by file with descriptive messages rather than one large commit, unless told otherwise.
