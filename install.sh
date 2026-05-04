#!/usr/bin/env bash
# Install git-meh from the latest GitHub release into ~/.local/bin.
set -euo pipefail

REPO="ryanhellyer/gitmeh"
dest_dir="${HOME}/.local/bin"
dest="${dest_dir}/git-meh"
marker='# git-meh PATH (added by install.sh)'
path_line='export PATH="${HOME}/.local/bin:${PATH}"'

select_artifact_name() {
	local os arch
	os=$(uname -s)
	arch=$(uname -m)
	case "${os}:${arch}" in
	Linux:x86_64)  echo "git-meh-linux-x86_64"  ;;
	Linux:aarch64|Linux:arm64) echo "git-meh-linux-arm64" ;;
	Darwin:x86_64) echo "git-meh-macos-x86_64"  ;;
	Darwin:arm64)  echo "git-meh-macos-arm64"   ;;
	*)
		echo "error: unsupported system (${os} ${arch})." >&2
		echo "       Supported: Linux x86_64 / arm64, macOS x86_64 / arm64." >&2
		exit 1
		;;
	esac
}

download() {
	local url=$1 dst=$2
	if command -v curl >/dev/null 2>&1; then
		curl -fsSL -o "${dst}" "${url}"
	elif command -v wget >/dev/null 2>&1; then
		wget -qO "${dst}" "${url}"
	else
		echo "error: need curl or wget to download the binary." >&2
		exit 1
	fi
}

verify_binary_kind() {
	local src=$1 artifact=$2
	if ! command -v file >/dev/null 2>&1; then
		return 0
	fi
	local desc
	desc=$(file -b "${src}" 2>/dev/null || true)
	if [[ -z "${desc}" ]]; then
		return 0
	fi
	case "${artifact}" in
	git-meh-linux-x86_64)
		echo "${desc}" | grep -qi 'ELF' || { echo "error: ${artifact} should be an ELF binary; file(1) says: ${desc}" >&2; exit 1; }
		echo "${desc}" | grep -qiE 'x86-64|x86_64' || { echo "error: ${artifact} should be x86-64; file(1) says: ${desc}" >&2; exit 1; }
		;;
	git-meh-linux-arm64)
		echo "${desc}" | grep -qi 'ELF' || { echo "error: ${artifact} should be an ELF binary; file(1) says: ${desc}" >&2; exit 1; }
		echo "${desc}" | grep -qiE 'aarch64|ARM aarch64|ARM, EABI64' || { echo "error: ${artifact} should be ARM aarch64; file(1) says: ${desc}" >&2; exit 1; }
		;;
	git-meh-macos-x86_64)
		echo "${desc}" | grep -qi 'Mach-O' || { echo "error: ${artifact} should be a Mach-O binary; file(1) says: ${desc}" >&2; exit 1; }
		echo "${desc}" | grep -qi 'x86_64' || { echo "error: ${artifact} should be x86_64; file(1) says: ${desc}" >&2; exit 1; }
		;;
	git-meh-macos-arm64)
		echo "${desc}" | grep -qi 'Mach-O' || { echo "error: ${artifact} should be a Mach-O binary; file(1) says: ${desc}" >&2; exit 1; }
		echo "${desc}" | grep -qi 'arm64' || { echo "error: ${artifact} should be arm64; file(1) says: ${desc}" >&2; exit 1; }
		;;
	esac
}

artifact=$(select_artifact_name)
tmp=$(mktemp)
trap 'rm -f "${tmp}"' EXIT

url="https://github.com/${REPO}/releases/latest/download/${artifact}"
echo "Downloading ${artifact} from GitHub releases ..."
download "${url}" "${tmp}"

verify_binary_kind "${tmp}" "${artifact}"

mkdir -p "${dest_dir}"
install -m 0755 "${tmp}" "${dest}"
ln -sf git-meh "${dest_dir}/gitmeh"

echo "Installed: ${dest}"
echo "Symlink:   ${dest_dir}/gitmeh -> git-meh"

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
		echo "Then: git meh   (or just: gitmeh)"
	else
		rc=$(choose_rc)
		touch "${rc}"
		printf '\n%s\n%s\n' "${marker}" "${path_line}" >>"${rc}"
		echo "Added ~/.local/bin to PATH in ${rc}"
		echo "Run:  source ${rc}"
		echo "Then: git meh   (or just: gitmeh)"
	fi
fi
