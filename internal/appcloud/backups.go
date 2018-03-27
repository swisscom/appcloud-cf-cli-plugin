package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// Backups lists all backups for a service instance.
func (p *Plugin) Backups(c plugin.CliConnection, serviceInstanceName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Getting backups for service instance %s as %s...", terminal.EntityNameColor(serviceInstanceName), terminal.EntityNameColor(un))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return errors.Wrap(err, "Service instance not found")
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve backups")
	}

	resString := strings.Join(resLines, "")
	var res BackupsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return fmt.Errorf("Error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	backups := res.Resources
	if len(backups) > 0 {
		table := p.ui.Table([]string{"created at", "GUID", "last operation"})
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
		err := table.Print()
		if err != nil {
			return errors.Wrap(err, "Couldn't print table")
		}
	} else {
		p.ui.Say("No backups found")
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
