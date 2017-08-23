package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateBackupResponse is the response from the server from a create backup call
type CreateBackupResponse struct {
	ServiceBrokerResponse
	Backup
}

// CreateBackup creates a backup for a service instance
func (p *AppCloudPlugin) CreateBackup(c plugin.CliConnection, serviceInstanceName string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Creating backup for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve service instance %s", serviceInstanceName)
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve backups for %s", serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var bRes CreateBackupResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Backup in progress")
	return nil
}
