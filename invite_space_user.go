package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Inviting %s to space %s as %s...\n", cyanBold(invitee), cyanBold(spaceName), cyanBold(username))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return fmt.Errorf("Space %s not found", spaceName)
	}

	args := InviteSpaceUserRequest{
		Invitee: invitee,
		SpaceID: s.Guid,
		Roles:   strings.Split(roles, ","),
	}
	argsData, err := json.Marshal(args)
	if err != nil {
		return errors.New("Couldn't parse JSON data")
	}

	url := "/custom/space_invitations"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", url, "-d", string(argsData))

	if err != nil {
		return fmt.Errorf("Couldn't invite %s to space %s", invitee, spaceName)
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
