package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateBackupResponse is the response from the server from a create backup call
type SendSpaceInvitationResponse struct {
	ServerResponseError
	OrgInvitation
}

// CreateBackup creates a backup for a service instance
func (p *AppCloudPlugin) SendSpaceInvitation(c plugin.CliConnection, spaceName string, invitee string, roles string) error {

	fmt.Printf("sending invitation to space %s as %s\n", cyanBold(spaceName), cyanBold(roles))

	s, err := c.GetSpace(spaceName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space details %s, make sure space does exists", s.Name)
	}
	url := "/custom/space_invitations"
	jsondata := fmt.Sprintf("'{\"invitee\": \"%s\",\"roles\": [\"%s\"], \"space_id\": \"%s\"}'", invitee, roles, s.Guid)

	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", "-d", jsondata, url)

	if err != nil {
		return fmt.Errorf("Couldn't send invitation to SPACE %s", spaceName)
	}

	resString := strings.Join(resLines, "")
	var bRes SendInvitationResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Sent invitation to " + invitee + " for SPACE " + spaceName + " successfully")
	return nil
}
