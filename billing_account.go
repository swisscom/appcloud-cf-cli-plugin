package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// BillingAccount is an entity of the Swisscom Application Cloud which handles billing.
type BillingAccount struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Status           string `json:"status"`
		StatusReasonData struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"status_reason_data"`
		Prohibited        bool   `json:"prohibited"`
		Name              string `json:"name"`
		CustomerNumber    string `json:"customer_number"`
		CustomerType      string `json:"customer_type"`
		OrganizationCount int    `json:"organization_count"`
		MaxOrganizations  int    `json:"max_organizations"`
	} `json:"entity"`
}

// BillingAccountResponse is a response from the server to a billing account request.
type BillingAccountResponse struct {
	BillingAccount
	ServerResponseError
}

// getBillingAccount retrieves a billing account by name.
func getBillingAccount(c plugin.CliConnection, name string) (BillingAccount, error) {
	accURL := url.QueryEscape("/custom/accounts?q=name:%s")
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", accURL)

	if err != nil {
		return BillingAccount{}, fmt.Errorf("Billing Account %s not found", name)
	}

	resString := strings.Join(resLines, "")
	var res BillingAccountResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return BillingAccount{}, errors.New("Couldn't read JSON response from server")
	}

	return res.BillingAccount, nil
}
