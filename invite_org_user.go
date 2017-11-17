package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// InviteOrgUserRequest is a request to invite a user to an org.
type InviteOrgUserRequest struct {
	Invitee string   `json:"invitee"`
	OrgID   string   `json:"organization_id"`
	Roles   []string `json:"roles"`
}

// InviteOrgUser invites a user to join an org with a specific set of roles.
func (p *AppCloudPlugin) InviteOrgUser(c plugin.CliConnection, invitee string, orgName string, roles string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Inviting %s to org %s as %s...\n", cyanBold(invitee), cyanBold(orgName), cyanBold(username))

	o, err := c.GetOrg(orgName)
	if err != nil {
		return fmt.Errorf("Org %s not found", orgName)
	}

	args := InviteOrgUserRequest{
		Invitee: invitee,
		OrgID:   o.Guid,
		Roles:   strings.Split(roles, ","),
	}
	if len(args.Roles) == 1 && args.Roles[0] == "" {
		args.Roles = []string{}
	}
	argsData, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("Couldn't parse JSON data")
	}

	url := "/custom/organization_invitations"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", url, "-d", string(argsData))

	if err != nil {
		return fmt.Errorf("Couldn't invite %s to org %s", invitee, orgName)
	}

	resString := strings.Join(resLines, "")
	var res InvitationResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.Entity.Status != "SENT" && res.Entity.Status != "CONFIRMED" {
		return fmt.Errorf("Couldn't send invitation. Current status: %s", res.Entity.Status)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Invitation sent")
	return nil
}
