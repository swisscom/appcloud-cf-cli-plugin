package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// ResendBillingAccountInvitation resends an existing billing account invitation.
func (p *AppCloudPlugin) ResendBillingAccountInvitation(c plugin.CliConnection, billingAccountName string, invitee string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Resending invitation for %s to billing account %s as %s...\n", cyanBold(invitee), cyanBold(billingAccountName), cyanBold(username))

	ba, err := getBillingAccount(c, billingAccountName)
	if err != nil {
		return fmt.Errorf("Billing account %s not found", billingAccountName)
	}

	url := fmt.Sprintf("/custom/accounts/%s/invitations", ba.Metadata.GUID)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't get invitations for billing account %s", billingAccountName)
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
			invURL := fmt.Sprintf("/custom/account_invitations/%s/resend", i.Metadata.GUID)
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
