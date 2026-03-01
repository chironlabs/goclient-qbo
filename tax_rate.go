package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// EffectiveTaxRate holds a time-bounded tax rate value.
type EffectiveTaxRate struct {
	RateValue   json.Number `json:",omitempty"`
	EffectiveDate *Date     `json:",omitempty"`
	EndDate       *Date     `json:",omitempty"`
}

// TaxRate represents a QuickBooks TaxRate object as returned by the API.
// TaxRates are read-only — use the TaxService API to create them.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type TaxRate struct {
	ID               string             `json:"Id,omitempty"`
	SyncToken        string             `json:",omitempty"`
	MetaData         *MetaData          `json:",omitempty"`
	Name             string             `json:",omitempty"`
	Description      *string            `json:",omitempty"`
	RateValue        json.Number        `json:",omitempty"`
	AgencyRef        *ReferenceType     `json:",omitempty"`
	TaxCode          *string            `json:",omitempty"`
	DisplayType      *string            `json:",omitempty"`
	Active           *bool              `json:",omitempty"`
	EffectiveTaxRate []EffectiveTaxRate `json:",omitempty"`
}

// FindTaxRates gets the full list of TaxRates in the QuickBooks account.
func (c *Client) FindTaxRates() ([]TaxRate, error) {
	var resp struct {
		QueryResponse struct {
			TaxRates      []TaxRate `json:"TaxRate"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM TaxRate", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no tax rates could be found")
	}

	taxRates := make([]TaxRate, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM TaxRate ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.TaxRates == nil {
			return nil, errors.New("no tax rates could be found")
		}

		taxRates = append(taxRates, resp.QueryResponse.TaxRates...)
	}

	return taxRates, nil
}

// FindTaxRateByID returns a tax rate with a given Id.
func (c *Client) FindTaxRateByID(id string) (*TaxRate, error) {
	var resp struct {
		TaxRate TaxRate
		Time    Date
	}

	if err := c.get("taxrate/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.TaxRate, nil
}

// QueryTaxRates accepts an SQL query and returns all tax rates found using it.
func (c *Client) QueryTaxRates(query string) ([]TaxRate, error) {
	var resp struct {
		QueryResponse struct {
			TaxRates      []TaxRate `json:"TaxRate"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TaxRates == nil {
		return nil, errors.New("could not find any tax rates")
	}

	return resp.QueryResponse.TaxRates, nil
}
