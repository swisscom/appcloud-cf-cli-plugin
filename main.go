package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/cf/flags"
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
				Name:     "create-ssl-certificate",
				HelpText: "A new certificate will be issued and immediately installed",
				UsageDetails: plugin.Usage{
					Usage: "create-ssl-certificate ROUTE",
				},
			},
			{
				Name:     "turn-ssl-on",
				HelpText: "Certificate will be disabled for given route",
				UsageDetails: plugin.Usage{
					Usage: "turn-ssl-on ROUTE",
				},
			},
			{
				Name:     "turn-ssl-off",
				HelpText: "Certificate will be enabled for given route",
				UsageDetails: plugin.Usage{
					Usage: "turn-ssl-off ROUTE",
				},
			},
			{
				Name:     "revoke-ssl-certificate",
				HelpText: "Certificate will be revoked",
				UsageDetails: plugin.Usage{
					Usage: "revoke-ssl-certificate ROUTE",
				},
			},
			{
				Name:     "abort-ssl-certificate",
				HelpText: "Abort ssl certificate creation process",
				UsageDetails: plugin.Usage{
					Usage: "abort-ssl-certificate ROUTE",
				},
			},
			{
				Name:     "list-ssl-certificates",
				HelpText: "List available certificates for space",
				UsageDetails: plugin.Usage{
					Usage: "list-ssl-certificates",
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
	case "create-ssl-certificate":
		if len(args) < 3 {
			fmt.Println("Incorrect Usage: the required arguments `DOMAIN`and `ROUTE` was not provided")
			return
		}

		err = p.CreateSSLCertificate(cliConnection, args[2])
	case "turn-ssl-on":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}

		err = p.TurnSSLOn(cliConnection, args[1])
	case "turn-ssl-off":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}
		err = p.TurnSSLOff(cliConnection, args[1])
	case "revoke-ssl-certificate":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}

		err = p.RevokeSSLCertificate(cliConnection, args[1])
	case "abort-ssl-certificate":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}

		err = p.AbortSSLCertificateProcess(cliConnection, args[1])
	case "list-ssl-certificates":
		if len(args) < 1 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}

		err = p.ListSSLCertificates(cliConnection)
	case "tree":
		fc, err := parseArguments(args)
		if err != nil {
			fmt.Println("Incorrect Usage: Level option must be an int")
			return
		}
		value := fc.Int("l")

		err = p.Tree(cliConnection, value)
	}

	if err != nil {
		fmt.Printf("\n%s\n", redBold(err.Error()))
	}
}

func parseArguments(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewIntFlagWithDefault("level", "l", "Level of output", 3)
	err := fc.Parse(args...)

	return fc, err
}

func main() {
	plugin.Start(new(AppCloudPlugin))
}
