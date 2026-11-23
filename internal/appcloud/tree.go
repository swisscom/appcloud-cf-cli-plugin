package appcloud

import (
	"code.cloudfoundry.org/cli/v8/plugin"
)

// These are strings used to create the tree structure.
const (
	indent           = "    "
	indentWithBranch = "│   "
)

// TreeResponse is the response of a server to a tree request.
type TreeResponse struct {
	Resources []TreeOrg `json:"resources"`
	ServerResponsePagination
	ServerResponseError
}

// TreeOrg is an org node of the tree structure.
type TreeOrg struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Status string      `json:"status"`
	Spaces []TreeSpace `json:"spaces"`
}

// TreeSpace is a space node of the tree structure.
type TreeSpace struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Applications     []TreeApp     `json:"applications"`
	ServiceInstances []TreeService `json:"service_instances"`
}

// TreeApp is an app node of the tree structure.
type TreeApp struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	BackupIconURL string `json:"buildpack_icon_url"`
}

// TreeService is a service instance node of the tree structure.
type TreeService struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ServiceIconURL string `json:"service_icon_url"`
}

// Tree renders the org tree for the current user.
func (p *Plugin) Tree(c plugin.CliConnection, depth int) error {
	p.ui.Say("This command is no longer available")
	return nil
}
