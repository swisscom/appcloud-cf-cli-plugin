package main

import (
	"encoding/json"
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
func (p *AppCloudPlugin) CreateBackup(c plugin.CliConnection, serviceInstanceName string) {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Creating backup for service instance %s as %s...\n", cyanBold(serviceInstanceName), cyanBold(username))

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve service instance %s\n", cyanBold(serviceInstanceName))
		return
	}

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve backups for %s\n", cyanBold(serviceInstanceName))
		return
	}

	resString := strings.Join(resLines, "")
	var bRes CreateBackupResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		fmt.Printf("\nCouldn't read JSON response\n")
		return
	}

	if bRes.ErrorCode != "" {
		fmt.Printf("\n%s\n", bRes.Description)
		return
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Backup in progress")
}
