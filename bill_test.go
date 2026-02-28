package quickbooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBill(t *testing.T) {
	byteValue := json.RawMessage(`
{
  "Bill": {
    "SyncToken": "2",
    "domain": "QBO",
    "APAccountRef": {
      "name": "Accounts Payable (A/P)",
      "value": "33"
    },
    "VendorRef": {
      "name": "Norton Lumber and Building Materials",
      "value": "46"
    },
    "TxnDate": "2014-11-06",
    "TotalAmt": 103.55,
    "CurrencyRef": {
      "name": "United States Dollar",
      "value": "USD"
    },
    "LinkedTxn": [
      {
        "TxnId": "118",
        "TxnType": "BillPaymentCheck"
      }
    ],
    "SalesTermRef": {
      "value": "3"
    },
    "DueDate": "2014-12-06",
    "sparse": false,
    "Line": [
      {
        "DetailType": "AccountBasedExpenseLineDetail",
        "Amount": 103.55,
        "Id": "1",
        "AccountBasedExpenseLineDetail": {
          "TaxCodeRef": {
            "value": "TAX"
          },
          "AccountRef": {
            "name": "Job Expenses:Job Materials:Decks and Patios",
            "value": "64"
          },
          "BillableStatus": "Billable",
          "CustomerRef": {
            "name": "Travis Waldron",
            "value": "26"
          }
        },
        "Description": "Lumber"
      }
    ],
    "Balance": 0,
    "Id": "25",
    "MetaData": {
      "CreateTime": "2014-11-06T15:37:25-08:00",
      "LastUpdatedTime": "2015-02-09T10:11:11-08:00"
    }
  },
  "time": "2015-02-09T10:17:20.251-08:00"
}
		`)

	var r struct {
		Bill Bill
		Time Date
	}
	err := json.Unmarshal(byteValue, &r)
	assert.NoError(t, err)
	assert.Equal(t, "2", r.Bill.SyncToken)
	assert.Equal(t, "Accounts Payable (A/P)", r.Bill.APAccountRef.Name)
	assert.Equal(t, "33", r.Bill.APAccountRef.Value)
	assert.Equal(t, "Norton Lumber and Building Materials", r.Bill.VendorRef.Name)
	assert.Equal(t, "46", r.Bill.VendorRef.Value)
	assert.Equal(t, "2014-11-06T00:00:00+00:00", r.Bill.TxnDate.String())
	totalAmt, _ := r.Bill.TotalAmt.Float64()
	assert.Equal(t, 103.55, totalAmt)
	assert.Equal(t, "United States Dollar", r.Bill.CurrencyRef.Name)
	assert.Equal(t, "USD", r.Bill.CurrencyRef.Value)
	// LinkedTxn
	assert.Equal(t, "3", r.Bill.SalesTermRef.Value)
	assert.Equal(t, "2014-12-06T00:00:00+00:00", r.Bill.DueDate.String())
	assert.Equal(t, 1, len(r.Bill.Line))
	balance, _ := r.Bill.Balance.Int64()
	assert.Equal(t, int64(0), balance)
	assert.Equal(t, "25", r.Bill.ID)
	assert.Equal(t, "2014-11-06T15:37:25-08:00", r.Bill.MetaData.CreateTime.String())
	assert.Equal(t, "2015-02-09T10:11:11-08:00", r.Bill.MetaData.LastUpdatedTime.String())
}
