package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// V3Domain is a Cloud Foundry domain.
type V3Domain struct {
	Name string `json:"name"`
}

// V3DomainsResponse is a response from the server to a domains request.
type V3DomainsResponse struct {
	Resources []V3Domain `json:"resources"`
	// ServerResponsePagination
	V3ServerResponseErrors
}

// getOrgDomains retrieves all org domains.
func getOrgDomains(c plugin.CliConnection) ([]V3Domain, error) {
	org, err := c.GetCurrentOrg()
	if err != nil {
		return []V3Domain{}, err
	}

	url := fmt.Sprintf("/v3/organizations/%s/domains", org.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return []V3Domain{}, err
	}

	resString := strings.Join(resLines, "")
	var res V3DomainsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return []V3Domain{}, errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if len(res.Errors) > 0 {
		return []V3Domain{}, fmt.Errorf("Error response from server: %s", res.Errors[0].Detail)
	}

	return res.Resources, nil
}
