package main

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// ResendOrgInvitation resends an existing org invitation.
func (p *AppCloudPlugin) ResendOrgInvitation(c plugin.CliConnection, invitee string, orgName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Resending invitation for %s to org %s as %s...", terminal.EntityNameColor(invitee), terminal.EntityNameColor(orgName), terminal.EntityNameColor(un))

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

	for _, i := range res.Resources {
		if i.Entity.Invitee == invitee {
			invURL := fmt.Sprintf("/custom/organization_invitations/%s/resend", i.Metadata.GUID)
			invResLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", invURL)
			invResString := strings.Join(invResLines, "")
			var invRes InvitationResponse

			err = json.Unmarshal([]byte(invResString), &invRes)
			if err != nil {
				return errors.Wrap(err, "Couldn't read JSON response from server")
			}

			if invRes.ErrorCode != "" {
				return fmt.Errorf("Error response from server: %s", invRes.Description)
			}
		}
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))
	p.ui.Say("Invitation resent")

	return nil
}
