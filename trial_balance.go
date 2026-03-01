package quickbooks

import (
	"encoding/json"
	"time"
)

type TrialBalanceHeader struct {
	ReportName  string
	Option      []NameValue
	DateMacro   string
	ReportBasis string
	StartPeriod string
	Currency    string
	EndPeriod   string
	Time        time.Time
}

type TrialBalanceColumn struct {
	ColType  string
	ColTitle string
}

type TrialBalanceValueRow struct {
	Value string `json:"value"`
}

type TrialBalanceRowHeader struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type TrialBalanceRow struct {
	RowHeader TrialBalanceRowHeader
	Values    []TrialBalanceValueRow
}

type TrialBalanceSummaryRow struct {
	Group string `json:"group"`
	Type  string `json:"type"`
	Rows  []TrialBalanceValueRow
}

type TrialBalance struct {
	Header  TrialBalanceHeader
	Columns []TrialBalanceColumn
	Rows    []TrialBalanceRow
	Total   TrialBalanceSummaryRow
}

func (tb *TrialBalance) UnmarshalJSON(data []byte) error {
	var t struct {
		Header TrialBalanceHeader
		Rows   struct {
			Row []json.RawMessage
		}
		Columns struct {
			Column []TrialBalanceColumn
		}
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	tb.Header = t.Header
	tb.Columns = t.Columns.Column

	for _, raw := range t.Rows.Row {
		var r map[string]json.RawMessage
		if err := json.Unmarshal(raw, &r); err != nil {
			return err
		}

		cdr, ok := r["ColData"]
		if ok {
			// Data row: the entire ColData array is one TrialBalanceRow.
			var colData []json.RawMessage
			if err := json.Unmarshal(cdr, &colData); err != nil {
				return err
			}

			var actualRow TrialBalanceRow
			for _, elem := range colData {
				var v map[string]json.RawMessage
				if err := json.Unmarshal(elem, &v); err != nil {
					return err
				}

				if _, hasID := v["id"]; hasID {
					// Element with "id" is the account name/id header.
					if err := json.Unmarshal(elem, &actualRow.RowHeader); err != nil {
						return err
					}
				} else {
					var val TrialBalanceValueRow
					if err := json.Unmarshal(elem, &val); err != nil {
						return err
					}
					actualRow.Values = append(actualRow.Values, val)
				}
			}
			tb.Rows = append(tb.Rows, actualRow)
		} else {
			// Summary row: unmarshal the full raw message, not the missing cdr.
			var sr struct {
				Group   string `json:"group"`
				Type    string `json:"type"`
				Summary struct {
					ColData []TrialBalanceValueRow
				}
			}

			if err := json.Unmarshal(raw, &sr); err != nil {
				return err
			}

			tb.Total = TrialBalanceSummaryRow{
				Group: sr.Group,
				Type:  sr.Type,
				Rows:  sr.Summary.ColData,
			}
		}
	}

	return nil
}

// TrialBalanceQueryParams holds the optional query parameters for the TrialBalance report.
type TrialBalanceQueryParams struct {
	// Cash or Accrual
	AccountingMethod *string
	StartDate        *string
	EndDate          *string
	// Predefined date range, e.g. "This Month", "Last Fiscal Year". Ignored when StartDate/EndDate are set.
	DateMacro *string
	// ascend or descend
	SortOrder *string
	// Total, Month, Week, Days, Quarter, Year, Customers, Vendors, Classes, Departments, Employees, ProductsAndServices
	SummarizeColumnBy *string
}

func (p *TrialBalanceQueryParams) toMap() map[string]string {
	m := map[string]string{}
	if p.AccountingMethod != nil {
		m["accounting_method"] = *p.AccountingMethod
	}
	if p.StartDate != nil {
		m["start_date"] = *p.StartDate
	}
	if p.EndDate != nil {
		m["end_date"] = *p.EndDate
	}
	if p.DateMacro != nil {
		m["date_macro"] = *p.DateMacro
	}
	if p.SortOrder != nil {
		m["sort_order"] = *p.SortOrder
	}
	if p.SummarizeColumnBy != nil {
		m["summarize_column_by"] = *p.SummarizeColumnBy
	}
	return m
}

// GetTrialBalance fetches a TrialBalance report from the QBO API.
// Pass nil for params to use the API defaults.
func (c *Client) GetTrialBalance(params *TrialBalanceQueryParams) (*TrialBalance, error) {
	var queryParams map[string]string
	if params != nil {
		queryParams = params.toMap()
	}
	var tb TrialBalance
	if err := c.get("reports/TrialBalance", &tb, queryParams); err != nil {
		return nil, err
	}
	return &tb, nil
}
