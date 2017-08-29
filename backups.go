package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// BackupsResponse is the response of the server to a get backups call
type BackupsResponse struct {
	ServiceBrokerResponse
	TotalResult int      `json:"total_results"`
	TotalPages  int      `json:"total_pages"`
	PrevURL     string   `json:"prev_url"`
	NextURL     string   `json:"next_url"`
	Resources   []Backup `json:"resources"`
}

// Backups lists all backups for a service instance
func (p *AppCloudPlugin) Backups(c plugin.CliConnection, serviceInstanceName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Getting backups for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve service instance %s", serviceInstanceName)
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve backups for %s", serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var bRes BackupsResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	backups := bRes.Resources
	if len(backups) == 0 {
		fmt.Println("No backups found")
		return nil
	}

	fmt.Println(bold("     GUID                                   created at             last operation"))
	for i, b := range backups {
		fmt.Printf("#%v   %s   %s   %s\n", i, b.Metadata.GUID, b.Metadata.CreatedAt, formatStatus(b.Entity.Status))
	}
	return nil
}

// formatStatus formats a status more nicely
func formatStatus(s string) string {
	formatted := strings.Replace(s, "_", " ", -1)
	formatted = strings.ToLower(formatted)

	return formatted
}
