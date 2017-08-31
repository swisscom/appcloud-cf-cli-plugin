package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// InvitationsResponse is a response from the server containing invitations.
type InvitationsResponse struct {
	Resources []Invitation `json:"resources"`
	ServerResponsePagination
	ServerResponseError
}

// Invitations retrieves a user's invitations.
func (p *AppCloudPlugin) Invitations(c plugin.CliConnection) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Getting invitations as %s...\n", cyanBold(username))

	invitations, err := getAllInvitations(c)
	if err != nil {
		return err
	}

	fmt.Print(greenBold("OK\n\n"))

	if len(invitations) == 0 {
		fmt.Println("No invitations found")
		return nil
	}

	fmt.Println(bold("GUID                                   entity type       entity"))
	for _, inv := range invitations {
		entityType, entityName := entityTypeAndName(inv)
		fmt.Printf("%s   %s   %s\n", inv.Metadata.GUID, entityType, entityName)
	}
	return nil
}

// getAllInvitations retrieves the invitations for an entity type.
func getAllInvitations(c plugin.CliConnection) ([]Invitation, error) {
	entityTypes := []string{
		"account",
		"organization",
		"space",
	}

	var invitations []Invitation
	for _, t := range entityTypes {
		url := fmt.Sprintf("/custom/%s_invitations", t)
		resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
		if err != nil {
			return []Invitation{}, fmt.Errorf("Couldn't retrieve %s invitations", t)
		}

		resString := strings.Join(resLines, "")
		var res InvitationsResponse
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return []Invitation{}, errors.New("Couldn't read JSON response from server")
		}

		if res.ErrorCode != "" {
			return []Invitation{}, errors.New(res.Description)
		}

		for _, i := range res.Resources {
			invitations = append(invitations, i)
		}
	}

	return invitations, nil
}

// entityTypeAndName returns the entity type and its name for an invitation.
func entityTypeAndName(inv Invitation) (string, string) {
	var entityType string
	var entityName string
	if inv.Entity.AccountID != "" {
		entityType = "Billing Account"
		entityName = inv.Entity.AccountName
	}
	if inv.Entity.OrganizationID != "" {
		entityType = "Org            "
		entityName = inv.Entity.OrganizationName
	}
	if inv.Entity.SpaceID != "" {
		entityType = "Space          "
		entityName = fmt.Sprintf("%s / %s", inv.Entity.OrganizationName, inv.Entity.SpaceName)
	}

	return entityType, entityName
}
