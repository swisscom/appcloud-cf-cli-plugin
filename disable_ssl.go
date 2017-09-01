package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// DisableSSL uninstalls an existing SSL certificate.
func (p *AppCloudPlugin) DisableSSL(c plugin.CliConnection, domain string, hostname string) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fullDomain := domain
	if hostname != "" {
		fullDomain = strings.Join([]string{hostname, domain}, ".")
	}

	fmt.Printf("Disabling SSL for %s as %s...\n", cyanBold(fullDomain), cyanBold(username))

	s, err := c.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve current space")
	}

	req := SSLCertificateRequest{
		SpaceID:        s.Guid,
		FullDomainName: fullDomain,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return errors.New("Couldn't parse JSON data")
	}
	url := "/custom/certifications/uninstall"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "PUT", url, "-d", string(reqData))

	if err != nil {
		return fmt.Errorf("Couldn't disable SSL for %s", fullDomain)
	}

	resString := strings.Join(resLines, "")
	var res SSLCertificateResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return errors.New(res.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("SSL disabled")
	return nil
}
