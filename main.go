package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/cf/flags"
	"strings"
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
				Name:     "tree",
				HelpText: "View organization tree",
				UsageDetails: plugin.Usage{
					Usage: "tree [--level | -l]\n   tree \n   tree -l 2 \n   tree --level 1",
					Options: map[string]string{
						"--level, l": "Level of output",
					},
				},
			},
			{
				Name:     "invitation",
				HelpText: "View invitations",
				UsageDetails: plugin.Usage{
					Usage: "invitation [--type | -t]\n   invitation\n   invitation -t account \n   invitation --type space",
					Options: map[string]string{
						"--type, t": "Type of invitation",
					},
				},
			},
			{
				Name:     "invitation-accept",
				HelpText: "Accept invitation",
				UsageDetails: plugin.Usage{
					Usage: "invitation-accept TYPE GUID",
				},
			},
			{
				Name:     "docker-repositories",
				HelpText: "List docker-repositories",
				UsageDetails: plugin.Usage{
					Usage: "docker-repositories [--org | -o]\n   docker-repositories\n   docker-repositories -o my-org \n   docker-repositories --org my-org",
					Options: map[string]string{
						"--type, t": "Type of invitation",
					},
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
	case "tree":
		fc, err := parseArguments(args)
		if err != nil {
			fmt.Println("Incorrect Usage: Level option must be an int")
			return
		}
		value := fc.Int("l")

		err = p.Tree(cliConnection, value)

	case "invitation":
		fc, err := parseArguments(args)
		if err != nil {
			fmt.Println("Incorrect Usage: Type option must be a string")
			return
		}
		value := strings.ToLower(fc.String("t"))

		valid := map[string]bool{"account": true, "organization": true, "space": true, "all": true}

		if !valid[value] {
			fmt.Println("Incorrect Usage: If type option is used, it must be account, organization, or space")
			return
		}

		err = p.Invitations(cliConnection, value)

	case "invitation-accept":
		if len(args) < 3 {
			fmt.Println("Incorrect Usage: the required argument GUID and/or TYPE was not provided")
			return
		}

		value := strings.ToLower(args[1])

		valid := map[string]bool{"account": true, "organization": true, "space": true}

		if !valid[value] {
			fmt.Println("Incorrect Usage: Invitation type is must be account, organization, or space")
			return
		}

		err = p.AcceptInvitation(cliConnection, args[1], args[2])
	case "docker-repositories":
		fc, err := parseArguments(args)
		if err != nil {
			fmt.Println("Incorrect Usage: Organization option must be a string")
			return
		}
		value := strings.ToLower(fc.String("o"))

		err = p.DockerRepository(cliConnection, value)
	}

	if err != nil {
		fmt.Printf("\n%s\n", redBold(err.Error()))
	}
}

func parseArguments(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewIntFlagWithDefault("level", "l", "Level of output", 3)
	fc.NewStringFlagWithDefault("type", "t", "Type of invitation", "all")
	fc.NewStringFlagWithDefault("org", "o", "Organization", "none")
	err := fc.Parse(args...)

	return fc, err
}

func main() {
	plugin.Start(new(AppCloudPlugin))
}
