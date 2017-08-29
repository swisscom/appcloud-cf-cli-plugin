package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"strings"
	"encoding/json"
)

func (p *AppCloudPlugin) Invitations(c plugin.CliConnection, invitationType string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("\nRetrieving your invitations as %s...\n\n", cyanBold(username))

	output := ""
	if invitationType == "all" {
		output += InvitationOutput(c, "account")
		output += InvitationOutput(c, "organization")
		output += InvitationOutput(c, "space")
	} else {
		output += InvitationOutput(c, invitationType)
	}

	if output == "" {
		output = "You have no invitations."
	}
	fmt.Println(output)
	return nil
}

func InvitationOutput(c plugin.CliConnection, invitationType string) string {

	url := fmt.Sprintf("/custom/%s_invitations", invitationType)

	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Sprint("Couldn't retrieve invitations for " + invitationType)
	}

	resString := strings.Join(resLines, "")

	var output string
	switch invitationType {
	case "account":
		var res AccountResponse
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return fmt.Sprint("Couldn't read JSON response")
		}

		for i := 0; i < res.TotalResults; i++ {
			output += fmt.Sprintf("%s invited you to billing account %s - %s\n", res.Resources[i].AccountEntity.ActorUsername, res.Resources[i].AccountEntity.AccountName, res.Resources[i].Metadata.GUID)
		}

	case "organization":
		var res OrganizationResponse
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return fmt.Sprint("Couldn't read JSON response")
		}

		for i := 0; i < res.TotalResults; i++ {
			output += fmt.Sprintf("%s invited you to organization %s - %s\n", res.Resources[i].OrganizationEntity.ActorUsername , res.Resources[i].OrganizationEntity.OrganizationName, res.Resources[i].Metadata.GUID)
		}

	case "space":
		var res SpaceResponse
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return fmt.Sprint("Couldn't read JSON response")
		}

		for i := 0; i < res.TotalResults; i++ {
			output += fmt.Sprintf("%s invited you to space %s - %s\n", res.Resources[i].SpaceEntity.ActorUsername, res.Resources[i].SpaceEntity.SpaceName, res.Resources[i].Metadata.GUID)
		}
	}

	return output

}

func (p *AppCloudPlugin) AcceptInvitation(c plugin.CliConnection, invitationType string ,invitationGUID string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("\nAccepting your %s invitation %s as %s...\n", invitationType, invitationGUID, cyanBold(username))

	output := AcceptInvitationOutput(c, invitationType, invitationGUID)

	fmt.Println(output)
	return nil
}

func AcceptInvitationOutput(c plugin.CliConnection, invitationType string, invitationGUID string) string {

	url := fmt.Sprintf("/custom/%s_invitations/%s/confirm", invitationType, invitationGUID)

	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url, "-X", "POST")
	if err != nil {
		return fmt.Sprint("Couldn't retrieve organizations")
	}

	resString := strings.Join(resLines, "")

	var output string
	switch invitationType {
	case "account":
		var res AccountInvitation
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return fmt.Sprint("Couldn't read JSON response")
		}

		output += fmt.Sprintf("Invitation has been accepted to account %s", res.AccountEntity.AccountName)

	case "organization":
		var res OrganizationInvitation
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return fmt.Sprint("Couldn't read JSON response")
		}

		output += fmt.Sprintf("Invitation has been accepted to organization %s", res.OrganizationEntity.OrganizationName)

	case "space":
		var res SpaceInvitation
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return fmt.Sprint("Couldn't read JSON response")
		}

		output += fmt.Sprintf("Invitation has been accepted to space %s", res.SpaceEntity.SpaceName)
	}

	return output
}
