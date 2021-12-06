package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// BillingAccountInvitations retrieves all invitations for a billing account.
func (p *Plugin) BillingAccountInvitations(c plugin.CliConnection, billingAccountName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Getting invitations for billing account %s as %s...", terminal.EntityNameColor(billingAccountName), terminal.EntityNameColor(un))

	ba, err := getBillingAccount(c, billingAccountName)
	if err != nil {
		return err
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

	p.ui.Say(terminal.SuccessColor("OK\n"))

	invitations := res.Resources
	if len(invitations) > 0 {
		table := p.ui.Table([]string{"invitee", "roles", "status"})
		for _, inv := range res.Resources {
			table.Add(inv.Entity.Invitee, strings.Join(inv.Entity.Roles, ","), inv.Entity.Status)
		}
		err := table.Print()
		if err != nil {
			return errors.Wrap(err, "Couldn't print table")
		}
	} else {
		p.ui.Say("No invitations found")
	}

	return nil
}
