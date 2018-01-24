package main

import (
	"encoding/json"

	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// ServiceEvents retrieves events for a service instance.
func (p *AppCloudPlugin) ServiceEvents(c plugin.CliConnection, serviceInstanceName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Getting events for service instance %s as %s...", terminal.EntityNameColor(serviceInstanceName), terminal.EntityNameColor(un))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return errors.Wrap(err, "Service instance not found")
	}

	url := fmt.Sprintf("/v2/events?q=actee:%s", s.Guid)
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

	if res.ErrorCode != "" {
		return fmt.Errorf("Error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	events := res.Resources
	if len(events) > 0 {
		table := p.ui.Table([]string{"time", "event", "actor"})
		for _, e := range events {
			table.Add(e.Metadata.CreatedAt.Format(time.RFC3339), e.Entity.Type, e.Entity.ActorName)
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
