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
			{
				Name:     "invite-space-user",
				HelpText: "Invite a user to a space",
				UsageDetails: plugin.Usage{
					Usage: "invite-space-user SPACE_NAME INVITEE ROLE1(,ROLE2(,ROLE3))",
				},
			},
			{
				Name:     "invite-org-user",
				HelpText: "Invite a user to an organization",
				UsageDetails: plugin.Usage{
					Usage: "invite-org-user ORG_NAME INVITEE ROLE1(,ROLE2(,ROLE3))",
				},
			},
		},
	}
}

// Run initiates the plugin
func (p *AppCloudPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	var err error

	switch args[0] {
	case "create-backup":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument SERVICE_INSTANCE was not provided")
			return
		}

		err = p.CreateBackup(cliConnection, args[1])
	case "backups":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument SERVICE_INSTANCE was not provided")
			return
		}

		err = p.Backups(cliConnection, args[1])
	case "invite-space-user":
		if len(args) < 4 {
			fmt.Println("Incorrect Usage: the required arguments SPACE_NAME, INVITEE and/or ROLES were not provided")
			return
		}

		err = p.InviteSpaceUser(cliConnection, args[1], args[2], args[3])
	case "invite-org-user":
		if len(args) < 4 {
			fmt.Println("Incorrect Usage: the required arguments ORG_NAME, INVITEE and/or ROLES were not provided")
			return
		}

		err = p.InviteOrgUser(cliConnection, args[1], args[2], args[3])
	}

    if err != nil {
		fmt.Printf("\n%s\n", redBold(err.Error()))
	}
}

func main() {
	plugin.Start(new(AppCloudPlugin))
}
