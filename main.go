package main

import (
	"fmt"

	"strings"

	"code.cloudfoundry.org/cli/cf/flags"
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

			// Backups
			{
				Name:     "backups",
				HelpText: "List backups for a service instance",
				UsageDetails: plugin.Usage{
					Usage: "backups SERVICE_INSTANCE",
				},
			},
			{
				Name:     "create-backup",
				HelpText: "Create a backup for a service instance",
				UsageDetails: plugin.Usage{
					Usage: "create-backup SERVICE_INSTANCE",
				},
			},
			{
				Name:     "restore-backup",
				HelpText: "Restore a backup on a service instance",
				UsageDetails: plugin.Usage{
					Usage: "restore-backup SERVICE_INSTANCE BACKUP_GUID",
				},
			},

			// Invitations
			{
				Name:     "invitations",
				HelpText: "List your currently pending invitations",
				UsageDetails: plugin.Usage{
					Usage: "invitations",
				},
			},
			{
				Name:     "invite-space-user",
				HelpText: "Invite a user to a space",
				UsageDetails: plugin.Usage{
					Usage: "invite-space-user USERNAME ORG SPACE ROLE1[,ROLE2[,ROLE3]]",
				},
			},
			{
				Name:     "invite-org-user",
				HelpText: "Invite a user to an org",
				UsageDetails: plugin.Usage{
					Usage: "invite-org-user USERNAME ORG ROLE1[,ROLE2]]",
				},
			},
			{
				Name:     "resend-org-invitation",
				HelpText: "Resend an existing org invitation",
				UsageDetails: plugin.Usage{
					Usage: "resend-org-invitation USERNAME ORG",
				},
			},
			{
				Name:     "resend-space-invitation",
				HelpText: "Resend an existing space invitation",
				UsageDetails: plugin.Usage{
					Usage: "resend-org-invitation USERNAME ORG SPACE",
				},
			},
			{
				Name:     "accept-invitation",
				HelpText: "Accept a pending invitation",
				UsageDetails: plugin.Usage{
					Usage: "accept-invitation INVITATION_GUID",
				},
			},

			// SSL certificates
			{
				Name:     "ssl-certificates",
				HelpText: "List SSL certificates for the current space",
				UsageDetails: plugin.Usage{
					Usage: "list-ssl-certificates [--space SPACE]",
					Options: map[string]string{
						"--space, s": "Space",
					},
				},
			},
			{
				Name:     "create-ssl-certificate",
				HelpText: "Create and enable an SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "create-ssl-certificate DOMAIN [--hostname HOSTNAME] [--path PATH]",
					Options: map[string]string{
						"--hostname, n": "Hostname for the HTTP route (required for shared domains)",
						"--path, p":     "Path for the HTTP route",
					},
				},
			},
			{
				Name:     "revoke-ssl-certificate",
				HelpText: "Revoke an existing SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "revoke-ssl-certificate DOMAIN [--hostname HOSTNAME] [--path PATH]",
					Options: map[string]string{
						"--hostname, n": "Hostname for the HTTP route (required for shared domains)",
						"--path, p":     "Path for the HTTP route",
					},
				},
			},
			{
				Name:     "enable-ssl",
				HelpText: "Enable an existing SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "enable-ssl DOMAIN [--hostname HOSTNAME] [--path PATH]",
					Options: map[string]string{
						"--hostname, n": "Hostname for the HTTP route (required for shared domains)",
						"--path, p":     "Path for the HTTP route",
					},
				},
			},
			{
				Name:     "disable-ssl",
				HelpText: "Disable an existing SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "disable-ssl DOMAIN [--hostname HOSTNAME] [--path PATH]",
					Options: map[string]string{
						"--hostname, n": "Hostname for the HTTP route (required for shared domains)",
						"--path, p":     "Path for the HTTP route",
					},
				},
			},
			{
				Name:     "ssl-enabled",
				HelpText: "Reports whether SSL is enabled for a route",
				UsageDetails: plugin.Usage{
					Usage: "ssl-enabled DOMAIN [--hostname HOSTNAME] [--path PATH]",
					Options: map[string]string{
						"--hostname, n": "Hostname for the HTTP route (required for shared domains)",
						"--path, p":     "Path for the HTTP route",
					},
				},
			},

			// Docker registry
			{
				Name:     "docker-repositories",
				HelpText: "List docker-repositories",
				UsageDetails: plugin.Usage{
					Usage: "docker-repositories [--org ORG]",
					Options: map[string]string{
						"--org, o": "Organization",
					},
				},
			},

			// Tree
			{
				Name:     "tree",
				HelpText: "View organization tree",
				UsageDetails: plugin.Usage{
					Usage: "tree [--depth DEPTH]",
					Options: map[string]string{
						"--depth, d": "Depth of the tree output",
					},
				},
			},
			{
				Name:     "service-events",
				HelpText: "Show service events for SERVICE_INSTANCE",
				UsageDetails: plugin.Usage{
					Usage: "service-events  SERVICE_INSTANCE",
				},
			},
		},
	}
}

