package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/v8/plugin"
	"github.com/pkg/errors"
)

type V3Space struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

type V3SpacesResponse struct {
	Resources []V3Space `json:"resources"`
	V3ServerResponseErrors
}

func getSpace(c plugin.CliConnection, name string) (V3Space, error) {
	org, err := c.GetCurrentOrg()
	if err != nil {
		return V3Space{}, err
	}

	url := fmt.Sprintf("/v3/spaces?names=%s&organization_guids=%s", name, org.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return V3Space{}, err
	}

	resString := strings.Join(resLines, "")
	var res V3SpacesResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return V3Space{}, errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if len(res.Errors) > 0 {
		return V3Space{}, fmt.Errorf("Error response from server: %s", res.Errors[0].Detail)
	}
	if len(res.Resources) == 0 {
		return V3Space{}, fmt.Errorf("Space '%s' not found", name)
	}

	return res.Resources[0], nil
}
