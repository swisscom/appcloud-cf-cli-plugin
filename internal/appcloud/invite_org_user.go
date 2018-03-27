package appcloud

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// InviteOrgUserRequest is a request to invite a user to an org.
type InviteOrgUserRequest struct {
	Invitee string   `json:"invitee"`
	OrgID   string   `json:"organization_id"`
	Roles   []string `json:"roles"`
}

// InviteOrgUser invites a user to join an org with a specific set of roles.
func (p *Plugin) InviteOrgUser(c plugin.CliConnection, invitee string, orgName string, roles string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Inviting %s to org %s as %s...", terminal.EntityNameColor(invitee), terminal.EntityNameColor(orgName), terminal.EntityNameColor(un))

	o, err := c.GetOrg(orgName)
	if err != nil {
		return errors.Wrap(err, "Org not found")
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
		return errors.Wrap(err, "Couldn't parse JSON data")
	}

	url := "/custom/organization_invitations"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-H", "Content-Type: application/json", "-X", "POST", url, "-d", string(argsData))

	if err != nil {
		return errors.Wrap(err, "Couldn't invite user to org")
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
