package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateBackupResponse is the response from the server from a create backup call
type ResendInvitationResponse struct {
	ServerResponse
	OrgInvitation
 	MyInvitations []OrganizationInvitation `json:"resources"`
}

// CreateBackup creates a backup for a service instance
func (p *AppCloudPlugin) ResendOrgInvitation(c plugin.CliConnection, orgName string, invitee string, roles string) error {

	fmt.Printf("Resending invitation to ORG %s as %s\n", cyanBold(orgName),cyanBold(roles))

	s, err := c.GetOrg(orgName)

	if err != nil {
		return fmt.Errorf("Couldn't retrieve org details %s, make sure org does exists", orgName)
	}

	url := fmt.Sprintf("/custom/organizations/%s/invitations",s.Guid)
	fmt.Printf("Resending invitation to ORG %s as %s\n", cyanBold(url),cyanBold(roles))
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "GET", url)


	if err != nil {
		return fmt.Errorf("Couldn't get invitations to org %s", orgName)
	}

	resString := strings.Join(resLines, "")
	var bRes ResendInvitationResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	if err != nil {
		return fmt.Errorf("Couldn't Get invitation to invitee  %s", orgName)
	}

	for i := 0; i < len(bRes.MyInvitations); i++{
			if (bRes.MyInvitations[i].OrganizationEntity.Invitee==invitee){

			url1 := fmt.Sprintf("/custom/organization_invitations/%s/resend",bRes.MyInvitations[i].Metadata.GUID)
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
	fmt.Println("Invitation resent to "+invitee+" for org "+orgName+" successfully")
	return nil
}
