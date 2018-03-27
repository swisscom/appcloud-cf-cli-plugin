package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// ResendBillingAccountInvitation resends an existing billing account invitation.
func (p *Plugin) ResendBillingAccountInvitation(c plugin.CliConnection, invitee string, billingAccountName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Resending invitation for %s to billing account %s as %s...", terminal.EntityNameColor(invitee), terminal.EntityNameColor(billingAccountName), terminal.EntityNameColor(un))

	ba, err := getBillingAccount(c, billingAccountName)
	if err != nil {
		return errors.Wrap(err, "Billing account not found")
	}

	url := fmt.Sprintf("/custom/accounts/%s/invitations", ba.Metadata.GUID)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't get invitations for billing account")
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
			invURL := fmt.Sprintf("/custom/account_invitations/%s/resend", i.Metadata.GUID)
			invResLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", invURL)
			if err != nil {
				return errors.Wrap(err, "Couldn't resend billing account invitation")
			}
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
