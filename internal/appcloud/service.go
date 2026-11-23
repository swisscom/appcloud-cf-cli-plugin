package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/v8/plugin"
	"github.com/pkg/errors"
)

type V3Service struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

type V3ServicesResponse struct {
	Resources []V3Service `json:"resources"`
	// ServerResponsePagination
	V3ServerResponseErrors
}

func getService(c plugin.CliConnection, name string) (V3Service, error) {
	space, err := c.GetCurrentSpace()
	if err != nil {
		return V3Service{}, err
	}

	url := fmt.Sprintf("/v3/service_instances?names=%s&space_guids=%s", name, space.Guid)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return V3Service{}, err
	}

	resString := strings.Join(resLines, "")
	var res V3ServicesResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return V3Service{}, errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if len(res.Errors) > 0 {
		return V3Service{}, fmt.Errorf("Error response from server: %s", res.Errors[0].Detail)
	}
	if len(res.Resources) == 0 {
		return V3Service{}, fmt.Errorf("Service '%s' not found", name)
	}

	return res.Resources[0], nil
}
