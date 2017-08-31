package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// ResendSpaceInvitation resends an existing space invitation.
func (p *AppCloudPlugin) ResendSpaceInvitation(c plugin.CliConnection, spaceName string, invitee string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Resending invitation for %s to space %s as %s\n", cyanBold(invitee), cyanBold(spaceName), cyanBold(username))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return fmt.Errorf("Space %s not found", spaceName)
	}

	url := fmt.Sprintf("/custom/spaces/%s/invitations", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't get invitations for space %s", spaceName)
	}

	resString := strings.Join(resLines, "")
	var res InvitationsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	for _, i := range res.Resources {
		if i.Entity.Invitee == invitee {
			invURL := fmt.Sprintf("/custom/space_invitations/%s/resend", i.Metadata.GUID)
			invResLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", invURL)
			invResString := strings.Join(invResLines, "")
			var invRes InvitationResponse

			err = json.Unmarshal([]byte(invResString), &invRes)
			if err != nil {
				return errors.New("Couldn't read JSON response from server")
			}
		}
	}

	fmt.Println(greenBold("OK\n\n"))
	fmt.Println("Invitation resent")
	return nil
}
