package quickbooks

import (
	"errors"
	"strconv"
)

// TaxAgency represents a QuickBooks TaxAgency object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type TaxAgency struct {
	ID                      string         `json:"Id,omitempty"`
	SyncToken               string         `json:",omitempty"`
	MetaData                *MetaData      `json:",omitempty"`
	DisplayName             string         `json:",omitempty"`
	TaxRegistrationNumber   *string        `json:",omitempty"`
	TaxTrackedOnPurchases   *bool          `json:",omitempty"`
	TaxOnPurchasesAccountRef *ReferenceType `json:",omitempty"`
	TaxTrackedOnSales       *bool          `json:",omitempty"`
	TaxOnSalesAccountRef    *ReferenceType `json:",omitempty"`
	LastFileDate            *Date          `json:",omitempty"`
}

// TaxAgencyCreateInput contains the writable fields accepted when creating a TaxAgency.
// DisplayName is required; all other fields are optional.
type TaxAgencyCreateInput struct {
	DisplayName             string         `json:",omitempty"`
	TaxRegistrationNumber   *string        `json:",omitempty"`
	TaxTrackedOnPurchases   *bool          `json:",omitempty"`
	TaxOnPurchasesAccountRef *ReferenceType `json:",omitempty"`
	TaxTrackedOnSales       *bool          `json:",omitempty"`
	TaxOnSalesAccountRef    *ReferenceType `json:",omitempty"`
}

// CreateTaxAgency creates the given TaxAgency on the QuickBooks server, returning
// the resulting TaxAgency object.
func (c *Client) CreateTaxAgency(input *TaxAgencyCreateInput) (*TaxAgency, error) {
	var resp struct {
		TaxAgency TaxAgency
		Time      Date
	}

	if err := c.post("taxagency", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.TaxAgency, nil
}

// FindTaxAgencies gets the full list of TaxAgencies in the QuickBooks account.
func (c *Client) FindTaxAgencies() ([]TaxAgency, error) {
	var resp struct {
		QueryResponse struct {
			TaxAgencies   []TaxAgency `json:"TaxAgency"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM TaxAgency", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no tax agencies could be found")
	}

	taxAgencies := make([]TaxAgency, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM TaxAgency ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.TaxAgencies == nil {
			return nil, errors.New("no tax agencies could be found")
		}

		taxAgencies = append(taxAgencies, resp.QueryResponse.TaxAgencies...)
	}

	return taxAgencies, nil
}

// FindTaxAgencyByID returns a tax agency with a given Id.
func (c *Client) FindTaxAgencyByID(id string) (*TaxAgency, error) {
	var resp struct {
		TaxAgency TaxAgency
		Time      Date
	}

	if err := c.get("taxagency/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.TaxAgency, nil
}

// QueryTaxAgencies accepts an SQL query and returns all tax agencies found using it.
func (c *Client) QueryTaxAgencies(query string) ([]TaxAgency, error) {
	var resp struct {
		QueryResponse struct {
			TaxAgencies   []TaxAgency `json:"TaxAgency"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TaxAgencies == nil {
		return nil, errors.New("could not find any tax agencies")
	}

	return resp.QueryResponse.TaxAgencies, nil
}
