#!/usr/bin/env bash
# Install git-meh next to this script into ~/.local/bin as git-meh (for "git meh").
set -euo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
src="${script_dir}/git-meh"
dest_dir="${HOME}/.local/bin"
dest="${dest_dir}/git-meh"
marker='# git-meh PATH (added by install.sh)'
path_line='export PATH="${HOME}/.local/bin:${PATH}"'

if [[ ! -f "${src}" ]]; then
	echo "error: git-meh not found next to install.sh (expected: ${src})" >&2
	exit 1
fi

mkdir -p "${dest_dir}"
install -m 0755 "${src}" "${dest}"

echo "Installed: ${dest}"

path_has_local_bin() {
	case ":${PATH}:" in
	*:"${HOME}/.local/bin":*) return 0 ;;
	*) return 1 ;;
	esac
}

already_marked() {
	local f
	for f in "${HOME}/.zshrc" "${HOME}/.bashrc" "${HOME}/.bash_profile" "${HOME}/.profile"; do
		[[ -f "${f}" ]] || continue
		if grep -qF "${marker}" "${f}" 2>/dev/null; then
			echo "${f}"
			return 0
		fi
	done
	return 1
}

choose_rc() {
	local shell_base
	shell_base=$(basename "${SHELL:-bash}" 2>/dev/null || echo bash)
	if [[ "${shell_base}" == zsh ]]; then
		echo "${HOME}/.zshrc"
		return
	fi
	if [[ "${shell_base}" == bash ]]; then
		if [[ "$(uname -s)" == Darwin ]] && [[ -f "${HOME}/.bash_profile" ]]; then
			echo "${HOME}/.bash_profile"
			return
		fi
		echo "${HOME}/.bashrc"
		return
	fi
	echo "${HOME}/.profile"
}

hash -r 2>/dev/null || true

if command -v git-meh >/dev/null 2>&1; then
	echo "On PATH as: $(command -v git-meh)"
elif path_has_local_bin; then
	echo "Open a new terminal (or run: hash -r) so your shell picks up git-meh."
else
	existing=""
	if existing=$(already_marked); then
		echo "Run:  source ${existing}"
		echo "Then: git meh"
	else
		rc=$(choose_rc)
		touch "${rc}"
		printf '\n%s\n%s\n' "${marker}" "${path_line}" >>"${rc}"
		echo "Added ~/.local/bin to PATH in ${rc}"
		echo "Run:  source ${rc}"
		echo "Then: git meh"
	fi
fi
