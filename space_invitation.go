package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type SpaceInvitationArgs struct {
	Invitee string   `json:"invitee"`
	SpaceID string   `json:"space_id"`
	Roles   []string `json:"roles"`
}

// SpaceInvitation sends an invitation to an invitee to join the space with specified roles
func (p *AppCloudPlugin) InviteSpaceUser(c plugin.CliConnection, spaceName string, invitee string, roles string) error {
	fmt.Printf("Creating an invitation for space %s for invitee %s and providing roles %s...\n", cyanBold(spaceName), cyanBold(invitee), cyanBold(roles))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space for name %s", spaceName)
	}

	roleArr := strings.Split(roles, ",")
	args := SpaceInvitationArgs{
		Invitee: invitee,
		SpaceID: s.Guid,
		Roles:   roleArr,
	}
	jsonData, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("Couldn't Marshal the args data into json")
	}
	jsonStr := string(jsonData)
	url := "/custom/space_invitations"
	resLines, err := c.CliCommand("curl", "-H", "Content-Type: application/json", "-X", "POST", "-d", jsonStr, url)
	if err != nil {
		return fmt.Errorf("Couldn't invite %s to space %s: %s", invitee, spaceName, err.Error())
	}

	resString := strings.Join(resLines, "")
	var bRes SpaceInvResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.Entity.Status != "SENT" {
		return fmt.Errorf("Couldn't invite %s to space %s. The status received from the response was %s", invitee, spaceName, bRes.Entity.Status)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("The user has been invited")
	return nil
}
