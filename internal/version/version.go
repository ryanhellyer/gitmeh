package version

// Version is the release version. Override at link time, for example:
//
//	go build -ldflags "-X gitmeh/internal/version.Version=1.0.0"
var Version = "3.0"
