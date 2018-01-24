package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// SSLCertificates Lists available SSL certificates
func (p *AppCloudPlugin) SSLCertificates(c plugin.CliConnection) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	s, err := c.GetCurrentSpace()
	if err != nil {
		return errors.Wrap(err, "Couldn't get current space")
	}

	p.ui.Say("Getting SSL certificates for space %s as %s...", terminal.EntityNameColor(s.Name), terminal.EntityNameColor(un))

	url := fmt.Sprintf("/custom/spaces/%s/certificates", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't get SSL certificates")
	}

	resString := strings.Join(resLines, "")
	var res SSLCertificatesResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return fmt.Errorf("Error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	if len(res.Resources) > 0 {
		table := p.ui.Table([]string{"full domain name", "status"})
		for _, cert := range res.Resources {
			table.Add(cert.Entity.FullDomainName, formatStatus(cert.Entity.Status))
		}
		err := table.Print()
		if err != nil {
			return errors.Wrap(err, "Couldn't print table")
		}
	} else {
		p.ui.Say("No SSL certificates found")
	}

	return nil
}
