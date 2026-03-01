package quickbooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttachable(t *testing.T) {
	byteValue := json.RawMessage(`

{
	"Attachable": {
		"SyncToken": "0",
		"domain": "QBO",
		"AttachableRef": [
			{
				"IncludeOnSend": false,
				"EntityRef": {
					"type": "Invoice",
					"value": "95"
				}
			}
		],
		"Note": "This is an attached note.",
		"sparse": false,
		"Id": "200900000000000008541",
		"MetaData": {
			"CreateTime": "2015-11-17T11:05:15-08:00",
			"LastUpdatedTime": "2015-11-17T11:05:15-08:00"
		}
	},
	"time": "2015-11-17T11:05:15.797-08:00"
}
		`)

	var r struct {
		Attachable Attachable
		Time       Date
	}

	require.NoError(t, json.Unmarshal(byteValue, &r))

	assert.Equal(t, "0", r.Attachable.SyncToken)
	assert.False(t, *r.Attachable.AttachableRef[0].IncludeOnSend)
	assert.Equal(t, "95", r.Attachable.AttachableRef[0].EntityRef.Value)
	assert.Equal(t, "This is an attached note.", r.Attachable.Note)
	assert.Equal(t, "200900000000000008541", r.Attachable.ID)
	assert.Equal(t, "2015-11-17T11:05:15-08:00", r.Attachable.MetaData.CreateTime.String())
	assert.Equal(t, "2015-11-17T11:05:15-08:00", r.Attachable.MetaData.LastUpdatedTime.String())
}
