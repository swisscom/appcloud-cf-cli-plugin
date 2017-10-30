package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// BillingAccountInvitations retrieves all invitations for a billing account.
func (p *AppCloudPlugin) BillingAccountInvitations(c plugin.CliConnection, billingAccountName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Gettings invitations to billing account %s as %s...\n", cyanBold(billingAccountName), cyanBold(username))

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

	fmt.Println(greenBold("OK\n\n"))

	invitations := res.Resources
	if len(invitations) > 0 {
		table := NewTable([]string{"invitee", "roles", "status"})
		for _, inv := range res.Resources {
			table.Add(inv.Entity.Invitee, strings.Join(inv.Entity.Roles, ","), inv.Entity.Status)
		}
		table.Print()
	} else {
		fmt.Println("No invitations found")
	}
	return nil
}
