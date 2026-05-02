#!/usr/bin/env bash
# Install git-meh next to this script into ~/.local/bin as git-meh (for "git meh").
set -euo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
src="${script_dir}/git-meh"
dest_dir="${HOME}/.local/bin"
dest="${dest_dir}/git-meh"

if [[ ! -f "${src}" ]]; then
	echo "error: git-meh not found next to install.sh (expected: ${src})" >&2
	exit 1
fi

mkdir -p "${dest_dir}"
install -m 0755 "${src}" "${dest}"

echo "Installed: ${dest}"

if command -v git-meh >/dev/null 2>&1; then
	echo "On PATH as: $(command -v git-meh)"
else
	echo
	echo "~/.local/bin is not on your PATH. Add something like this to ~/.profile or ~/.bashrc:"
	echo "  export PATH=\"\${HOME}/.local/bin:\${PATH}\""
	echo "Then open a new shell or run: source ~/.profile"
fi
echo
echo "Usage: git meh"
