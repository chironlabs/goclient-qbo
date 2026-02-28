package quickbooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	byteValue := json.RawMessage(`
{	
	"Account": {
		"FullyQualifiedName": "MyJobs",
		"domain": "QBO",
		"Name": "MyJobs",
		"Classification": "Asset",
		"AccountSubType": "AccountsReceivable",
		"CurrencyRef": {
			"name": "United States Dollar",
			"value": "USD"
		},
		"CurrentBalanceWithSubAccounts": 0,
		"sparse": false,
		"MetaData": {
			"CreateTime": "2014-12-31T09:29:05-08:00",
			"LastUpdatedTime": "2014-12-31T09:29:05-08:00"
		},
		"AccountType": "Accounts Receivable",
		"CurrentBalance": 0,
		"Active": true,
		"SyncToken": "0",
		"Id": "94",
		"SubAccount": false
	},
	"time": "2014-12-31T09:29:05.717-08:00"
}
		`)

	var r struct {
		Account Account
		Time    Date
	}

	err := json.Unmarshal(byteValue, &r)
	require.NoError(t, err)

	assert.Equal(t, "MyJobs", r.Account.FullyQualifiedName)
	assert.Equal(t, "MyJobs", r.Account.Name)
	assert.Equal(t, "Asset", r.Account.Classification)
	assert.Equal(t, "AccountsReceivable", r.Account.AccountSubType)
	assert.Equal(t, json.Number("0"), r.Account.CurrentBalanceWithSubAccounts)
	assert.Equal(t, "2014-12-31T09:29:05-08:00", r.Account.MetaData.CreateTime.String())
	assert.Equal(t, "2014-12-31T09:29:05-08:00", r.Account.MetaData.LastUpdatedTime.String())
	assert.Equal(t, AccountsReceivableAccountType, r.Account.AccountType)
	assert.Equal(t, json.Number("0"), r.Account.CurrentBalance)
	assert.True(t, r.Account.Active != nil && *r.Account.Active == true)
	assert.Equal(t, "0", r.Account.SyncToken)
	assert.Equal(t, "94", r.Account.ID)
	assert.False(t, r.Account.SubAccount)
}
