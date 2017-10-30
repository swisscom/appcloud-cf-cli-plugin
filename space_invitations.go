package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// SpaceInvitations retrieves all invitations for a space.
func (p *AppCloudPlugin) SpaceInvitations(c plugin.CliConnection, spaceName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Gettings invitations to space %s as %s...\n", cyanBold(spaceName), cyanBold(username))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return fmt.Errorf("Space %s not found", spaceName)
	}

	url := fmt.Sprintf("/custom/spaces/%s/invitations", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't get invitations for space %s", spaceName)
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
