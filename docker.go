package main

import (
	"fmt"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/cf/errors"
	"strings"
	"encoding/json"
)

func (p *AppCloudPlugin) DockerRepository(c plugin.CliConnection, org string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	if org == "none" {
		organization, err := c.GetCurrentOrg()
		if err != nil {
			return errors.New("Could not get current organization")
		}

		org = organization.Guid
		fmt.Println(organization.Guid)
	}

	fmt.Printf("\nRetrieving your repositories for %s as %s...\n\n", org, cyanBold(username))

	url := fmt.Sprintf("/custom/organizations/%s/docker-repositories", org)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.New("Couldn't retrieve repositories for " + org)
	}

	resString := strings.Join(resLines, "")

	var dRes DockerRepository
	err = json.Unmarshal([]byte(resString), &dRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	output := ""
	for i := 0; i < len(dRes.Repositories); i++ {
		output += dRes.Repositories[i]
	}

	if output != "" {
		fmt.Println(bold("Docker Repositories\n"))
		fmt.Println(output)
	} else {
		fmt.Println("There are no docker repositories for this organization.")
	}

	return nil
}
