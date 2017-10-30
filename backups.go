package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
)

// Backups lists all backups for a service instance.
func (p *AppCloudPlugin) Backups(c plugin.CliConnection, serviceInstanceName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Getting backups for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Service instance %s not found", serviceInstanceName)
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve backups for %s", serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var res BackupsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	backups := res.Resources
	if len(backups) > 0 {
		table := NewTable([]string{"created at", "GUID", "last operation"})
		for _, b := range backups {
			var newestRestoreDate time.Time
			overallStatus := b.Entity.Status

			for _, r := range b.Entity.Restores {
				if r.Metadata.CreatedAt.After(newestRestoreDate) {
					overallStatus = fmt.Sprintf("RESTORE %s", r.Entity.Status)
					newestRestoreDate = r.Metadata.CreatedAt
				}
			}

			table.Add(b.Metadata.CreatedAt.Format(time.RFC3339), b.Metadata.GUID, formatStatus(overallStatus))
		}
		table.Print()
	} else {
		fmt.Println("No backups found")
	}
	return nil
}

// formatStatus formats a status more nicely.
func formatStatus(s string) string {
	if s == "VALID_INST" {
		s = "INSTALLED"
	}
	if s == "VALID_UNINST" {
		s = "UNINSTALLED"
	}

	formatted := strings.Replace(s, "_", " ", -1)
	formatted = strings.ToLower(formatted)

	return formatted
}
