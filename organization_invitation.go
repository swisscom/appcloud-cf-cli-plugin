package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type OrgInvitationArgs struct {
   Invitee string `json:"invitee"`
   OrgID string `json:"organization_id"`
   Roles []string `json:"roles"`
}

// OrgInvitation sends an invitation to an invitee to join the organization with specified roles
func (p *AppCloudPlugin) InviteOrgUser(c plugin.CliConnection, organizationName string, invitee string, roles string) error {
	
	fmt.Printf("Creating an invitation for organization %s for invitee %s and providing roles %s...\n", cyanBold(organizationName), cyanBold(invitee), cyanBold(roles))

	s, err := c.GetOrg(organizationName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve organization for name %s", organizationName)
	}

        roleArr := strings.Split(roles,",")
        args := OrgInvitationArgs{
           Invitee: invitee,
           OrgID: s.Guid,
           Roles: roleArr,
        }
        jsonData, err := json.Marshal(args)
        if err != nil {
                return fmt.Errorf("Couldn't Marshal the args data into json")
        }
        jsonStr := string(jsonData)
	url := "/custom/organization_invitations"
	resLines, err := c.CliCommand("curl", "-H", "Content-Type: application/json", "-X", "POST", "-d", jsonStr, url)
	
	if err != nil {
		return fmt.Errorf("Couldn't invite %s to organization %s: %s", invitee, organizationName, err.Error())
	}

	resString := strings.Join(resLines, "")
	var bRes OrganizationInvResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.Entity.Status != "SENT" {
		return fmt.Errorf("Couldn't invite %s to organization %s. The status received from the response was %s", invitee, organizationName, bRes.Entity.Status)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("The user has been invited")
	return nil
}
