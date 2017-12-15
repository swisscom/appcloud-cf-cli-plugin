package main

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// DeleteBackup deletes a backup of a service instance.
func (p *AppCloudPlugin) DeleteBackup(c plugin.CliConnection, serviceInstanceName string, backupGUID string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Deleting backup of service instance %s as %s...", terminal.EntityNameColor(serviceInstanceName), terminal.EntityNameColor(un))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return errors.Wrap(err, "Service instances not found")
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups/%s", s.Guid, backupGUID)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "DELETE", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't delete backup")
	}

	resString := strings.Join(resLines, "")
	var res ServerResponseError
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	p.ui.Say("Service backup deleted")

	return nil
}
