package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateBackupResponse is the response from the server from a create backup call
type SendInvitationResponse struct {
	OrgInvitation
	ServerResponseError
}

// CreateBackup creates a backup for a service instance
func (p *AppCloudPlugin) SendOrgInvitation(c plugin.CliConnection, orgName string, invitee string, roles string) error {

	fmt.Printf("sending invitation to ORG %s as %s\n", cyanBold(orgName), cyanBold(roles))

	s, err := c.GetOrg(orgName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve org details %s, make sure org does exists", orgName)
	}
	url := "/custom/organization_invitations"
	jsondata := fmt.Sprintf("'{\"invitee\": \"%s\",\"roles\": [\"%s\"], \"organization_id\": \"%s\"}'", invitee, roles, s.Guid)

	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", "-d", jsondata, url)

	if err != nil {
		return fmt.Errorf("Couldn't send invitation to org %s", orgName)
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

	fmt.Println("Sent invitation to " + invitee + " for org " + orgName + " successfully")
	return nil
}
