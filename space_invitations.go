package main

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// SpaceInvitations retrieves all invitations for a space.
func (p *AppCloudPlugin) SpaceInvitations(c plugin.CliConnection, spaceName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Gettings invitations to space %s as %s...", terminal.EntityNameColor(spaceName), terminal.EntityNameColor(un))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return errors.Wrap(err, "Space not found")
	}

	url := fmt.Sprintf("/custom/spaces/%s/invitations", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't get invitations for space")
	}

	resString := strings.Join(resLines, "")
	var res InvitationsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	invitations := res.Resources
	if len(invitations) > 0 {
		table := p.ui.Table([]string{"invitee", "roles", "status"})
		for _, inv := range res.Resources {
			table.Add(inv.Entity.Invitee, strings.Join(inv.Entity.Roles, ","), inv.Entity.Status)
		}
		table.Print()
	} else {
		p.ui.Say("No invitations found")
	}
	return nil
}
