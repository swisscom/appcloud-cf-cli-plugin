package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateSSLCertResponse is the response from the server from a create certificate call
type CreateSSLCertResponse struct {
	SSLCertificate
	ServerResponse
}

// CreateSSLCertificate creates a SSL certificate for provided route
func (p *AppCloudPlugin) CreateSSLCertificate(c plugin.CliConnection, fullDomain string) error {
	fmt.Printf("Creating SSL certificate for route %s ...\n", cyanBold(fullDomain))
	fmt.Print(greenBold("OK\n\n"))
	// Get the current targeted space details 
	s, err := c.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}
	req := fmt.Sprintf("'{\"space_id\": \"%s\",\"full_domain_name\": \"%s\"}'", s.Guid, fullDomain)

	url := "/custom/certifications/create"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "PUT","-d",req, url)

	if err != nil {
		return fmt.Errorf("Couldn't create SSL certificate for domain:  %s", fullDomain)
	}

	resString := strings.Join(resLines, "")
	var bRes CreateSSLCertResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}


	fmt.Println("SSL certificate creation suceeded")
	return nil
}
