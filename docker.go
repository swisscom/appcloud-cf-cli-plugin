package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
)

func (p *AppCloudPlugin) DockerRepository(c plugin.CliConnection, org string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	orgId := ""

	if org == "none" {
		organization, err := c.GetCurrentOrg()
		if err != nil {
			return errors.New("Could not get current organization")
		}
		orgId = organization.Guid
		org = organization.Name
	} else {
		curOrg, err := c.GetOrg(org)
		if err != nil {
			return errors.New("Could not get organization")
		}

		orgId = curOrg.Guid
	}

	if org == "" {
		fmt.Println("No organization set.")
		return nil
	}

	fmt.Printf("\nRetrieving your repositories for %s as %s...\n\n", org, cyanBold(username))

	url := fmt.Sprintf("/custom/organizations/%s/docker-repositories", orgId)
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
		output += dRes.Repositories[i] + "\n"
	}

	if output != "" {
		fmt.Println(bold("Docker Repositories\n"))
		fmt.Println(output)
	} else {
		fmt.Println("There are no docker repositories for this organization.")
	}

	fmt.Print(greenBold("OK\n\n"))

	return nil
}
