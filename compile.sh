#!/usr/bin/env bash
# Cross-build git-meh for Linux and macOS (amd64 + arm64 each).
# Output names are explicit so install scripts / users can pick the right file.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"
export CGO_ENABLED=0

build() {
	local goos=$1
	local goarch=$2
	local out=$3
	echo "==> ${out}  (GOOS=${goos} GOARCH=${goarch})"
	GOOS="${goos}" GOARCH="${goarch}" go build -o "${out}" .
}

# Linux (servers, desktops, most VMs)
build linux amd64 git-meh-linux-x86_64
build linux arm64 git-meh-linux-arm64

# macOS (Intel vs Apple Silicon)
build darwin amd64 git-meh-macos-x86_64
build darwin arm64 git-meh-macos-arm64

# Native binary for this machine (install.sh expects "git-meh" here until it learns the names above)
echo "==> git-meh  (native, current OS/arch)"
go build -o git-meh .

echo
echo "Done. Artifacts:"
ls -la git-meh git-meh-linux-x86_64 git-meh-linux-arm64 git-meh-macos-x86_64 git-meh-macos-arm64 2>/dev/null || true
