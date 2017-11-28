# Swisscom Application Cloud cf CLI Plugin

A Plugin for the cf [CLI](https://github.com/cloudfoundry/cli) which extends it by the additional features of the [Swisscom Application Cloud](https://developer.swisscom.com).

## Installation

1. Install [Go](https://golang.org/)
1. `git clone` this repo into your `$GOPATH`
1. Run the [Development requirements](https://github.com/cloudfoundry/cli/tree/master/plugin/plugin_examples#development-requirements) commands for cf CLI plugins
1. Run `go get -d`
1. Run `go build`
1. Run `cf install-plugin appcloud-cf-cli-plugin`

## Commands

Simply run `cf` to see a list of the commands which are exposed by this plugin.
