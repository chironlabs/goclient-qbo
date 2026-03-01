package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// VendorCredit represents a QuickBooks VendorCredit object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, TotalAmt, Balance) are populated by the service.
type VendorCredit struct {
	ID                  string         `json:"Id,omitempty"`
	SyncToken           string         `json:",omitempty"`
	MetaData            *MetaData      `json:",omitempty"`
	VendorRef           ReferenceType  `json:",omitempty"`
	APAccountRef        *ReferenceType `json:",omitempty"`
	Line                []Line         `json:",omitempty"`
	TxnDate             *Date          `json:",omitempty"`
	DocNumber           *string        `json:",omitempty"`
	PrivateNote         *string        `json:",omitempty"`
	TotalAmt            json.Number    `json:",omitempty"`
	Balance             json.Number    `json:",omitempty"`
	CurrencyRef         *ReferenceType `json:",omitempty"`
	ExchangeRate        json.Number    `json:",omitempty"`
	DepartmentRef       *ReferenceType `json:",omitempty"`
	IncludeInAnnualTPAR *bool          `json:",omitempty"`
	LinkedTxn           []LinkedTxn    `json:",omitempty"`
}

// VendorCreditCreateInput contains the writable fields accepted when creating a VendorCredit.
// VendorRef and Line are required; all other fields are optional.
type VendorCreditCreateInput struct {
	VendorRef           ReferenceType  `json:",omitempty"`
	APAccountRef        *ReferenceType `json:",omitempty"`
	Line                []Line         `json:",omitempty"`
	TxnDate             *Date          `json:",omitempty"`
	DocNumber           *string        `json:",omitempty"`
	PrivateNote         *string        `json:",omitempty"`
	CurrencyRef         *ReferenceType `json:",omitempty"`
	ExchangeRate        json.Number    `json:",omitempty"`
	DepartmentRef       *ReferenceType `json:",omitempty"`
	IncludeInAnnualTPAR *bool          `json:",omitempty"`
}

// CreateVendorCredit creates the given VendorCredit on the QuickBooks server, returning
// the resulting VendorCredit object.
func (c *Client) CreateVendorCredit(input *VendorCreditCreateInput) (*VendorCredit, error) {
	var resp struct {
		VendorCredit VendorCredit
		Time         Date
	}

	if err := c.post("vendorcredit", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.VendorCredit, nil
}

// DeleteVendorCredit deletes the vendor credit.
func (c *Client) DeleteVendorCredit(vendorCredit *VendorCredit) error {
	if vendorCredit.ID == "" || vendorCredit.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("vendorcredit", vendorCredit, nil, map[string]string{"operation": "delete"})
}

// FindVendorCredits gets the full list of VendorCredits in the QuickBooks account.
func (c *Client) FindVendorCredits() ([]VendorCredit, error) {
	var resp struct {
		QueryResponse struct {
			VendorCredits []VendorCredit `json:"VendorCredit"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM VendorCredit", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no vendor credits could be found")
	}

	vendorCredits := make([]VendorCredit, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM VendorCredit ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.VendorCredits == nil {
			return nil, errors.New("no vendor credits could be found")
		}

		vendorCredits = append(vendorCredits, resp.QueryResponse.VendorCredits...)
	}

	return vendorCredits, nil
}

// FindVendorCreditByID finds the vendor credit by the given id.
func (c *Client) FindVendorCreditByID(id string) (*VendorCredit, error) {
	var resp struct {
		VendorCredit VendorCredit
		Time         Date
	}

	if err := c.get("vendorcredit/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.VendorCredit, nil
}

// QueryVendorCredits accepts an SQL query and returns all vendor credits found using it.
func (c *Client) QueryVendorCredits(query string) ([]VendorCredit, error) {
	var resp struct {
		QueryResponse struct {
			VendorCredits []VendorCredit `json:"VendorCredit"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.VendorCredits == nil {
		return nil, errors.New("could not find any vendor credits")
	}

	return resp.QueryResponse.VendorCredits, nil
}

// UpdateVendorCredit updates the vendor credit.
func (c *Client) UpdateVendorCredit(vendorCredit *VendorCredit) (*VendorCredit, error) {
	if vendorCredit.ID == "" {
		return nil, errors.New("missing vendor credit id")
	}

	existingVendorCredit, err := c.FindVendorCreditByID(vendorCredit.ID)
	if err != nil {
		return nil, err
	}

	vendorCredit.SyncToken = existingVendorCredit.SyncToken

	payload := struct {
		*VendorCredit
		Sparse bool `json:"sparse"`
	}{
		VendorCredit: vendorCredit,
		Sparse:       true,
	}

	var vendorCreditData struct {
		VendorCredit VendorCredit
		Time         Date
	}

	if err = c.post("vendorcredit", payload, &vendorCreditData, nil); err != nil {
		return nil, err
	}

	return &vendorCreditData.VendorCredit, err
}
