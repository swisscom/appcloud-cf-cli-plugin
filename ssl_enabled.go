package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// SSLEnabled tells the user whether SSL is enabled for a full domain name.
func (p *AppCloudPlugin) SSLEnabled(c plugin.CliConnection, domain string, hostname string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fullDomain := domain
	if hostname != "" {
		fullDomain = strings.Join([]string{hostname, domain}, ".")
	}

	fmt.Printf("Checking SSL status for %s as %s...\n", cyanBold(fullDomain), cyanBold(username))

	s, err := c.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve current space")
	}

	if s.Guid == "" {
		return fmt.Errorf("No space targeted, use %s to target a space", cyanBold("'cf target -s SPACE'"))
	}

	url := fmt.Sprintf("/custom/spaces/%s/certificates", s.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve SSL certificates for space %s", s.Name)
	}

	resString := strings.Join(resLines, "")
	var res SSLCertificatesResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Print(greenBold("OK\n\n"))
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
			return errors.New("Couldn't get shared domains")
		}

		for _, d := range sharedDomains {
			if d.Entity.Name == domain {
				enabled = true
				break
			}
		}
	}

	if enabled {
		fmt.Printf("SSL is enabled for '%s'\n", fullDomain)
	} else {
		fmt.Printf("SSL is not enabled for '%s'\n", fullDomain)
	}
	return nil
}
