package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// ListSSLCertResponse is the response from the server when retrieved all certificates 
type ListSSLCertResponse struct {
	SSLCertificates
	ServerResponse
}

// ListSSLCertificates Lists available SSL certificates
func (p *AppCloudPlugin) ListSSLCertificates(c plugin.CliConnection) error {
	// Get the current targeted space details 
	s, err := c.GetCurrentSpace()
	fmt.Printf("Listing SSL certificates for provided space %s ...\n", cyanBold(s.Name))
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}

	url := fmt.Sprintf("/custom/spaces/%s/certificates",s.SpaceFields.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "GET", url)
	fmt.Println("response:",resLines);
	if err != nil {
		return fmt.Errorf("Couldn't list all SSL certificate for given space:  %s", s.SpaceFields.Name)
	}

	resString := strings.Join(resLines, "")
	var bRes ListSSLCertResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	fmt.Print(greenBold("OK\n\n"))

	fmt.Println("Available SSL certificates retrieved successfully")
	return nil
}
