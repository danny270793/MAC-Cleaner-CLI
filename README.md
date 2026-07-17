# 🧹 MAC Cleaner CLI

A small, fast command-line tool that frees up disk space on macOS by clearing caches you don't need — Gradle, `~/Library/Caches`, `~/.pub-cache`, outdated VS Code extension versions, Xcode DerivedData, CoreSimulator caches, the Go module cache, Cargo's registry cache, npm's cache, the pnpm store, and Docker.

> **macOS only.** This tool targets macOS-specific paths (`~/Library/Caches`) and behavior, and is not intended to run on Linux or Windows.

[![CI](https://github.com/danny270793/MAC-Cleaner-CLI/actions/workflows/ci.yml/badge.svg)](https://github.com/danny270793/MAC-Cleaner-CLI/actions/workflows/ci.yml)

## Why

Development tools quietly pile up gigabytes of caches over time: old Gradle downloads, stale VS Code extension versions, dangling Docker images. `maccleaner` clears them out, one flag at a time, and tells you how much space each step will free before it runs.

## Install

### Prebuilt binary

Downloads the latest release for your Mac's architecture and installs it to `~/.local/bin/maccleaner`:

```sh
curl -fsSL https://raw.githubusercontent.com/danny270793/MAC-Cleaner-CLI/main/scripts/install.sh | bash
```

Make sure `~/.local/bin` is on your `PATH` — the script will tell you if it isn't.

### From source

Requires Go, managed via [asdf](https://asdf-vm.com/) — see [docs/asdf.md](docs/asdf.md) for setup.

```sh
git clone git@github.com:danny270793/MAC-Cleaner-CLI.git
cd MAC-Cleaner-CLI
asdf install
```

## Usage

Run without building, straight from source:

```sh
./scripts/start.sh --all
```

Or build a binary first:

```sh
./scripts/build.sh
./build/maccleaner --all
```

### Flags

| Flag                     | Description                                                                  |
| ------------------------ | ----------------------------------------------------------------------------- |
| `--all`                   | Run every cleaner                                                            |
| `--docker`                | Run `docker system prune`                                                    |
| `--gradle`                | Clear `~/.gradle/caches` and `~/.gradle/wrapper/dists`                       |
| `--library-caches`        | Clear the contents of `~/Library/Caches`                                     |
| `--pub-cache`             | Clear the contents of `~/.pub-cache`                                         |
| `--vscode-extensions`     | Remove outdated versions of installed VS Code extensions, keeping the latest |
| `--xcode-derived-data`    | Clear `~/Library/Developer/Xcode/DerivedData`                                |
| `--core-simulator-caches` | Clear `~/Library/Developer/CoreSimulator/Caches`                             |
| `--go-mod-cache`          | Clear the Go module cache via `go clean -modcache`                          |
| `--cargo-cache`           | Clear `~/.cargo/registry/cache` and `~/.cargo/registry/src`                  |
| `--npm-cache`             | Clear the contents of `~/.npm`                                               |
| `--pnpm-store`            | Prune the pnpm store via `pnpm store prune`                                 |
| `--dry-run`               | Show what would be cleaned (with sizes) without actually cleaning anything  |
| `--auto-approve`          | Skip the confirmation prompt shown before each cleaner                      |
| `--version`               | Print the version and exit                                                   |
| `--help`                  | Show usage and exit                                                          |

At least one flag is required — running with none prints usage instead of doing nothing silently.

### Example

```
$ ./build/maccleaner --gradle --docker
[ ] cleaning gradle (1.2GB)
cleaned /Users/you/.gradle/caches
cleaned /Users/you/.gradle/wrapper/dists
[x] cleaned gradle

[ ] cleaning docker
WARNING! This will remove:
  - all stopped containers
  - all networks not used by at least one container
  - all dangling images
  - unused build cache

Are you sure you want to continue? [y/N] y
[x] cleaned docker
```

Cleaners that clear a cache folder (Gradle, Library Caches, Pub Cache, Xcode DerivedData, CoreSimulator Caches, Cargo Cache, npm Cache) remove the folder's *contents* only — the folder itself is left in place. The Go module cache and pnpm store are cleared via their own tools (`go clean -modcache`, `pnpm store prune`) instead of a raw delete. Docker's reclaimable space can't be measured up front, so no estimate is shown for it.

## Development

```sh
./scripts/build.sh    # build the binary into ./build
./scripts/start.sh    # run from source, arguments are passed through
./scripts/test.sh     # run the test suite
./scripts/format.sh   # format the source with gofmt
```

See [AGENTS.md](AGENTS.md) for contribution conventions (Conventional Commits, branch naming, etc.) — the single source of truth for both human and AI contributors.
