package appcloud

import (
	"os"

	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin"
)

// Plugin is the Swisscom Application Cloud cf CLI plugin.
type Plugin struct {
	ui terminal.UI
}

// GetMetadata retrieves the metadata for the plugin.
func (p *Plugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Swisscom Application Cloud",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
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
			{
				Name:     "delete-backup",
				HelpText: "Delete a backup of a service instance",
				UsageDetails: plugin.Usage{
					Usage: "delete-backup SERVICE_INSTANCE BACKUP_GUID",
				},
			},

			// Receive Invitations
			{
				Name:     "invitations",
				HelpText: "List your pending invitations",
				UsageDetails: plugin.Usage{
					Usage: "invitations",
				},
			},
			{
				Name:     "accept-invitation",
				HelpText: "Accept a pending invitation",
				UsageDetails: plugin.Usage{
					Usage: "accept-invitation INVITATION_GUID",
				},
			},
			{
				Name:     "decline-invitation",
				HelpText: "Decline a pending invitation",
				UsageDetails: plugin.Usage{
					Usage: "decline-invitation INVITATION_GUID",
				},
			},

			// Send invitations
			{
				Name:     "billing-account-invitations",
				HelpText: "List all invitations for a billing account",
				UsageDetails: plugin.Usage{
					Usage: "billing-account-invitations BILLING_ACCOUNT",
				},
			},
			{
				Name:     "org-invitations",
				HelpText: "List all invitations for an org",
				UsageDetails: plugin.Usage{
					Usage: "org-invitations ORG",
				},
			},
			{
				Name:     "space-invitations",
				HelpText: "List all invitations for a space",
				UsageDetails: plugin.Usage{
					Usage: "space-invitations SPACE",
				},
			},
			{
				Name:     "invite-billing-account-user",
				HelpText: "Invite a user to a billing account as an 'accountOwner'",
				UsageDetails: plugin.Usage{
					Usage: "invite-org-user USERNAME BILLING_ACCOUNT",
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
				Name:     "invite-space-user",
				HelpText: "Invite a user to a space",
				UsageDetails: plugin.Usage{
					Usage: "invite-space-user USERNAME SPACE ROLE1[,ROLE2[,ROLE3]]",
				},
			},
			{
				Name:     "resend-billing-account-invitation",
				HelpText: "Resend an existing billing account invitation",
				UsageDetails: plugin.Usage{
					Usage: "resend-billing-account-invitation USERNAME BILLING_ACCOUNT",
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
					Usage: "resend-space-invitation USERNAME SPACE",
				},
			},

			// SSL certificates
			{
				Name:     "ssl-certificates",
				HelpText: "List SSL certificates for the current space",
				UsageDetails: plugin.Usage{
					Usage: "ssl-certificates",
				},
			},
			{
				Name:     "create-ssl-certificate",
				HelpText: "Create and enable a SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "create-ssl-certificate DOMAIN [--hostname HOSTNAME] [--key-type KEY_TYPE]",
					Options: map[string]string{
						"-hostname, -n": "Hostname for the HTTP route",
						"-key-type, -kt": "Key type for the certificate (e.g. RSA, ECDSA)",
					},
				},
			},
			{
				Name:     "revoke-ssl-certificate",
				HelpText: "Revoke an existing SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "revoke-ssl-certificate DOMAIN [--hostname HOSTNAME]",
					Options: map[string]string{
						"-hostname, -n": "Hostname for the HTTP route",
					},
				},
			},
			{
				Name:     "enable-ssl",
				HelpText: "Enable an existing SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "enable-ssl DOMAIN [--hostname HOSTNAME]",
					Options: map[string]string{
						"-hostname, -n": "Hostname for the HTTP route",
					},
				},
			},
			{
				Name:     "disable-ssl",
				HelpText: "Disable an existing SSL certificate for a route",
				UsageDetails: plugin.Usage{
					Usage: "disable-ssl DOMAIN [--hostname HOSTNAME]",
					Options: map[string]string{
						"-hostname, -n": "Hostname for the HTTP route",
					},
				},
			},
			{
				Name:     "ssl-enabled",
				HelpText: "Reports whether SSL is enabled for a route",
				UsageDetails: plugin.Usage{
					Usage: "ssl-enabled DOMAIN [--hostname HOSTNAME]",
					Options: map[string]string{
						"-hostname, -n": "Hostname for the HTTP route",
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
						"--depth, -d": "Depth of the tree output (0: orgs, 1: spaces, 2: apps and service instances)",
					},
				},
			},

			// Service events
			{
				Name:     "service-events",
				HelpText: "Show recent service instance events",
				UsageDetails: plugin.Usage{
					Usage: "service-events SERVICE_INSTANCE",
				},
			},
		},
	}
}

