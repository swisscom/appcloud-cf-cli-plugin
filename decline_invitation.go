package main

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// DeclineInvitation declines a pending invitation.
func (p *AppCloudPlugin) DeclineInvitation(c plugin.CliConnection, invitationGUID string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Declining invitation as %s...", terminal.EntityNameColor(un))

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
	url := fmt.Sprintf("/custom/%s_invitations/%s/reject", t, inv.Metadata.GUID)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't decline invitation")
	}

	resString := strings.Join(resLines, "")
	var res InvitationResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return fmt.Errorf("Error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))
	p.ui.Say("Invitation declined")

	return nil
}
