package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// UnInstallSSLCertificate is the response from the server from a certificate uninstallation call
type UnInstallSSLCertResponse struct {
	SSLCertificate
	ServerResponse
}

// TurnSSLOff uninstalls given SSL certificate
func (p *AppCloudPlugin) TurnSSLOff(c plugin.CliConnection, fullDomain string) error {
	fmt.Printf("Uninstalling SSL certificate for route %s ...\n", cyanBold(fullDomain))
	fmt.Print(greenBold("OK\n\n"))
	// Get the current targeted space details
	s, err := c.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}
	req := fmt.Sprintf("'{\"space_id\": \"%s\",\"full_domain_name\": \"%s\"}'", s.Guid, fullDomain)

	url := "/custom/certifications/uninstall"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "PUT", "-d", req, url)

	if err != nil {
		return fmt.Errorf("Couldn't uninstall SSL certificate for route:  %s", fullDomain)
	}

	resString := strings.Join(resLines, "")
	var bRes UnInstallSSLCertResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Println("SSL certificate uninstalled successfully")
	return nil
}
