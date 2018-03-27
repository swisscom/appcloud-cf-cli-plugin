package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// CreateBackup creates a backup for a service instance.
func (p *Plugin) CreateBackup(c plugin.CliConnection, serviceInstanceName string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Creating backup for service instance %s as %s...", terminal.EntityNameColor(serviceInstanceName), terminal.EntityNameColor(un))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return errors.Wrap(err, "Service instance not found")
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't create backup for service instance")
	}

	resString := strings.Join(resLines, "")
	var res BackupResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return fmt.Errorf("Error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	p.ui.Say("Create in progress. Use '%s' to check operation status.", terminal.CommandColor("cf backups"))

	return nil
}
