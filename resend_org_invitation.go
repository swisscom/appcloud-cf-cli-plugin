package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// ResendOrgInvitation resends an existing org invitation.
func (p *AppCloudPlugin) ResendOrgInvitation(c plugin.CliConnection, orgName string, invitee string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Resending invitation for %s to org %s as %s...\n", cyanBold(invitee), cyanBold(orgName), cyanBold(username))

	o, err := c.GetOrg(orgName)
	if err != nil {
		return fmt.Errorf("Org %s not found", orgName)
	}

	url := fmt.Sprintf("/custom/organizations/%s/invitations", o.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't get invitations for org %s", orgName)
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
			invURL := fmt.Sprintf("/custom/organization_invitations/%s/resend", i.Metadata.GUID)
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
