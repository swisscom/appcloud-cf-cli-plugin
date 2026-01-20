package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// SSLEnabled tells the user whether SSL is enabled for a full domain name.
func (p *Plugin) SSLEnabled(c plugin.CliConnection, domain string, hostname string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	fqdn := domain
	if hostname != "" {
		fqdn = strings.Join([]string{hostname, domain}, ".")
	}

	p.ui.Say("Checking SSL status for %s as %s...", terminal.EntityNameColor(fqdn), terminal.EntityNameColor(un))

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
		return fmt.Errorf("error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	var enabled bool
	for _, cert := range res.Resources {
		if cert.Entity.FullDomainName == fqdn {
			enabled = true
			break
		}
	}

	if !enabled {
		sharedDomains, err := getOrgDomains(c)
		if err != nil {
			return errors.Wrap(err, "Couldn't get org domains")
		}

		for _, d := range sharedDomains {
			if d.Name == domain {
				enabled = true
				break
			}
		}
	}

	if enabled {
		p.ui.Say("SSL is enabled for '%s'", fqdn)
	} else {
		p.ui.Say("SSL is not enabled for '%s'", fqdn)
	}

	return nil
}
