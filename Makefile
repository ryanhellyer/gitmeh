.PHONY: dev build test lint clean cross all

dev:
	go build -ldflags="-X gitmeh/internal/config.isDev=true" -o git-meh .
	ln -sf git-meh gitmeh

build:
	go build -o git-meh .
	ln -sf git-meh gitmeh

test:
	go test ./... -count=1

lint:
	command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.1
	golangci-lint run ./...
	govulncheck ./...

clean:
	rm -f git-meh gitmeh git-meh.exe gitmeh.exe git-meh-linux-* git-meh-macos-*

cross: clean
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o git-meh-linux-x86_64      .
	CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 go build -o git-meh-linux-arm64       .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o git-meh-macos-x86_64      .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o git-meh-macos-arm64       .
	go build -o git-meh .

all: lint test cross
