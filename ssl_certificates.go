package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// SSLCertificates Lists available SSL certificates
func (p *AppCloudPlugin) SSLCertificates(c plugin.CliConnection) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	o, oerr := c.GetCurrentOrg()
	if oerr != nil {
		return fmt.Errorf("Couldn't retrieve current org %s", oerr.Error())
	}
	if o.Guid == "" {
		return fmt.Errorf("No org and space targeted, use %s to target an org and space", cyanBold("'cf target -o ORG -s SPACE'"))
	}

	s, err := c.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve current space %s", err.Error())
	}

	if s.Guid == "" {
		return fmt.Errorf("No space targeted, use %s to target a space", cyanBold("'cf target -s SPACE'"))
	}

	fmt.Printf("Getting SSL certificates for space %s as %s...\n", cyanBold(s.Name), cyanBold(username))

	url := fmt.Sprintf("/custom/spaces/%s/certificates", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't get certificates for space %s", s.Name)
	}

	resString := strings.Join(resLines, "")
	var res SSLCertificatesResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return fmt.Errorf("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Println(greenBold("OK\n\n"))

	if len(res.Resources) > 0 {
		table := NewTable([]string{"full domain name", "status"})
		for _, cert := range res.Resources {
			table.Add(cert.Entity.FullDomainName, formatStatus(cert.Entity.Status))
		}
		table.Print()
	} else {
		fmt.Println("No SSL certificates found")
	}

	return nil
}
