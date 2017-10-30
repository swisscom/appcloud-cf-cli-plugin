package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

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

	if len(invitations) > 0 {
		table := NewTable([]string{"GUID", "entity type", "entity"})
		for _, inv := range invitations {
			entityType, entityName := invitationEntityTypeAndName(inv)
			table.Add(inv.Metadata.GUID, formatEntityType(entityType), entityName)
		}
		table.Print()
	} else {
		fmt.Println("No invitations found")
	}
	return nil
}

// formatEntityType formats an entity name more nicely.
func formatEntityType(t string) string {
	switch t {
	case "account":
		return "Billing Account"
	case "organization":
		return "Org            "
	case "space":
		return "Space          "
	default:
		return "unknown        "
	}
}
