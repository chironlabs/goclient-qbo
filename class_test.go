package quickbooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClass(t *testing.T) {
	sampleClass := json.RawMessage(`{
		"Class": {
		"FullyQualifiedName": "France", 
		"domain": "QBO", 
		"Name": "France", 
		"SyncToken": "0", 
		"SubClass": false, 
		"sparse": false, 
		"Active": true, 
		"Id": "57280", 
		"MetaData": {
		"CreateTime": "2015-07-22T13:57:27-07:00", 
		"LastUpdatedTime": "2015-07-22T13:57:27-07:00"
		}
		}, 
		"time": "2015-07-22T13:57:27.84-07:00"
		}`)

	resp := struct {
		Class Class
		Time  string `json:"time"`
	}{
		Class: Class{},
		Time:  "",
	}
	err := json.Unmarshal(sampleClass, &resp)

	class := resp.Class
	assert.NoError(t, err)
	assert.Equal(t, "57280", class.ID)
	assert.Equal(t, "France", class.Name)
	assert.True(t, *class.Active)
}
