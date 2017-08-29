package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
)

func (p *AppCloudPlugin) Tree(c plugin.CliConnection, level int) error {
	username, err := c.Username()
	if err != nil {
		username = "you"
	}

	fmt.Printf("\nRetrieving your organization tree as %s...\n", cyanBold(username))

	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "/custom/organizations")
	if err != nil {
		return errors.New("Couldn't retrieve organizations")
	}

	resString := strings.Join(resLines, "")
	var oRes OrgResponse
	err = json.Unmarshal([]byte(resString), &oRes)
	if err != nil {
		return errors.New("Couldn't read JSON response")
	}

	fmt.Print(greenBold("OK\n\n"))

	orgs := oRes.Resources
	if len(orgs) == 0 {
		fmt.Println("No organizations found")
		return nil
	}

	TreeOutput(oRes, level)

	return nil
}

func TreeOutput(oRes OrgResponse, level int) {

	output := bold("Org Tree\n\n")
	output += bold("Organizations\n")

	for i := 0; i < oRes.TotalResults; i++ {

		if i+1 == oRes.TotalResults {
			output += "└"
		} else {
			output += "├"
		}

		output += fmt.Sprintf("─ %s\n", oRes.Resources[i].Name)

		if len(oRes.Resources[i].Spaces) > 0 && level > 1 {
			if i+1 == oRes.TotalResults {
				output += bold("    Spaces\n")
			} else {
				output += bold("|   Spaces\n")
			}
			for j := 0; j < len(oRes.Resources[i].Spaces); j++ {
				if j+1 == len(oRes.Resources[i].Spaces) && i+1 == oRes.TotalResults {
					if i+1 == oRes.TotalResults {
						output += "    └"
					} else {
						output += "│   └"
					}
				} else {
					output += "│   ├"
				}

				output += fmt.Sprintf("─ %s\n", oRes.Resources[i].Spaces[j].Name)

				if len(oRes.Resources[i].Spaces[j].Applications) > 0 && level > 2 {
					if j+1 == len(oRes.Resources[i].Spaces) && i+1 == oRes.TotalResults {
						output += bold("        Applications\n")
					} else {
						output += bold("|   |   Applications\n")
					}
					for k := 0; k < len(oRes.Resources[i].Spaces[j].Applications); k++ {
						if k+1 == len(oRes.Resources[i].Spaces[j].Applications) {
							if j+1 == len(oRes.Resources[i].Spaces) && i+1 == oRes.TotalResults {
								output += "        └"
							} else {
								output += "│   │   └"
							}
						} else {
							output += "│   │   ├"
						}

						output += fmt.Sprintf("─ %s\n", oRes.Resources[i].Spaces[j].Applications[k].Name)
					}
				}

				if len(oRes.Resources[i].Spaces[j].ServiceInstances) > 0 && level > 2 {
					output += bold("|   |   Services\n")
					for l := 0; l < len(oRes.Resources[i].Spaces[j].ServiceInstances); l++ {
						if l+1 == len(oRes.Resources[i].Spaces[j].ServiceInstances) {
							if j+1 == len(oRes.Resources[i].Spaces) && i+1 == oRes.TotalResults {
								output += "        └"
							} else {
								output += "│   │   └"
							}
						} else {
							output += "│   │   ├"
						}

						output += fmt.Sprintf("─ %s\n", oRes.Resources[i].Spaces[j].ServiceInstances[l].Name)
					}
				}
			}
		}
	}

	fmt.Println(output)
}
