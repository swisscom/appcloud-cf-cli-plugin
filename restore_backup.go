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
	Restore
	ServerResponseError
}

// RestoreBackup creates a backup for a service instance
func (p *AppCloudPlugin) RestoreBackup(c plugin.CliConnection, serviceInstanceName string, backupGUID string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Restoring backup on the service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Service instance %s not found", serviceInstanceName)
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups/%s/restores", s.Guid, backupGUID)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		return fmt.Errorf("Couldn't restore %s on %s", backupGUID, serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var bRes RestoreBackup
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Printf("Restore in progress. Use '%s' to check operation status.\n", yellowBold("cf backups"))
	return nil
}
