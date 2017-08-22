package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

// AppCloudPlugin is the Swisscom Application Cloud cf CLI plugin
type AppCloudPlugin struct{}

// GetMetadata retrieves the metadata for the plugin
func (p *AppCloudPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Swisscom Application Cloud",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "create-backup",
				HelpText: "Create a backup for a service instance",
				UsageDetails: plugin.Usage{
					Usage: "create-backup SERVICE_INSTANCE",
				},
			},
			{
				Name:     "backups",
				HelpText: "List all backups for a service instance",
				UsageDetails: plugin.Usage{
					Usage: "backups SERVICE_INSTANCE",
				},
			},
		},
	}
}

// Run initiates the plugin
func (p *AppCloudPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case "create-backup":
		if len(args) < 2 {
			fmt.Println("Please provide a service instance name")
			return
		}

		p.CreateBackup(cliConnection, args[1])
	case "backups":
		if len(args) < 2 {
			fmt.Println("Please provide a service instance name")
			return
		}

		p.Backups(cliConnection, args[1])
	}
}

func main() {
	plugin.Start(new(AppCloudPlugin))
}
