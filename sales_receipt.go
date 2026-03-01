package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// SalesReceipt represents a QuickBooks SalesReceipt object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, TotalAmt, Balance, TxnSource) are populated by the service.
type SalesReceipt struct {
	ID                           string         `json:"Id,omitempty"`
	SyncToken                    string         `json:",omitempty"`
	MetaData                     *MetaData      `json:",omitempty"`
	CustomerRef                  *ReferenceType `json:",omitempty"`
	CustomerMemo                 *MemoRef       `json:",omitempty"`
	DocNumber                    *string        `json:",omitempty"`
	TxnDate                      *Date          `json:",omitempty"`
	DepartmentRef                *ReferenceType `json:",omitempty"`
	PrivateNote                  *string        `json:",omitempty"`
	Line                         []Line         `json:",omitempty"`
	TxnTaxDetail                 *TxnTaxDetail  `json:",omitempty"`
	BillAddr                     *Address       `json:",omitempty"`
	ShipAddr                     *Address       `json:",omitempty"`
	ClassRef                     *ReferenceType `json:",omitempty"`
	ShipMethodRef                *ReferenceType `json:",omitempty"`
	ShipDate                     *Date          `json:",omitempty"`
	TrackingNum                  *string        `json:",omitempty"`
	TotalAmt                     json.Number    `json:",omitempty"`
	CurrencyRef                  *ReferenceType `json:",omitempty"`
	ExchangeRate                 json.Number    `json:",omitempty"`
	DepositToAccountRef          *ReferenceType `json:",omitempty"`
	ApplyTaxAfterDiscount        *bool          `json:",omitempty"`
	PrintStatus                  *string        `json:",omitempty"`
	EmailStatus                  *string        `json:",omitempty"`
	BillEmail                    *EmailAddress  `json:",omitempty"`
	BillEmailCC                  *EmailAddress  `json:"BillEmailCc,omitempty"`
	BillEmailBCC                 *EmailAddress  `json:"BillEmailBcc,omitempty"`
	DeliveryInfo                 *DeliveryInfo  `json:",omitempty"`
	Balance                      json.Number    `json:",omitempty"`
	TxnSource                    *string        `json:",omitempty"`
	PaymentMethodRef             *ReferenceType `json:",omitempty"`
	CustomField                  []CustomField  `json:",omitempty"`
}

// SalesReceiptCreateInput contains the writable fields accepted when creating a SalesReceipt.
// Line is required; all other fields are optional.
type SalesReceiptCreateInput struct {
	Line                         []Line         `json:",omitempty"`
	CustomerRef                  *ReferenceType `json:",omitempty"`
	CustomerMemo                 *MemoRef       `json:",omitempty"`
	DocNumber                    *string        `json:",omitempty"`
	TxnDate                      *Date          `json:",omitempty"`
	DepartmentRef                *ReferenceType `json:",omitempty"`
	PrivateNote                  *string        `json:",omitempty"`
	TxnTaxDetail                 *TxnTaxDetail  `json:",omitempty"`
	BillAddr                     *Address       `json:",omitempty"`
	ShipAddr                     *Address       `json:",omitempty"`
	ClassRef                     *ReferenceType `json:",omitempty"`
	ShipMethodRef                *ReferenceType `json:",omitempty"`
	ShipDate                     *Date          `json:",omitempty"`
	TrackingNum                  *string        `json:",omitempty"`
	CurrencyRef                  *ReferenceType `json:",omitempty"`
	ExchangeRate                 json.Number    `json:",omitempty"`
	DepositToAccountRef          *ReferenceType `json:",omitempty"`
	ApplyTaxAfterDiscount        *bool          `json:",omitempty"`
	PrintStatus                  *string        `json:",omitempty"`
	EmailStatus                  *string        `json:",omitempty"`
	BillEmail                    *EmailAddress  `json:",omitempty"`
	BillEmailCC                  *EmailAddress  `json:"BillEmailCc,omitempty"`
	BillEmailBCC                 *EmailAddress  `json:"BillEmailBcc,omitempty"`
	PaymentMethodRef             *ReferenceType `json:",omitempty"`
	CustomField                  []CustomField  `json:",omitempty"`
}

