package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// ServiceEventResponse is the response from the server from a service event call
type ServiceEventResponse struct {
	ServiceEvent
	ServerResponse
}

type ServiceInstancesResponse struct {
	SpaceSummary
	ServerResponse
}
type ServiceInstancesGUIDResponse struct {
	si []string
	ServerResponse
}

// ServiceEvents Returns events for given service instance
func (p *AppCloudPlugin) ServiceEvents(c plugin.CliConnection, serviceInstance string) error {
	s, err := c.GetCurrentSpace()
	o, orgErr :=c.GetCurrentOrg()
	if orgErr != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}
	uName, nameErr := c.Username()

	if nameErr != nil {
		return fmt.Errorf("Couldn't retrieve username")
	}

	fmt.Printf("Getting events for service %s in org %s / space %s as %s ...\n", cyanBold(serviceInstance),
	cyanBold(o.Name), cyanBold(s.Name), cyanBold(uName))
	fmt.Print(greenBold("OK\n\n"))
	// Get the current targeted space details
	if err != nil {
		return fmt.Errorf("Couldn't retrieve space")
	}
	url := fmt.Sprintf("/custom/spaces/%s/summary",s.Guid)
	servicesList, err := c.CliCommandWithoutTerminalOutput("curl", url)
	
	if err != nil {
		return fmt.Errorf("Couldn't retrieve service instances for space %s",s.Name)
	}

	servicesListString := strings.Join(servicesList, "")
	var bRes ServiceInstancesResponse
	err = json.Unmarshal([]byte(servicesListString), &bRes)
	if err != nil {
		return  fmt.Errorf("Couldn't read JSON response: %s", err.Error())
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}
	serviceGUID := ""
	for i := 0; i < len(bRes.Services); i++ {
		if(bRes.Services[i].Name == serviceInstance){
			serviceGUID = bRes.Services[i].GUID
			break
		}
	}

	url = fmt.Sprintf("/v2/events?q=actee:%s",serviceGUID)
	serviceEvents, err := c.CliCommandWithoutTerminalOutput("curl", url)
	
	serviceEventString := strings.Join(serviceEvents, "")
	var bEventRes ServiceEventResponse
	err = json.Unmarshal([]byte(serviceEventString), &bEventRes)
	if err != nil {
		return fmt.Errorf("Couldn't read JSON response: %s", err.Error())
	}
	if(len(bEventRes.ServiceEvent.Resources) >0){
		fmt.Println(bold("time")+"                   "+bold("event")+"                          "+bold("actor")+"                   "+bold("description"))
	}
	for i := 0; i < len(bEventRes.ServiceEvent.Resources); i++ {
		fmt.Println(cyanBold(bEventRes.ServiceEvent.Resources[i].Entity.TimeStamp)+
		"   "+bEventRes.ServiceEvent.Resources[i].Entity.Type+
		"  "+bEventRes.ServiceEvent.Resources[i].Entity.ActorUserName+
		""+"")
	}

	if bRes.ErrorCode != "" {
		return errors.New(bRes.Description)
	}
	return nil
}
