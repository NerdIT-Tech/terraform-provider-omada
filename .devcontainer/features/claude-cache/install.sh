#!/usr/bin/env bash
set -euo pipefail

# _REMOTE_USER / _REMOTE_USER_HOME are injected by the devcontainer CLI
# during feature installation (run as root).
target_dir="${_REMOTE_USER_HOME:-/home/vscode}/.claude"

mkdir -p "$target_dir"

if [ -n "${_REMOTE_USER:-}" ]; then
  chown -R "${_REMOTE_USER}:${_REMOTE_USER}" "$target_dir"
fi
