package appcloud

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// Invitations retrieves a user's invitations.
func (p *Plugin) Invitations(c plugin.CliConnection) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Getting invitations as %s...", terminal.EntityNameColor(un))

	invitations, err := getAllInvitations(c)
	if err != nil {
		return err
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	if len(invitations) > 0 {
		table := p.ui.Table([]string{"GUID", "entity type", "entity"})
		for _, inv := range invitations {
			entityType, entityName := invitationEntityTypeAndName(inv)
			table.Add(inv.Metadata.GUID, formatEntityType(entityType), entityName)
		}
		err := table.Print()
		if err != nil {
			return errors.Wrap(err, "Couldn't print table")
		}
	} else {
		p.ui.Say("No invitations found")
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
