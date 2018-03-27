package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// AcceptInvitation accepts a pending invitation.
func (p *Plugin) AcceptInvitation(c plugin.CliConnection, invitationGUID string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Accepting invitation as %s...", terminal.EntityNameColor(un))

	invs, err := getAllInvitations(c)
	if err != nil {
		return err
	}

	var inv Invitation
	for _, i := range invs {
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
		return errors.Wrap(err, "Couldn't accept invitation")
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
	p.ui.Say("Invitation accepted")

	return nil
}
