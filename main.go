package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/swisscom/appcloud-cf-cli-plugin/internal/appcloud"
)

func main() {
	plugin.Start(new(appcloud.Plugin))
}
