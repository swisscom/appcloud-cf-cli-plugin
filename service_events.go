package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
)

// ServiceEvents retrieves events for a service instance.
func (p *AppCloudPlugin) ServiceEvents(c plugin.CliConnection, serviceInstanceName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Getting events for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Service instance %s not found", serviceInstanceName)
	}

	url := fmt.Sprintf("/v2/events?q=actee:%s", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve events for %s", serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var res EventsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	events := res.Resources
	if len(events) > 0 {
		table := NewTable([]string{"time", "event", "actor"})
		for _, e := range events {
			table.Add(e.Metadata.CreatedAt.Format(time.RFC3339), e.Entity.Type, e.Entity.ActorName)
		}
		table.Print()
	} else {
		fmt.Println("No events found")
	}
	return nil
}
