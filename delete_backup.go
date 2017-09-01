package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// DeleteBackup deletes a backup of a service instance.
func (p *AppCloudPlugin) DeleteBackup(c plugin.CliConnection, serviceInstanceName string, backupGUID string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Deleting backup of service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		return fmt.Errorf("Service instance %s not found", serviceInstanceName)
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups/%s", s.Guid, backupGUID)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "DELETE", url)
	if err != nil {
		return fmt.Errorf("Couldn't delete %s of %s", backupGUID, serviceInstanceName)
	}

	resString := strings.Join(resLines, "")
	var res ServerResponseError
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Service backup deleted")
	return nil
}
