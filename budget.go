package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// BudgetDetail holds a single budget line entry.
type BudgetDetail struct {
	BudgetDate    *Date          `json:",omitempty"`
	Amount        json.Number    `json:",omitempty"`
	AccountRef    *ReferenceType `json:",omitempty"`
	CustomerRef   *ReferenceType `json:",omitempty"`
	ClassRef      *ReferenceType `json:",omitempty"`
	DepartmentRef *ReferenceType `json:",omitempty"`
}

// Budget represents a QuickBooks Budget object as returned by the API.
// Budgets are read-only via the standard CRUD API.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type Budget struct {
	ID             string         `json:"Id,omitempty"`
	SyncToken      string         `json:",omitempty"`
	MetaData       *MetaData      `json:",omitempty"`
	Name           string         `json:",omitempty"`
	FiscalYear     *string        `json:",omitempty"`
	StartDate      *Date          `json:",omitempty"`
	EndDate        *Date          `json:",omitempty"`
	BudgetType     *string        `json:",omitempty"`
	BudgetEntryType *string       `json:",omitempty"`
	Active         *bool          `json:",omitempty"`
	BudgetDetail   []BudgetDetail `json:",omitempty"`
}

// FindBudgets gets the full list of Budgets in the QuickBooks account.
func (c *Client) FindBudgets() ([]Budget, error) {
	var resp struct {
		QueryResponse struct {
			Budgets       []Budget `json:"Budget"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Budget", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no budgets could be found")
	}

	budgets := make([]Budget, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Budget ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Budgets == nil {
			return nil, errors.New("no budgets could be found")
		}

		budgets = append(budgets, resp.QueryResponse.Budgets...)
	}

	return budgets, nil
}

// QueryBudgets accepts an SQL query and returns all budgets found using it.
func (c *Client) QueryBudgets(query string) ([]Budget, error) {
	var resp struct {
		QueryResponse struct {
			Budgets       []Budget `json:"Budget"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Budgets == nil {
		return nil, errors.New("could not find any budgets")
	}

	return resp.QueryResponse.Budgets, nil
}
