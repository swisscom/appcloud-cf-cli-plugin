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
	fmt.Println(greenBold("OK\n\n"))
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}

	url := fmt.Sprintf("/custom/spaces/%s/certificates",s.SpaceFields.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "GET", url)

	if err != nil {
		return fmt.Errorf("Couldn't list all SSL certificate for given space:  %s", s.SpaceFields.Name)
	}

	resString := strings.Join(resLines, "")
	var bRes ListSSLCertResponse
	err = json.Unmarshal([]byte(resString), &bRes)
	if err != nil {
		return fmt.Errorf("Couldn't read JSON response: %s", err.Error())
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}
	fmt.Println("===============================================================================================================================")
	fmt.Println(greenBold("SSL Certificate                    Status           not_valid_before           not_valid_after                automatic_renewal"))
	fmt.Println("===============================================================================================================================")
	for i := 0; i < len(bRes.Resources); i++ {
		fmt.Println(bRes.Resources[i].Entity.FullDomainName+"          "+bRes.Resources[i].Entity.Status+"     "+bRes.Resources[i].Entity.NotValidBefore+"       "+bRes.Resources[i].Entity.NotValidAfter+" "+bRes.Resources[i].Entity.AutomaticRenewal);
	}
	if(len(bRes.Resources)==0){
		fmt.Println("No certificates found")
	}
	fmt.Println("===============================================================================================================================")

	fmt.Println("Available SSL certificates retrieved successfully")
	return nil
}
