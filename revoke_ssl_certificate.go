package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// RevokeSSLCertResponse is the response from the server from a certificate revokation call
type RevokeSSLCertResponse struct {
	SSLCertificate
	ServerResponse
}

// RevokeSSLCertificate revokes given SSL certificate
func (p *AppCloudPlugin) RevokeSSLCertificate(c plugin.CliConnection, fullDomain string) error {
	fmt.Printf("Revoking SSL certificate for route %s ...\n", cyanBold(fullDomain))
	fmt.Print(greenBold("OK\n\n"))
	// Get the current targeted space details 
	s, err := c.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}
	req := fmt.Sprintf("'{\"space_id\": \"%s\",\"full_domain_name\": \"%s\"}'", s.Guid, fullDomain)

	//req :=  "'{\"space_id\": \""+s.SpaceFields.Guid+"\","+"\"full_domain_name\":\""+ fullDomain+ "\"}'"

	url := "/custom/certifications/revoke"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "PUT","-d",req, url)

	if err != nil {
		return fmt.Errorf("Couldn't revoke SSL certificate for route:  %s", fullDomain)
	}

	resString := strings.Join(resLines, "")
	var bRes InstallSSLCertResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}

	
	fmt.Println("SSL certificate revoked successfully")
	return nil
}
