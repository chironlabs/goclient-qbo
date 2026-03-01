package quickbooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrialBalance(t *testing.T) {
	s := json.RawMessage(`{
  "Header": {
    "ReportName": "TrialBalance",
    "Option": [
      {
        "Name": "NoReportData",
        "Value": "false"
      }
    ],
    "DateMacro": "this month-to-date",
    "ReportBasis": "Accrual",
    "StartPeriod": "2016-03-01",
    "Currency": "USD",
    "EndPeriod": "2016-03-14",
    "Time": "2016-03-14T10:11:07-07:00"
  },
  "Rows": {
    "Row": [
      {
        "ColData": [
          {
            "id": "35",
            "value": "Checking"
          },
          {
            "value": "4151.74"
          },
          {
            "value": ""
          }
        ]
      },
      {
        "ColData": [
          {
            "id": "13",
            "value": "Meals and Entertainment"
          },
          {
            "value": ""
          },
          {
            "value": "46.00"
          }
        ]
      },
      {
        "ColData": [
          {
            "id": "93",
            "value": "QuickBooks Payments Fees"
          },
          {
            "value": "0.44"
          },
          {
            "value": ""
          }
        ]
      },
      {
        "group": "GrandTotal",
        "type": "Section",
        "Summary": {
          "ColData": [
            {
              "value": "TOTAL"
            },
            {
              "value": "36587.47"
            },
            {
              "value": "36587.47"
            }
          ]
        }
      }
    ]
  },
  "Columns": {
    "Column": [
      {
        "ColType": "Account",
        "ColTitle": ""
      },
      {
        "ColType": "Money",
        "ColTitle": "Debit"
      },
      {
        "ColType": "Money",
        "ColTitle": "Credit"
      }
    ]
  }
}`)

	var tb TrialBalance
	err := json.Unmarshal(s, &tb)
	assert.NoError(t, err)

	// Header
	assert.Equal(t, "TrialBalance", tb.Header.ReportName)
	assert.Equal(t, "Accrual", tb.Header.ReportBasis)
	assert.Equal(t, "USD", tb.Header.Currency)

	// Columns
	assert.Len(t, tb.Columns, 3)
	assert.Equal(t, "Account", tb.Columns[0].ColType)
	assert.Equal(t, "Money", tb.Columns[1].ColType)
	assert.Equal(t, "Debit", tb.Columns[1].ColTitle)
	assert.Equal(t, "Credit", tb.Columns[2].ColTitle)

	// Data rows
	assert.Len(t, tb.Rows, 3)

	assert.Equal(t, "35", tb.Rows[0].RowHeader.ID)
	assert.Equal(t, "Checking", tb.Rows[0].RowHeader.Value)
	assert.Len(t, tb.Rows[0].Values, 2)
	assert.Equal(t, "4151.74", tb.Rows[0].Values[0].Value)
	assert.Equal(t, "", tb.Rows[0].Values[1].Value)

	assert.Equal(t, "13", tb.Rows[1].RowHeader.ID)
	assert.Equal(t, "Meals and Entertainment", tb.Rows[1].RowHeader.Value)
	assert.Len(t, tb.Rows[1].Values, 2)
	assert.Equal(t, "", tb.Rows[1].Values[0].Value)
	assert.Equal(t, "46.00", tb.Rows[1].Values[1].Value)

	assert.Equal(t, "93", tb.Rows[2].RowHeader.ID)
	assert.Equal(t, "QuickBooks Payments Fees", tb.Rows[2].RowHeader.Value)
	assert.Len(t, tb.Rows[2].Values, 2)
	assert.Equal(t, "0.44", tb.Rows[2].Values[0].Value)

	// Summary / total row
	assert.Equal(t, "GrandTotal", tb.Total.Group)
	assert.Equal(t, "Section", tb.Total.Type)
	assert.Len(t, tb.Total.Rows, 3)
	assert.Equal(t, "TOTAL", tb.Total.Rows[0].Value)
	assert.Equal(t, "36587.47", tb.Total.Rows[1].Value)
	assert.Equal(t, "36587.47", tb.Total.Rows[2].Value)
}
