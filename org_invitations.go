package main

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// OrgInvitations retrieves all invitations for an org.
func (p *AppCloudPlugin) OrgInvitations(c plugin.CliConnection, orgName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Gettings invitations to org %s as %s...", terminal.EntityNameColor(orgName), terminal.EntityNameColor(un))

	o, err := c.GetOrg(orgName)
	if err != nil {
		return errors.Wrap(err, "Org not found")
	}

	url := fmt.Sprintf("/custom/organizations/%s/invitations", o.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't get invitations for org")
	}

	resString := strings.Join(resLines, "")
	var res InvitationsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return fmt.Errorf("Error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	invitations := res.Resources
	if len(invitations) > 0 {
		table := p.ui.Table([]string{"invitee", "roles", "status"})
		for _, inv := range res.Resources {
			table.Add(inv.Entity.Invitee, strings.Join(inv.Entity.Roles, ","), inv.Entity.Status)
		}
		table.Print()
	} else {
		p.ui.Say("No invitations found")
	}
	return nil
}
