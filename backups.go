package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// Backups lists all backups for a service instance
func (p *AppCloudPlugin) Backups(c plugin.CliConnection, serviceInstanceName string) {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Getting backups for service instance %s as %s...\n", serviceInstanceName, username)

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve service instance %s\n", serviceInstanceName)
		return
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve backups for %s\n", serviceInstanceName)
		return
	}

	resString := strings.Join(resLines, "")

	var bRes BackupsResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		fmt.Printf("\nCouldn't read JSON response")
		return
	}

	fmt.Print("OK\n\n")

	backups := bRes.Resources
	if len(backups) == 0 {
		fmt.Println("No backups found")
		return
	}

	fmt.Println("     created at             last operation")
	for i, b := range backups {
		fmt.Printf("#%v   %s   %s\n", i, b.Metadata.CreatedAt, formatStatus(b.Entity.Status))
	}
}

// formatStatus formats a status more nicely
func formatStatus(s string) string {
	switch s {
	case "CREATE_SUCCEEDED":
		return "create succeeded"
	}

	return s
}
