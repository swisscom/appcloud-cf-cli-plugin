package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
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
func (p *AppCloudPlugin) Tree(c plugin.CliConnection, depth int) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	p.ui.Say("Getting organization tree as %s...", terminal.EntityNameColor(un))

	url := "/custom/organizations"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve organization tree")
	}

	resString := strings.Join(resLines, "")
	var res TreeResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	orgs := res.Resources
	if len(orgs) == 0 {
		p.ui.Say("No organizations found")
		return nil
	}

	renderTree(p.ui, orgs, depth)
	return nil
}

// renderTree renders the org tree with a specified level of depth.
func renderTree(ui terminal.UI, orgs []TreeOrg, depth int) {
	output := terminal.HeaderColor("Orgs\n")

	for i, o := range orgs {
		lastOrg := i == len(orgs)-1

		if lastOrg {
			output += fmt.Sprintf("└─ %s\n", o.Name)
		} else {
			output += fmt.Sprintf("├─ %s\n", o.Name)
		}

		if len(o.Spaces) > 0 && depth > 0 {
			if lastOrg {
				output += "    " + terminal.HeaderColor("Spaces\n")
			} else {
				output += "│   " + terminal.HeaderColor("Spaces\n")
			}

			for j, s := range o.Spaces {
				lastSpace := j == len(o.Spaces)-1

				if lastOrg {
					output += "    "
				} else {
					output += "│   "
				}

				if lastSpace {
					output += fmt.Sprintf("└─ %s\n", s.Name)
				} else {
					output += fmt.Sprintf("├─ %s\n", s.Name)
				}

				if len(s.Applications) > 0 && depth > 1 {
					if lastOrg {
						output += "    "
					} else {
						output += "│   "
					}

					if lastSpace {
						output += "    " + terminal.HeaderColor("Apps\n")
					} else {
						output += "│   " + terminal.HeaderColor("Apps\n")
					}

					for k, a := range s.Applications {
						lastApp := k == len(s.Applications)-1

						if lastOrg {
							output += "    "
						} else {
							output += "│   "
						}

						if lastSpace {
							output += "    "
						} else {
							output += "│   "
						}

						if lastApp {
							output += fmt.Sprintf("└─ %s\n", a.Name)
						} else {
							output += fmt.Sprintf("├─ %s\n", a.Name)
						}
					}
				}

				if len(s.ServiceInstances) > 0 && depth > 1 {
					if lastOrg {
						output += "    "
					} else {
						output += "│   "
					}

					if lastSpace {
						output += "    " + terminal.HeaderColor("Services\n")
					} else {
						output += "│   " + terminal.HeaderColor("Services\n")
					}

					for k, si := range s.ServiceInstances {
						lastService := k == len(s.ServiceInstances)-1

						if lastOrg {
							output += "    "
						} else {
							output += "│   "
						}

						if lastSpace {
							output += "    "
						} else {
							output += "│   "
						}

						if lastService {
							output += fmt.Sprintf("└─ %s\n", si.Name)
						} else {
							output += fmt.Sprintf("├─ %s\n", si.Name)
						}
					}
				}
			}
		}
	}

	ui.Say(output)
}
