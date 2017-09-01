package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateBackup creates a backup for a service instance.
func (p *AppCloudPlugin) CreateBackup(c plugin.CliConnection, serviceInstanceName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Creating backup for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Service instance %s not found", serviceInstanceName)
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		return fmt.Errorf("Couldn't create backup for %s", serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var res BackupResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Printf("Create in progress. Use '%s' to check operation status.\n", yellowBold("cf backups"))
	return nil
}
