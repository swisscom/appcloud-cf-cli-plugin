package main

import (
	"encoding/json"

	"fmt"
	"strings"

	"github.com/pkg/errors"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// SSLEnabled tells the user whether SSL is enabled for a full domain name.
func (p *AppCloudPlugin) SSLEnabled(c plugin.CliConnection, domain string, hostname string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	fullDomain := domain
	if hostname != "" {
		fullDomain = strings.Join([]string{hostname, domain}, ".")
	}

	p.ui.Say("Checking SSL status for %s as %s...", terminal.EntityNameColor(fullDomain), terminal.EntityNameColor(un))

	s, err := c.GetCurrentSpace()
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve current space")
	}

	url := fmt.Sprintf("/custom/spaces/%s/certificates", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve SSL certificates for space")
	}

	resString := strings.Join(resLines, "")
	var res SSLCertificatesResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	var enabled bool
	for _, cert := range res.Resources {
		if cert.Entity.FullDomainName == fullDomain {
			enabled = true
			break
		}
	}

	if !enabled {
		sharedDomains, err := getSharedDomains(c)
		if err != nil {
			return errors.Wrap(err, "Couldn't get shared domains")
		}

		for _, d := range sharedDomains {
			if d.Entity.Name == domain {
				enabled = true
				break
			}
		}
	}

	if enabled {
		p.ui.Say("SSL is enabled for '%s'", fullDomain)
	} else {
		p.ui.Say("SSL is not enabled for '%s'", fullDomain)
	}

	return nil
}
