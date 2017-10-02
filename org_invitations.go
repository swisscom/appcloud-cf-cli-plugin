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

	if len(res.Resources) == 0 {
		fmt.Println("No invitations found")
		return nil
	}

	fmt.Println(bold("Invitee                                   status                 roles"))
	for _, inv := range res.Resources {
		fmt.Printf("%s                 %s               %s\n", inv.Entity.Invitee, inv.Entity.Status, inv.Entity.Roles)
	}
	return nil
}
