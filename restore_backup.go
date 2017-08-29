package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// RestoreBackup is a service instance restore
type RestoreBackup struct {
	ServiceBrokerResponse
	Restore
}

// RestoreBackup creates a backup for a service instance
func (p *AppCloudPlugin) RestoreBackup(c plugin.CliConnection, serviceInstanceName string, backupGUID string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Creating a restore for the service backup for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve service instance %s", serviceInstanceName)
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups/%s/restores", s.Guid, backupGUID)
	fmt.Print(url)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve the restore %s for %s", backupGUID, serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var bRes RestoreBackup
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Restore in progress")
	return nil
}
