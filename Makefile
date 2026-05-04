.PHONY: build test lint clean cross all

build:
	go build -o git-meh .
	ln -sf git-meh gitmeh

test:
	go test ./... -count=1

lint:
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
