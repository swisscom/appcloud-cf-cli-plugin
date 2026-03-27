.PHONY: all build clean

all: clean build

build:
	@echo "Did you update the version in main.go/PluginMetadata?"
	CGO_ENABLED=0 GOOS=linux   GOARCH=386   go build -o bin/appcloud-cf-cli-plugin_linux32
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o bin/appcloud-cf-cli-plugin_linux64
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o bin/appcloud-cf-cli-plugin_osx
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o bin/appcloud-cf-cli-plugin_osx_arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=386   go build -o bin/appcloud-cf-cli-plugin_win32.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/appcloud-cf-cli-plugin_win64.exe

clean:
	@rm bin/appcloud-cf-cli-plugin_* -f

sha:
	@sha1sum bin/appcloud-cf-cli-plugin_*

release:
	@echo release plugin - see https://github.com/swisscom/appcloud-cf-cli-plugin/releases
	@echo open plugin repo PR - see https://github.com/cloudfoundry/cli-plugin-repo/pull/519
