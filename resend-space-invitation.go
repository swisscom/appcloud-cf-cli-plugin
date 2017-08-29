package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateBackupResponse is the response from the server from a create backup call
type ResendSpaceInvitationResponse struct {
	ServerResponse
	OrgInvitation
	MyInvitations []OrganizationInvitation `json:"resources"`

}

// CreateBackup creates a backup for a service instance
func (p *AppCloudPlugin) ResendSpaceInvitation(c plugin.CliConnection, spaceName string, invitee string, roles string) error {

	fmt.Printf("sending invitation to space %s as %s\n", cyanBold(spaceName),cyanBold(roles))
	
	s, err := c.GetSpace(spaceName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space details %s, make sure space does exists", s.Name)
	}
	url := fmt.Sprintf("/custom/spaces/%s/invitations",s.Guid)
	fmt.Printf(url)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl",  url)

	if err != nil {
		return fmt.Errorf("Couldn't get invitations to Invitee %s in this space ", invitee)
	}

	resString := strings.Join(resLines, "")
	var bRes ResendSpaceInvitationResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	for i := 0; i < len(bRes.MyInvitations); i++{
		if (bRes.MyInvitations[i].OrganizationEntity.Invitee==invitee){

			url1 := fmt.Sprintf("/custom/space_invitations/%s/resend",bRes.MyInvitations[i].Metadata.GUID)
			resLines2, err1 := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", url1)
			resString2 := strings.Join(resLines2, "")
			var bRes1 OrganizationInvitation

			err1 = json.Unmarshal([]byte(resString2), &bRes1)
			if err1 != nil {
				return errors.New("Couldn't read JSON response")
			}
		}
	}

	fmt.Println(greenBold("OK\n\n"))
	fmt.Println("Sent invitation to "+invitee+" for SPACE "+spaceName+" successfully")
	return nil
}
