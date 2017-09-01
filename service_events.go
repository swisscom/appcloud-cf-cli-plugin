package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	if len(events) == 0 {
		fmt.Println("No events found")
		return nil
	}

	fmt.Println(bold("time                   event"))
	for _, e := range events {
		fmt.Printf("%s   %s   %s\n", cyanBold(e.Metadata.CreatedAt), e.Entity.Type, e.Entity.ActorName)
	}
	return nil
}