// CreateSalesReceipt creates the given SalesReceipt on the QuickBooks server, returning
// the resulting SalesReceipt object.
func (c *Client) CreateSalesReceipt(input *SalesReceiptCreateInput) (*SalesReceipt, error) {
	var resp struct {
		SalesReceipt SalesReceipt
		Time         Date
	}

	if err := c.post("salesreceipt", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.SalesReceipt, nil
}

// DeleteSalesReceipt deletes the sales receipt.
func (c *Client) DeleteSalesReceipt(salesReceipt *SalesReceipt) error {
	if salesReceipt.ID == "" || salesReceipt.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("salesreceipt", salesReceipt, nil, map[string]string{"operation": "delete"})
}

// FindSalesReceipts gets the full list of SalesReceipts in the QuickBooks account.
func (c *Client) FindSalesReceipts() ([]SalesReceipt, error) {
	var resp struct {
		QueryResponse struct {
			SalesReceipts []SalesReceipt `json:"SalesReceipt"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM SalesReceipt", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no sales receipts could be found")
	}

	salesReceipts := make([]SalesReceipt, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM SalesReceipt ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.SalesReceipts == nil {
			return nil, errors.New("no sales receipts could be found")
		}

		salesReceipts = append(salesReceipts, resp.QueryResponse.SalesReceipts...)
	}

	return salesReceipts, nil
}

// FindSalesReceiptByID finds the sales receipt by the given id.
func (c *Client) FindSalesReceiptByID(id string) (*SalesReceipt, error) {
	var resp struct {
		SalesReceipt SalesReceipt
		Time         Date
	}

	if err := c.get("salesreceipt/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.SalesReceipt, nil
}

// QuerySalesReceipts accepts an SQL query and returns all sales receipts found using it.
func (c *Client) QuerySalesReceipts(query string) ([]SalesReceipt, error) {
	var resp struct {
		QueryResponse struct {
			SalesReceipts []SalesReceipt `json:"SalesReceipt"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.SalesReceipts == nil {
		return nil, errors.New("could not find any sales receipts")
	}

	return resp.QueryResponse.SalesReceipts, nil
}

// SendSalesReceipt sends the sales receipt to the SalesReceipt.BillEmail if emailAddress is left empty.
func (c *Client) SendSalesReceipt(salesReceiptId string, emailAddress string) error {
	queryParameters := make(map[string]string)

	if emailAddress != "" {
		queryParameters["sendTo"] = emailAddress
	}

	return c.post("salesreceipt/"+salesReceiptId+"/send", nil, nil, queryParameters)
}

// UpdateSalesReceipt updates the sales receipt.
func (c *Client) UpdateSalesReceipt(salesReceipt *SalesReceipt) (*SalesReceipt, error) {
	if salesReceipt.ID == "" {
		return nil, errors.New("missing sales receipt id")
	}

	existingSalesReceipt, err := c.FindSalesReceiptByID(salesReceipt.ID)
	if err != nil {
		return nil, err
	}

	salesReceipt.SyncToken = existingSalesReceipt.SyncToken

	payload := struct {
		*SalesReceipt
		Sparse bool `json:"sparse"`
	}{
		SalesReceipt: salesReceipt,
		Sparse:       true,
	}

	var salesReceiptData struct {
		SalesReceipt SalesReceipt
		Time         Date
	}

	if err = c.post("salesreceipt", payload, &salesReceiptData, nil); err != nil {
		return nil, err
	}

	return &salesReceiptData.SalesReceipt, err
}

// VoidSalesReceipt voids the sales receipt.
func (c *Client) VoidSalesReceipt(salesReceipt *SalesReceipt) error {
	if salesReceipt.ID == "" {
		return errors.New("missing sales receipt id")
	}

	existingSalesReceipt, err := c.FindSalesReceiptByID(salesReceipt.ID)
	if err != nil {
		return err
	}

	salesReceipt.SyncToken = existingSalesReceipt.SyncToken

	return c.post("salesreceipt", salesReceipt, nil, map[string]string{"operation": "void"})
}
