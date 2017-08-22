package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateBackup creates a backup for a service instance
func (p *AppCloudPlugin) CreateBackup(c plugin.CliConnection, serviceInstanceName string) {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("Creating backup for service instance %s as %s...\n", serviceInstanceName, username)

	s, err := c.GetService(serviceInstanceName)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve service instance %s\n", serviceInstanceName)
		return
	}

	fmt.Println("This may take a while...")

	url := fmt.Sprintf("/custom/service_instances/%s/backups", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "POST", url)
	if err != nil {
		fmt.Printf("\nCouldn't retrieve backups for %s\n", serviceInstanceName)
		return
	}

	fmt.Print("OK\n\n")
	fmt.Printf("Backup for service instance %s created\n", serviceInstanceName)
	fmt.Println(resLines)
}
