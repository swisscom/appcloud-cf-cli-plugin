package main

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// InviteBillingAccountUserRequest is the request to invite a user to a billing account.
type InviteBillingAccountUserRequest struct {
	Invitee   string   `json:"invitee"`
	AccountID string   `json:"account_id"`
	Roles     []string `json:"roles"`
}

// InviteBillingAccountUser invites a user to join a billing account with a specific set of roles.
func (p *AppCloudPlugin) InviteBillingAccountUser(c plugin.CliConnection, invitee string, billingAccountName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Inviting %s to billing account %s as %s...", terminal.EntityNameColor(invitee), terminal.EntityNameColor(billingAccountName), terminal.EntityNameColor(un))

	ba, err := getBillingAccount(c, billingAccountName)
	if err != nil {
		return errors.Wrap(err, "Billing Account not found")
	}

	args := InviteBillingAccountUserRequest{
		Invitee:   invitee,
		AccountID: ba.Metadata.GUID,
		Roles:     []string{"accountOwner"},
	}
	argsData, err := json.Marshal(args)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse JSON data")
	}

	url := "/custom/account_invitations"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", url, "-d", string(argsData))

	if err != nil {
		return errors.Wrap(err, "Couldn't invite user to billing account")
	}

	resString := strings.Join(resLines, "")
	var res InvitationResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.Entity.Status != StatusSent && res.Entity.Status != StatusConfirmed {
		return errors.Wrap(err, "Couldn't send invitation")
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	p.ui.Say("Invitation sent")

	return nil
}
