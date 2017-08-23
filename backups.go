package main

import (
	"encoding/json"
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
func (p *AppCloudPlugin) Backups(c plugin.CliConnection, serviceInstanceName string) {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Getting backups for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve service instance %s\n", cyanBold(serviceInstanceName))
		return
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve backups for %s\n", cyanBold(serviceInstanceName))
		return
	}

	resString := strings.Join(resLines, "")
	var bRes BackupsResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		fmt.Printf("\nCouldn't read JSON response: %s\n")
		return
	}

	if bRes.ErrorCode != "" {
		fmt.Printf("\n%s\n", bRes.Description)
		return
	}

	fmt.Print(greenBold("OK\n\n"))

	backups := bRes.Resources
	if len(backups) == 0 {
		fmt.Println("No backups found")
		return
	}

	fmt.Println(bold("     created at             last operation"))
	for i, b := range backups {
		fmt.Printf("#%v   %s   %s\n", i, b.Metadata.CreatedAt, formatStatus(b.Entity.Status))
	}
}

// formatStatus formats a status more nicely
func formatStatus(s string) string {
	formatted := strings.Replace(s, "_", " ", -1)
	formatted = strings.ToLower(formatted)

	return formatted
}
