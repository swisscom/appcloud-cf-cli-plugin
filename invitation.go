package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// Invitation is an invitation a user received to join a specific entity.
type Invitation struct {
	Metadata CFMetadata `json:"metadata"`
	Entity   struct {
		Invitee          string   `json:"invitee"`
		Roles            []string `json:"roles"`
		ActorUsername    string   `json:"actor_username"`
		ActorUserID      string   `json:"actor_user_id"`
		AccountID        string   `json:"account_id"`
		AccountName      string   `json:"account_name"`
		OrganizationID   string   `json:"organization_id"`
		OrganizationName string   `json:"organization_name"`
		SpaceID          string   `json:"space_id"`
		SpaceName        string   `json:"space_name"`
		Status           string   `json:"status"`
	} `json:"entity"`
}

// InvitationResponse is a response from the server to an invitation request.
type InvitationResponse struct {
	Invitation
	ServerResponseError
}

// InvitationsResponse is a response from the server containing invitations.
type InvitationsResponse struct {
	Resources []Invitation `json:"resources"`
	ServerResponsePagination
	ServerResponseError
}

// invitationEntityTypeAndName returns the entity type and its name for an invitation.
func invitationEntityTypeAndName(inv Invitation) (string, string) {
	var entityType string
	var entityName string
	if inv.Entity.AccountID != "" {
		entityType = "account"
		entityName = inv.Entity.AccountName
	}
	if inv.Entity.OrganizationID != "" {
		entityType = "organization"
		entityName = inv.Entity.OrganizationName
	}
	if inv.Entity.SpaceID != "" {
		entityType = "space"
		entityName = fmt.Sprintf("%s / %s", inv.Entity.OrganizationName, inv.Entity.SpaceName)
	}

	return entityType, entityName
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
			return []Invitation{}, errors.Wrap(err, "Couldn't retrieve invitations")
		}

		resString := strings.Join(resLines, "")
		var res InvitationsResponse
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			return []Invitation{}, errors.Wrap(err, "Couldn't read JSON response from server")
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
