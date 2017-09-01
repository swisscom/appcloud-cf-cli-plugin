package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// AcceptInvitation accepts a pending invitation.
func (p *AppCloudPlugin) AcceptInvitation(c plugin.CliConnection, invitationGUID string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Accepting invitation as %s...\n", cyanBold(username))

	invitations, err := getAllInvitations(c)
	if err != nil {
		return err
	}

	var inv Invitation
	for _, i := range invitations {
		if i.Metadata.GUID == invitationGUID {
			inv = i
			break
		}
	}

	if inv.Metadata.GUID == "" {
		return errors.New("Invitation not found")
	}

	t, _ := invitationEntityTypeAndName(inv)
	url := fmt.Sprintf("/custom/%s_invitations/%s/confirm", t, inv.Metadata.GUID)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		return fmt.Errorf("Couldn't accept invitation %s", invitationGUID)
	}

	resString := strings.Join(resLines, "")
	var res InvitationResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Print(greenBold("OK\n\n"))
	fmt.Println("Invitation accepted")
	return nil
}
