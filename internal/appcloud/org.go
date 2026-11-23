package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/v8/plugin"
	"github.com/pkg/errors"
)

type V3Organization struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

type V3OrganizationsResponse struct {
	Resources []V3Organization `json:"resources"`
	// ServerResponsePagination
	V3ServerResponseErrors
}

func getOrg(c plugin.CliConnection, name string) (V3Organization, error) {
	url := fmt.Sprintf("/v3/organizations?names=%s", name)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return V3Organization{}, err
	}

	resString := strings.Join(resLines, "")
	var res V3OrganizationsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return V3Organization{}, errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if len(res.Errors) > 0 {
		return V3Organization{}, fmt.Errorf("Error response from server: %s", res.Errors[0].Detail)
	}
	if len(res.Resources) == 0 {
		return V3Organization{}, fmt.Errorf("Organization '%s' not found", name)
	}

	return res.Resources[0], nil
}
