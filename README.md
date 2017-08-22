# Swisscom Application Cloud cf CLI Plugin

A Plugin for the cf CLI which extends it by the additional features of the Swisscom Application Cloud.

# Installation

1. Install [Go](https://golang.org/)
1. `git clone` this repo into your `$GOPATH`
1. Run the [Development requirements](https://github.com/cloudfoundry/cli/tree/master/plugin/plugin_examples#development-requirements) commands for cf CLI plugins
1. Run `go build`
1. Run `cf install-plugin appcloud-cf-cli-plugin`

# Commands

* `cf backups` Lists all backups for a service instance
* `cf create-backup` Creates a backup for a service instance
