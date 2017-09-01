package main

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
)

// Domain is a Cloud Foundry domain.
type Domain struct {
	Metadata CFMetadata `json:"metadata"`
	Entity   struct {
		Name string `json:"name"`
	} `json:"entity"`
}

// DomainsResponse is a response from the server to a domains request.
type DomainsResponse struct {
	Resources []Domain `json:"resources"`
	ServerResponsePagination
	ServerResponseError
}

// getSharedDomains retrieves all shared domains.
func getSharedDomains(c plugin.CliConnection) ([]Domain, error) {
	url := "/v2/shared_domains"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return []Domain{}, err
	}

	resString := strings.Join(resLines, "")
	var res DomainsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return []Domain{}, errors.New("Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return []Domain{}, errors.New(res.Description)
	}

	return res.Resources, nil
}
