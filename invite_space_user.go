package main

import (
	"encoding/json"

	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// InviteSpaceUserRequest is a request to invite a user to a space.
type InviteSpaceUserRequest struct {
	Invitee string   `json:"invitee"`
	SpaceID string   `json:"space_id"`
	Roles   []string `json:"roles"`
}

// InviteSpaceUser invites a user to join an space with a specific set of roles.
func (p *AppCloudPlugin) InviteSpaceUser(c plugin.CliConnection, invitee string, spaceName string, roles string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Inviting %s to space %s as %s...", terminal.EntityNameColor(invitee), terminal.EntityNameColor(spaceName), terminal.EntityNameColor(un))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return errors.Wrap(err, "Space not found")
	}

	args := InviteSpaceUserRequest{
		Invitee: invitee,
		SpaceID: s.Guid,
		Roles:   strings.Split(roles, ","),
	}
	if len(args.Roles) == 1 && args.Roles[0] == "" {
		args.Roles = []string{}
	}
	argsData, err := json.Marshal(args)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse JSON data")
	}

	url := "/custom/space_invitations"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", url, "-d", string(argsData))

	if err != nil {
		return errors.Wrap(err, "Couldn't invite to space")
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
