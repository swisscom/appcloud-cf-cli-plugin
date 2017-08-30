package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// AbortSSLCertProcessResponse is the response from the server from a certificate aborting process call
type AbortSSLCertProcessResponse struct {
	ServerResponseError
	SSLCertificate
}

// AbortSSLCertificateProcess: A running certificate creation will be interrupted and aborted
func (p *AppCloudPlugin) AbortSSLCertificateProcess(c plugin.CliConnection, fullDomain string) error {
	fmt.Printf("Aborting running SSL certificate creation process for route %s ...\n", cyanBold(fullDomain))
	fmt.Print(greenBold("OK\n\n"))
	// Get the current targeted space details
	s, err := c.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}

	req := "'{\"space_id\": \"" + s.SpaceFields.Guid + "\"," + "\"full_domain_name\":\"" + fullDomain + "\"}'"

	url := "/custom/certifications/abort"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "PUT", "-d", req, url)

	if err != nil {
		return fmt.Errorf("Couldn't abort running SSL certificate creation process for route:  %s", fullDomain)
	}

	resString := strings.Join(resLines, "")
	var bRes AbortSSLCertProcessResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Println("Aborting SSL certificate creation process completed successfully")
	return nil
}
