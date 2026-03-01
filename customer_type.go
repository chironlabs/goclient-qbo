package quickbooks

import (
	"errors"
)

// CustomerType represents a QuickBooks CustomerType object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type CustomerType struct {
	ID        string    `json:"Id,omitempty"`
	SyncToken string    `json:",omitempty"`
	MetaData  *MetaData `json:",omitempty"`
	Name      string    `json:",omitempty"`
	Active    *bool     `json:",omitempty"`
}

// FindCustomerTypeByID returns a customerType with a given Id.
func (c *Client) FindCustomerTypeByID(id string) (*CustomerType, error) {
	var r struct {
		CustomerType CustomerType
		Time         Date
	}

	if err := c.get("customertype/"+id, &r, nil); err != nil {
		return nil, err
	}

	return &r.CustomerType, nil
}

// QueryCustomerTypes accepts an SQL query and returns all customerTypes found using it
func (c *Client) QueryCustomerTypes(query string) ([]CustomerType, error) {
	var resp struct {
		QueryResponse struct {
			CustomerTypes []CustomerType `json:"CustomerType"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.CustomerTypes == nil {
		return nil, errors.New("could not find any customerTypes")
	}

	return resp.QueryResponse.CustomerTypes, nil
}
