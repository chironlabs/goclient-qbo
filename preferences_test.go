package quickbooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreferences(t *testing.T) {
	byteValue := json.RawMessage(`
{
	"Preferences": {
		"AccountingInfoPrefs": {
			"FirstMonthOfFiscalYear": "January",
			"UseAccountNumbers": true
		},
		"CurrencyPrefs": {
			"HomeCurrency": {
				"value": "EUR"
			},
			"MultiCurrencyEnabled": false
		},
		"domain": "QBO",
		"sparse": false,
		"Id": "1",
		"SyncToken": "6",
		"MetaData": {
			"CreateTime": "2017-10-25T01:05:43-07:00",
			"LastUpdatedTime": "2018-03-08T13:24:26-08:00"
		}
	},
	"time": "2018-03-12T08:22:43.280-07:00"
}
		`)
	var r struct {
		Preferences Preferences
		Time        Date
	}

	err := json.Unmarshal(byteValue, &r)
	require.NoError(t, err)

	require.NotNil(t, r.Preferences.CurrencyPrefs)
	require.NotNil(t, r.Preferences.CurrencyPrefs.HomeCurrency)
	assert.Equal(t, "EUR", r.Preferences.CurrencyPrefs.HomeCurrency.Value)
	assert.True(t, r.Preferences.CurrencyPrefs.MultiCurrencyEnabled != nil && !*r.Preferences.CurrencyPrefs.MultiCurrencyEnabled)
	assert.Equal(t, "1", r.Preferences.ID)
	assert.Equal(t, "6", r.Preferences.SyncToken)
}
