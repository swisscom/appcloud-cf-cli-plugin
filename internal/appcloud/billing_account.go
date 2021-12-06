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
type BillingAccountsResponse struct {
	BillingAccounts []BillingAccount `json:"resources"`
	ServerResponseError
}

// getBillingAccount retrieves a billing account by name or number.
func getBillingAccount(c plugin.CliConnection, name string) (BillingAccount, error) {
	expression := url.QueryEscape("name_number:" + name)
	accURL := fmt.Sprintf("/custom/accounts?q=%s", expression)
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", accURL)

	if err != nil {
		return BillingAccount{}, errors.Wrap(err, "Billing Account not found")
	}

	resString := strings.Join(resLines, "")
	var res BillingAccountsResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return BillingAccount{}, errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if len(res.BillingAccounts) == 0 {
		return BillingAccount{}, errors.New("Billing account not found")
	} else if len(res.BillingAccounts) > 1 {
		return BillingAccount{}, errors.New("Multiple billing accounts found - retry with the account's number or full name")
	}

	return res.BillingAccounts[0], nil
}
