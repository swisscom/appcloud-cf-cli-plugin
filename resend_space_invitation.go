package main

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// ResendSpaceInvitation resends an existing space invitation.
func (p *AppCloudPlugin) ResendSpaceInvitation(c plugin.CliConnection, invitee string, spaceName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Resending invitation for %s to space %s as %s...", terminal.EntityNameColor(invitee), terminal.EntityNameColor(spaceName), terminal.EntityNameColor(un))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return errors.Wrap(err, "Space not found")
	}

	url := fmt.Sprintf("/custom/spaces/%s/invitations", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't get invitations for space")
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
			invURL := fmt.Sprintf("/custom/space_invitations/%s/resend", i.Metadata.GUID)
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
