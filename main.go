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
				Name:     "create-ssl-certificate",
				HelpText: "A new certificate will be issued and immediately installed",
				UsageDetails: plugin.Usage{
					Usage: "create-ssl-certificate ROUTE",
				},
			},
			{
				Name:     "install-ssl-certificate",
				HelpText: "A new certificate will be issued and immediately installed",
				UsageDetails: plugin.Usage{
					Usage: "install-ssl-certificate ROUTE",
				},
			},
			{
				Name:     "uninstall-ssl-certificate",
				HelpText: "A new certificate will be issued and immediately installed",
				UsageDetails: plugin.Usage{
					Usage: "uninstall-ssl-certificate ROUTE",
				},
			},
			{
				Name:     "revoke-ssl-certificate",
				HelpText: "A new certificate will be issued and immediately installed",
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
				HelpText: "A new certificate will be issued and immediately installed",
				UsageDetails: plugin.Usage{
					Usage: "list-ssl-certificates",
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
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}

		err = p.CreateSSLCertificate(cliConnection, args[1])
	case "install-ssl-certificate":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}

		err = p.InstallSSLCertificate(cliConnection, args[1])
	case "uninstall-ssl-certificate":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument ROUTE was not provided")
			return
		}
		err = p.UnInstallSSLCertificate(cliConnection, args[1])
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
	}

	if err != nil {
		fmt.Printf("\n%s\n", redBold(err.Error()))
	}
}

func main() {
	plugin.Start(new(AppCloudPlugin))
}