// Run initiates the plugin
func (p *AppCloudPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	var err error

	switch args[0] {

	// Backups
	case "backups":
		if len(args) != 2 {
			fmt.Println("Incorrect Usage: the required argument SERVICE_INSTANCE was not provided")
			return
		}

		err = p.Backups(cliConnection, args[1])
	case "create-backup":
		if len(args) != 2 {
			fmt.Println("Incorrect Usage: the required argument SERVICE_INSTANCE was not provided")
			return
		}

		err = p.CreateBackup(cliConnection, args[1])
	case "restore-backup":
		if len(args) != 3 {
			fmt.Println("Incorrect Usage: the required arguments SERVICE_INSTANCE and/or BACKUP_GUID were not provided")
			return
		}

		err = p.RestoreBackup(cliConnection, args[1], args[2])

	// Invitations
	case "invite-space-user":
		if len(args) != 4 {
			fmt.Println("Incorrect Usage: the required arguments SPACE, INVITEE and/or ROLES were not provided")
			return
		}

		err = p.InviteSpaceUser(cliConnection, args[1], args[2], args[3])
	case "invite-org-user":
		if len(args) < 4 {
			fmt.Println("Incorrect Usage: the required arguments ORG, INVITEE and/or ROLES were not provided")
			return
		}

		err = p.InviteOrgUser(cliConnection, args[1], args[2], args[3])

	// SSL Certificates
	case "ssl-certificates":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required argument DOMAIN was not provided")
			return
		}

		err = p.ListSSLCertificates(cliConnection)
	case "create-ssl-certificate":
		if len(args) < 3 {
			fmt.Println("Incorrect Usage: the required arguments DOMAIN and ROUTE was not provided")
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
	case "send-org-invitation":
		if len(args) < 4 {
			fmt.Println("Incorrect Usage: the required arguments was not provided")
			return
		}
		err = p.SendOrgInvitation(cliConnection, args[1], args[2], args[3])
	case "resend-org-invitation":
		if len(args) < 4 {
			fmt.Println("Incorrect Usage: the required arguments was not provided")
			return
		}
		err = p.ResendOrgInvitation(cliConnection, args[1], args[2], args[3])
	case "resend-space-invitation":
		if len(args) < 4 {
			fmt.Println("Incorrect Usage: the required arguments was not provided")
			return
		}
		err = p.ResendSpaceInvitation(cliConnection, args[1], args[2], args[3])

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
	case "service-events":
		if len(args) < 2 {
			fmt.Println("Incorrect Usage: the required arguments was not provided")
			return
		}

		err = p.ServiceEvents(cliConnection, args[1])
	}

	if err != nil {
		fmt.Print(redBold("FAILED\n\n"))
		fmt.Println(err.Error())
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
