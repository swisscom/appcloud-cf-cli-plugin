.PHONY: all build clean

build:
	GOOS=linux   GOARCH=386   go build -o bin/appcloud-cf-cli-plugin_linux32
	GOOS=linux   GOARCH=amd64 go build -o bin/appcloud-cf-cli-plugin_linux64
	GOOS=darwin  GOARCH=amd64 go build -o bin/appcloud-cf-cli-plugin_osx
	GOOS=darwin  GOARCH=arm64 go build -o bin/appcloud-cf-cli-plugin_osx_arm64
	GOOS=windows GOARCH=386   go build -o bin/appcloud-cf-cli-plugin_win32.exe
	GOOS=windows GOARCH=amd64 go build -o bin/appcloud-cf-cli-plugin_win64.exe

clean:
	@rm bin/appcloud-cf-cli-plugin_* -f

all: clean build
