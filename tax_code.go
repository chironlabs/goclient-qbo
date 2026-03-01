package quickbooks

import (
	"errors"
	"strconv"
)

// TaxRateDetail holds the rate reference within a TaxRateList.
type TaxRateDetail struct {
	TaxRateRef    ReferenceType `json:",omitempty"`
	TaxTypeApplicable *string  `json:",omitempty"`
	TaxOrder      *int          `json:",omitempty"`
}

// TaxRateList holds a list of tax rate details for a TaxCode.
type TaxRateList struct {
	TaxRateDetail []TaxRateDetail `json:",omitempty"`
}

// TaxCode represents a QuickBooks TaxCode object as returned by the API.
// TaxCodes are read-only — use the TaxService API to create them.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type TaxCode struct {
	ID                  string       `json:"Id,omitempty"`
	SyncToken           string       `json:",omitempty"`
	MetaData            *MetaData    `json:",omitempty"`
	Name                string       `json:",omitempty"`
	Description         *string      `json:",omitempty"`
	Active              *bool        `json:",omitempty"`
	Taxable             *bool        `json:",omitempty"`
	TaxGroup            *bool        `json:",omitempty"`
	SalesTaxRateList    *TaxRateList `json:",omitempty"`
	PurchaseTaxRateList *TaxRateList `json:",omitempty"`
}

// FindTaxCodes gets the full list of TaxCodes in the QuickBooks account.
func (c *Client) FindTaxCodes() ([]TaxCode, error) {
	var resp struct {
		QueryResponse struct {
			TaxCodes      []TaxCode `json:"TaxCode"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM TaxCode", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no tax codes could be found")
	}

	taxCodes := make([]TaxCode, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM TaxCode ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.TaxCodes == nil {
			return nil, errors.New("no tax codes could be found")
		}

		taxCodes = append(taxCodes, resp.QueryResponse.TaxCodes...)
	}

	return taxCodes, nil
}

// FindTaxCodeByID returns a tax code with a given Id.
func (c *Client) FindTaxCodeByID(id string) (*TaxCode, error) {
	var resp struct {
		TaxCode TaxCode
		Time    Date
	}

	if err := c.get("taxcode/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.TaxCode, nil
}

// QueryTaxCodes accepts an SQL query and returns all tax codes found using it.
func (c *Client) QueryTaxCodes(query string) ([]TaxCode, error) {
	var resp struct {
		QueryResponse struct {
			TaxCodes      []TaxCode `json:"TaxCode"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TaxCodes == nil {
		return nil, errors.New("could not find any tax codes")
	}

	return resp.QueryResponse.TaxCodes, nil
}
