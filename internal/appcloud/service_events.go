package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// ServiceEvents retrieves events for a service instance.
func (p *Plugin) ServiceEvents(c plugin.CliConnection, serviceInstanceName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Getting events for service instance %s as %s...", terminal.EntityNameColor(serviceInstanceName), terminal.EntityNameColor(un))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return errors.Wrap(err, "Service instance not found")
	}

	url := fmt.Sprintf("/v3/audit_events?target_guids=%s", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve events for service instance")
	}

	resString := strings.Join(resLines, "")
	var res EventsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if len(res.Errors) > 0 {
		return fmt.Errorf("error response from server: %s", res.Errors[0].Detail)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	events := res.Resources
	if len(events) > 0 {
		table := p.ui.Table([]string{"time", "event", "actor"})
		for _, e := range events {
			table.Add(e.CreatedAt, e.Type, e.Actor.Name)
		}
		err := table.Print()
		if err != nil {
			return errors.Wrap(err, "Couldn't print table")
		}
	} else {
		p.ui.Say("No events found")
	}
	return nil
}
