# repos

[![Go Version](https://img.shields.io/github/go-mod/go-version/calvin180smith/repos)](https://go.dev/)
[![Latest Release](https://img.shields.io/github/v/release/calvin180smith/repos)](https://github.com/calvin180smith/repos/releases/latest)
[![CI](https://github.com/calvin180smith/repos/actions/workflows/ci.yaml/badge.svg)](https://github.com/calvin180smith/repos/actions/workflows/ci.yaml)

A CLI tool for discovering, listing, inspecting, and opening Git repositories on your local machine. Configure one or more directories to scan, and `repos` will find all Git repositories within them.

## Example Output

### `repos list`

```
+------------------+-------------------------------+---------------+
| NAME             | PATH                          | LAST MODIFIED |
+------------------+-------------------------------+---------------+
| repos            | /home/user/projects/repos     | 3 hours ago   |
| my-api           | /home/user/projects/my-api    | 2 days ago    |
| dotfiles         | /home/user/projects/dotfiles  | 5 days ago    |
+------------------+-------------------------------+---------------+
```

### `repos info`

```
Name:        my-api
Path:        /home/user/projects/my-api
Remote:      git@github.com:user/my-api.git
Branch:      main
Branches:    3
Last Commit: fix: update error handling (Thu May 21 10:32:00 2026 -0500) <user>)
```

## Installation

### Option A: Download a pre-built binary

1. Go to the [Releases page](https://github.com/calvin180smith/repos/releases/latest)
2. Download the archive for your OS and architecture
3. Extract the binary and place it somewhere in your `PATH`

**Linux / macOS:**

```bash
# Example for Linux amd64
tar -xzf repos_*_linux_amd64.tar.gz
sudo mv repos /usr/local/bin/
```

**Windows:**

Download the `.zip` file, extract `repos.exe`, and add its location to your `PATH`.

### Option B: Install with Go

```bash
go install github.com/calvin180smith/repos@latest
```

Requires Go 1.24+.

## Quick Start

```bash
# Add directories that contain your Git repositories
repos config set --path /home/user/projects

# List all discovered repositories
repos list
```

## Commands

### `repos config` — Manage scan directories

Settings are stored in `~/.repos.yaml`.

```bash
# Set directories to scan (replaces existing paths)
repos config set --path /path/to/dir1 --path /path/to/dir2

# Add a directory to the existing list
repos config add --path /path/to/dir3

# Display current configuration
repos config show
```

### `repos list` — List repositories

Lists all discovered Git repositories sorted by most recently modified.

```bash
repos list

# Limit the number of results
repos list --limit 10
```

### `repos info` — Show repository details

Displays detailed information about a repository including remote URL, current branch, number of local branches, and the last commit.

```bash
repos info my-repo
```

### `repos open` — Open a repository in VS Code

```bash
# Open a specific repository
repos open my-repo

# Open the most recently modified repository
repos open --latest
```

## Configuration

The configuration file is stored at `~/.repos.yaml` with the following format:

```yaml
paths:
  - /home/user/projects
  - /home/user/work
```

Duplicate paths are automatically prevented when using `config add`.