// Run initiates the plugin.
func (p *Plugin) Run(cliConnection plugin.CliConnection, args []string) {
	var err error

	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	p.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	switch args[0] {

	// Backups
	case "backups":
		if len(args) != 2 {
			p.ui.Say("Incorrect Usage: the required argument SERVICE_INSTANCE was not provided")
			return
		}

		err = p.Backups(cliConnection, args[1])
	case "create-backup":
		if len(args) != 2 {
			p.ui.Say("Incorrect Usage: the required argument SERVICE_INSTANCE was not provided")
			return
		}

		err = p.CreateBackup(cliConnection, args[1])
	case "restore-backup":
		if len(args) != 3 {
			p.ui.Say("Incorrect Usage: the required arguments SERVICE_INSTANCE and/or BACKUP_GUID were not provided")
			return
		}

		err = p.RestoreBackup(cliConnection, args[1], args[2])
	case "delete-backup":
		if len(args) != 3 {
			p.ui.Say("Incorrect Usage: the required arguments SERVICE_INSTANCE and/or BACKUP_GUID were not provided")
			return
		}

		err = p.DeleteBackup(cliConnection, args[1], args[2])

	// Receive Invitations
	case "invitations":
		err = p.Invitations(cliConnection)
	case "accept-invitation":
		if len(args) != 2 {
			p.ui.Say("Incorrect Usage: the required argument INVITATION_GUID was not provided")
			return
		}

		err = p.AcceptInvitation(cliConnection, args[1])
	case "decline-invitation":
		if len(args) != 2 {
			p.ui.Say("Incorrect Usage: the required argument INVITATION_GUID was not provided")
			return
		}

		err = p.DeclineInvitation(cliConnection, args[1])

	// Send invitations
	case "billing-account-invitations":
		if len(args) != 2 {
			p.ui.Say("Incorrect Usage: the required argument BILLING_ACCOUNT was not provided")
			return
		}

		err = p.BillingAccountInvitations(cliConnection, args[1])
	case "org-invitations":
		if len(args) != 2 {
			p.ui.Say("Incorrect Usage: the required argument ORG was not provided")
			return
		}

		err = p.OrgInvitations(cliConnection, args[1])
	case "space-invitations":
		if len(args) != 2 {
			p.ui.Say("Incorrect Usage: the required argument SPACE was not provided")
			return
		}

		err = p.SpaceInvitations(cliConnection, args[1])
	case "invite-billing-account-user":
		if len(args) != 3 {
			p.ui.Say("Incorrect Usage: the required arguments USERNAME and/or BILLING_ACCOUNT were not provided")
			return
		}

		err = p.InviteBillingAccountUser(cliConnection, args[1], args[2])
	case "invite-org-user":
		if len(args) != 4 {
			p.ui.Say("Incorrect Usage: the required arguments USERNAME, ORG and/or ROLES were not provided")
			return
		}

		err = p.InviteOrgUser(cliConnection, args[1], args[2], args[3])
	case "invite-space-user":
		if len(args) != 4 {
			p.ui.Say("Incorrect Usage: the required arguments USERNAME, SPACE and/or ROLES were not provided")
			return
		}

		err = p.InviteSpaceUser(cliConnection, args[1], args[2], args[3])
	case "resend-billing-account-invitation":
		if len(args) != 3 {
			p.ui.Say("Incorrect Usage: the required arguments USERNAME and/org BILLING_ACCOUNT were not provided")
			return
		}
		err = p.ResendBillingAccountInvitation(cliConnection, args[1], args[2])
	case "resend-org-invitation":
		if len(args) != 3 {
			p.ui.Say("Incorrect Usage: the required arguments USERNAME and/org ORG were not provided")
			return
		}
		err = p.ResendOrgInvitation(cliConnection, args[1], args[2])
	case "resend-space-invitation":
		if len(args) != 3 {
			p.ui.Say("Incorrect Usage: the required arguments USERNAME and/or SPACE were not provided")
			return
		}
		err = p.ResendSpaceInvitation(cliConnection, args[1], args[2])

	// SSL Certificates
	case "ssl-certificates":
		err = p.SSLCertificates(cliConnection)
	case "create-ssl-certificate":
		if len(args) < 2 || len(args) > 4 {
			p.ui.Say("Incorrect Usage: the required argument DOMAIN was not provided")
			return
		}

		var fc flags.FlagContext
		fc, err = parseSSLCertificateArgs(args)
		if err != nil {
			p.ui.Say(err.Error())
			return
		}

		err = p.CreateSSLCertificate(cliConnection, args[1], fc.String("hostname"), fc.String("key-type"))
	case "revoke-ssl-certificate":
		if len(args) < 2 {
			p.ui.Say("Incorrect Usage: the required argument DOMAIN was not provided")
			return
		}

		var fc flags.FlagContext
		fc, err = parseSSLCertificateArgs(args)
		if err != nil {
			p.ui.Say("Incorrect Usage: HOSTNAME must be a string")
			return
		}

		err = p.RevokeSSLCertificate(cliConnection, args[1], fc.String("hostname"))
	case "enable-ssl":
		if len(args) < 2 {
			p.ui.Say("Incorrect Usage: the required argument DOMAIN was not provided")
			return
		}

		var fc flags.FlagContext
		fc, err = parseSSLCertificateArgs(args)
		if err != nil {
			p.ui.Say("Incorrect Usage: HOSTNAME must be a string")
			return
		}

		err = p.EnableSSL(cliConnection, args[1], fc.String("hostname"))
	case "disable-ssl":
		if len(args) < 2 {
			p.ui.Say("Incorrect Usage: the required argument DOMAIN was not provided")
			return
		}

		var fc flags.FlagContext
		fc, err = parseSSLCertificateArgs(args)
		if err != nil {
			p.ui.Say("Incorrect Usage: HOSTNAME must be a string")
			return
		}

		err = p.DisableSSL(cliConnection, args[1], fc.String("hostname"))
	case "ssl-enabled":
		if len(args) < 2 {
			p.ui.Say("Incorrect Usage: the required argument DOMAIN was not provided")
			return
		}

		var fc flags.FlagContext
		fc, err = parseSSLCertificateArgs(args)
		if err != nil {
			p.ui.Say("Incorrect Usage: HOSTNAME must be a string")
			return
		}

		err = p.SSLEnabled(cliConnection, args[1], fc.String("hostname"))

	// Tree
	case "tree":
		var fc flags.FlagContext
		fc, err = parseTreeArgs(args)
		if err != nil {
			p.ui.Say("Incorrect Usage: DEPTH must be an integer")
			return
		}

		err = p.Tree(cliConnection, fc.Int("depth"))

	// Service events
	case "service-events":
		if len(args) < 2 {
			p.ui.Say("Incorrect Usage: the required argument SERVICE_INSTANCE was not provided")
			return
		}

		err = p.ServiceEvents(cliConnection, args[1])
	}

	if err != nil {
		p.ui.Say("\n")
		p.ui.Say(err.Error())
		p.ui.Say(terminal.FailureColor("FAILED"))
	}
}

// parseSSLCertificateArgs parses the arguments passed to a ssl certificate command.
func parseSSLCertificateArgs(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewStringFlag("hostname", "n", "Hostname for the HTTP route")
	fc.NewStringFlag("key-type", "kt", "Key type for the certificate")
	err := fc.Parse(args...)
	if err != nil {
		return nil, err
	}

	return fc, nil
}

// parseTreeArgs parses the arguments passed to a tree command.
func parseTreeArgs(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewIntFlagWithDefault("depth", "d", "Level of output", 2)
	err := fc.Parse(args...)
	if err != nil {
		return nil, err
	}

	return fc, err
}
