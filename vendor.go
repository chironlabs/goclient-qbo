package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Vendor represents a QuickBooks Vendor object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, Balance) are populated by the service.
type Vendor struct {
	ID               string           `json:"Id,omitempty"`
	SyncToken        string           `json:",omitempty"`
	MetaData         *MetaData        `json:",omitempty"`
	Title            string           `json:",omitempty"`
	GivenName        string           `json:",omitempty"`
	MiddleName       *string          `json:",omitempty"`
	Suffix           *string          `json:",omitempty"`
	FamilyName       string           `json:",omitempty"`
	PrimaryEmailAddr *EmailAddress    `json:",omitempty"`
	DisplayName      string           `json:",omitempty"`
	APAccountRef     *ReferenceType   `json:",omitempty"`
	TermRef          *ReferenceType   `json:",omitempty"`
	GSTIN            *string          `json:",omitempty"`
	Fax              *TelephoneNumber `json:",omitempty"`
	BusinessNumber   *string          `json:",omitempty"`
	CurrencyRef      *ReferenceType   `json:",omitempty"`
	HasTPAR          *bool            `json:",omitempty"`
	TaxReportingBasis *string         `json:",omitempty"`
	Mobile           *TelephoneNumber `json:",omitempty"`
	PrimaryPhone     *TelephoneNumber `json:",omitempty"`
	Active           *bool            `json:",omitempty"`
	AlternatePhone   *TelephoneNumber `json:",omitempty"`
	Vendor1099       *bool            `json:",omitempty"`
	BillRate         json.Number      `json:",omitempty"`
	WebAddr          *WebSiteAddress  `json:",omitempty"`
	CompanyName      string           `json:",omitempty"`
	TaxIdentifier    *string          `json:",omitempty"`
	AcctNum          *string          `json:",omitempty"`
	GSTRegistrationType *string       `json:",omitempty"`
	PrintOnCheckName *string          `json:",omitempty"`
	BillAddr         *Address         `json:",omitempty"`
	Balance          json.Number      `json:",omitempty"`
}

// VendorCreateInput contains the writable fields accepted when creating a Vendor.
// At least one of GivenName, FamilyName, DisplayName, or CompanyName is required.
type VendorCreateInput struct {
	GivenName           string           `json:",omitempty"`
	FamilyName          string           `json:",omitempty"`
	DisplayName         string           `json:",omitempty"`
	CompanyName         string           `json:",omitempty"`
	Title               string           `json:",omitempty"`
	MiddleName          *string          `json:",omitempty"`
	Suffix              *string          `json:",omitempty"`
	PrintOnCheckName    *string          `json:",omitempty"`
	Active              *bool            `json:",omitempty"`
	PrimaryEmailAddr    *EmailAddress    `json:",omitempty"`
	PrimaryPhone        *TelephoneNumber `json:",omitempty"`
	AlternatePhone      *TelephoneNumber `json:",omitempty"`
	Mobile              *TelephoneNumber `json:",omitempty"`
	Fax                 *TelephoneNumber `json:",omitempty"`
	WebAddr             *WebSiteAddress  `json:",omitempty"`
	BillAddr            *Address         `json:",omitempty"`
	APAccountRef        *ReferenceType   `json:",omitempty"`
	TermRef             *ReferenceType   `json:",omitempty"`
	CurrencyRef         *ReferenceType   `json:",omitempty"`
	AcctNum             *string          `json:",omitempty"`
	TaxIdentifier       *string          `json:",omitempty"`
	BusinessNumber      *string          `json:",omitempty"`
	GSTIN               *string          `json:",omitempty"`
	GSTRegistrationType *string          `json:",omitempty"`
	TaxReportingBasis   *string          `json:",omitempty"`
	HasTPAR             *bool            `json:",omitempty"`
	Vendor1099          *bool            `json:",omitempty"`
	BillRate            json.Number      `json:",omitempty"`
}

// CreateVendor creates the given Vendor on the QuickBooks server, returning
// the resulting Vendor object.
func (c *Client) CreateVendor(input *VendorCreateInput) (*Vendor, error) {
	var resp struct {
		Vendor Vendor
		Time   Date
	}

	if err := c.post("vendor", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Vendor, nil
}

// FindVendors gets the full list of Vendors in the QuickBooks account.
func (c *Client) FindVendors() ([]Vendor, error) {
	var resp struct {
		QueryResponse struct {
			Vendors       []Vendor `json:"Vendor"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Vendor", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no vendors could be found")
	}

	vendors := make([]Vendor, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Vendor ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Vendors == nil {
			return nil, errors.New("no vendors could be found")
		}

		vendors = append(vendors, resp.QueryResponse.Vendors...)
	}

	return vendors, nil
}

// FindVendorByID finds the vendor by the given id
func (c *Client) FindVendorByID(id string) (*Vendor, error) {
	var resp struct {
		Vendor Vendor
		Time   Date
	}

	if err := c.get("vendor/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Vendor, nil
}

// QueryVendors accepts an SQL query and returns all vendors found using it
func (c *Client) QueryVendors(query string) ([]Vendor, error) {
	var resp struct {
		QueryResponse struct {
			Vendors       []Vendor `json:"Vendor"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Vendors == nil {
		return nil, errors.New("could not find any vendors")
	}

	return resp.QueryResponse.Vendors, nil
}

// DeleteVendor deletes the vendor.
func (c *Client) DeleteVendor(vendor *Vendor) error {
	if vendor.ID == "" || vendor.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("vendor", vendor, nil, map[string]string{"operation": "delete"})
}

// UpdateVendor updates the vendor
func (c *Client) UpdateVendor(vendor *Vendor) (*Vendor, error) {
	if vendor.ID == "" {
		return nil, errors.New("missing vendor id")
	}

	existingVendor, err := c.FindVendorByID(vendor.ID)
	if err != nil {
		return nil, err
	}

	vendor.SyncToken = existingVendor.SyncToken

	payload := struct {
		*Vendor
		Sparse bool `json:"sparse"`
	}{
		Vendor: vendor,
		Sparse: true,
	}

	var vendorData struct {
		Vendor Vendor
		Time   Date
	}

	if err = c.post("vendor", payload, &vendorData, nil); err != nil {
		return nil, err
	}

	return &vendorData.Vendor, err
}
