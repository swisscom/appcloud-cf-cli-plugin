package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// BillingAccountInvitationArgs are the arguments required by the server to invite a user to a billing account.
type BillingAccountInvitationArgs struct {
	Invitee   string   `json:"invitee"`
	AccountID string   `json:"account_id"`
	Roles     []string `json:"roles"`
}

// InvitationResponse is a response from the server to an invitation request.
type InvitationResponse struct {
	Invitation
	ServerResponseError
}

// InviteBillingAccountUser invites a user to join a billing account with a specific set of roles.
func (p *AppCloudPlugin) InviteBillingAccountUser(c plugin.CliConnection, billingAccountName string, invitee string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Inviting %s to billing account %s as %s\n", cyanBold(invitee), cyanBold(billingAccountName), cyanBold(username))

	ba, err := GetBillingAccount(c, billingAccountName)
	if err != nil {
		return fmt.Errorf("Billing Account %s not found", billingAccountName)
	}

	args := BillingAccountInvitationArgs{
		Invitee:   invitee,
		AccountID: ba.Metadata.GUID,
		Roles:     []string{"accountOwner"},
	}
	argsData, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("Couldn't parse JSON data")
	}

	url := "/custom/account_invitations"
	resLines, err := c.CliCommand("curl", "-H", "Content-Type: application/json", "-X", "POST", url, "-d", string(argsData))

	if err != nil {
		return fmt.Errorf("Couldn't invite %s to billing account %s", invitee, billingAccountName)
	}

	resString := strings.Join(resLines, "")
	var res InvitationResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.Entity.Status != "SENT" {
		return fmt.Errorf("Couldn't send invitation. Current status: %s", res.Entity.Status)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Invitation sent")
	return nil
}
