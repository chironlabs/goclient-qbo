package quickbooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVendor(t *testing.T) {
	byteValue := json.RawMessage(`
{
  "Vendor": {
    "PrimaryEmailAddr": {
      "Address": "Books@Intuit.com"
    },
    "Vendor1099": false,
    "domain": "QBO",
    "GivenName": "Bessie",
    "DisplayName": "Books by Bessie",
    "BillAddr": {
      "City": "Palo Alto",
      "Line1": "15 Main St.",
      "PostalCode": "94303",
      "Lat": "37.445013",
      "Long": "-122.1391443",
      "CountrySubDivisionCode": "CA",
      "Id": "31"
    },
    "SyncToken": "0",
    "PrintOnCheckName": "Books by Bessie",
    "FamilyName": "Williams",
    "PrimaryPhone": {
      "FreeFormNumber": "(650) 555-7745"
    },
    "AcctNum": "1345",
    "CompanyName": "Books by Bessie",
    "WebAddr": {
      "URI": "http://www.booksbybessie.co"
    },
    "sparse": false,
    "Active": true,
    "Balance": 0,
    "Id": "30",
    "MetaData": {
      "CreateTime": "2014-09-12T10:07:56-07:00",
      "LastUpdatedTime": "2014-09-17T11:13:46-07:00"
    }
  },
  "time": "2015-07-28T13:33:09.453-07:00"
}
		`)

	var resp struct {
		Vendor Vendor
		Time   Date
	}

	require.NoError(t, json.Unmarshal(byteValue, &resp))
	assert.NotNil(t, resp.Vendor.PrimaryEmailAddr)
	assert.False(t, resp.Vendor.Vendor1099)
	assert.Equal(t, "Bessie", resp.Vendor.GivenName)
	assert.Equal(t, "Books by Bessie", resp.Vendor.DisplayName)
	assert.NotNil(t, resp.Vendor.BillAddr)
	assert.Equal(t, "0", resp.Vendor.SyncToken)
	assert.Equal(t, "Books by Bessie", resp.Vendor.PrintOnCheckName)
	assert.Equal(t, "Williams", resp.Vendor.FamilyName)
	assert.NotNil(t, resp.Vendor.PrimaryPhone)
	assert.Equal(t, "1345", resp.Vendor.AcctNum)
	assert.Equal(t, "Books by Bessie", resp.Vendor.CompanyName)
	assert.NotNil(t, resp.Vendor.WebAddr)
	assert.True(t, resp.Vendor.Active)
	assert.Equal(t, "0", resp.Vendor.Balance.String())
	assert.Equal(t, "30", resp.Vendor.ID)
	assert.Equal(t, "2014-09-12T10:07:56-07:00", resp.Vendor.MetaData.CreateTime.String())
	assert.Equal(t, "2014-09-17T11:13:46-07:00", resp.Vendor.MetaData.LastUpdatedTime.String())
}
