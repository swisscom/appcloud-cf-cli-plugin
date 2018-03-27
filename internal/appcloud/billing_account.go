package appcloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// BillingAccount is an entity of the Swisscom Application Cloud which handles billing.
type BillingAccount struct {
	Metadata struct {
		GUID string `json:"guid"`
	} `json:"metadata"`
}

// BillingAccountResponse is a response from the server to a billing account request.
type BillingAccountResponse struct {
	BillingAccount
	ServerResponseError
}

// getBillingAccount retrieves a billing account by name.
func getBillingAccount(c plugin.CliConnection, name string) (BillingAccount, error) {
	accURL := url.QueryEscape(fmt.Sprintf("/custom/accounts?q=name:%s", name))
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", accURL)

	if err != nil {
		return BillingAccount{}, errors.Wrap(err, "Billing Account not found")
	}

	resString := strings.Join(resLines, "")
	var res BillingAccountResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return BillingAccount{}, errors.Wrap(err, "Couldn't read JSON response from server")
	}

	return res.BillingAccount, nil
}
