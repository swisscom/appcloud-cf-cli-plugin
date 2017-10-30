package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// OrgInvitations retrieves all invitations for an org.
func (p *AppCloudPlugin) OrgInvitations(c plugin.CliConnection, orgName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Gettings invitations to org %s as %s...\n", cyanBold(orgName), cyanBold(username))

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
