package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// RefundReceipt represents a QuickBooks RefundReceipt object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, TotalAmt, Balance) are populated by the service.
type RefundReceipt struct {
	ID                    string         `json:"Id,omitempty"`
	SyncToken             string         `json:",omitempty"`
	MetaData              *MetaData      `json:",omitempty"`
	CustomerRef           *ReferenceType `json:",omitempty"`
	DepositToAccountRef   *ReferenceType `json:",omitempty"`
	PaymentMethodRef      *ReferenceType `json:",omitempty"`
	Line                  []Line         `json:",omitempty"`
	TxnDate               *Date          `json:",omitempty"`
	DocNumber             *string        `json:",omitempty"`
	PrivateNote           *string        `json:",omitempty"`
	CustomerMemo          *MemoRef       `json:",omitempty"`
	BillAddr              *Address       `json:",omitempty"`
	ShipAddr              *Address       `json:",omitempty"`
	ClassRef              *ReferenceType `json:",omitempty"`
	DepartmentRef         *ReferenceType `json:",omitempty"`
	CurrencyRef           *ReferenceType `json:",omitempty"`
	ExchangeRate          json.Number    `json:",omitempty"`
	ApplyTaxAfterDiscount *bool          `json:",omitempty"`
	PrintStatus           *string        `json:",omitempty"`
	EmailStatus           *string        `json:",omitempty"`
	BillEmail             *EmailAddress  `json:",omitempty"`
	TxnTaxDetail          *TxnTaxDetail  `json:",omitempty"`
	CustomField           []CustomField  `json:",omitempty"`
	TotalAmt              json.Number    `json:",omitempty"`
	Balance               json.Number    `json:",omitempty"`
}

// RefundReceiptCreateInput contains the writable fields accepted when creating a RefundReceipt.
// Line is required; all other fields are optional.
type RefundReceiptCreateInput struct {
	Line                  []Line         `json:",omitempty"`
	CustomerRef           *ReferenceType `json:",omitempty"`
	DepositToAccountRef   *ReferenceType `json:",omitempty"`
	PaymentMethodRef      *ReferenceType `json:",omitempty"`
	TxnDate               *Date          `json:",omitempty"`
	DocNumber             *string        `json:",omitempty"`
	PrivateNote           *string        `json:",omitempty"`
	CustomerMemo          *MemoRef       `json:",omitempty"`
	BillAddr              *Address       `json:",omitempty"`
	ShipAddr              *Address       `json:",omitempty"`
	ClassRef              *ReferenceType `json:",omitempty"`
	DepartmentRef         *ReferenceType `json:",omitempty"`
	CurrencyRef           *ReferenceType `json:",omitempty"`
	ExchangeRate          json.Number    `json:",omitempty"`
	ApplyTaxAfterDiscount *bool          `json:",omitempty"`
	PrintStatus           *string        `json:",omitempty"`
	EmailStatus           *string        `json:",omitempty"`
	BillEmail             *EmailAddress  `json:",omitempty"`
	TxnTaxDetail          *TxnTaxDetail  `json:",omitempty"`
	CustomField           []CustomField  `json:",omitempty"`
}

// CreateRefundReceipt creates the given RefundReceipt on the QuickBooks server, returning
// the resulting RefundReceipt object.
func (c *Client) CreateRefundReceipt(input *RefundReceiptCreateInput) (*RefundReceipt, error) {
	var resp struct {
		RefundReceipt RefundReceipt
		Time          Date
	}

	if err := c.post("refundreceipt", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.RefundReceipt, nil
}

// DeleteRefundReceipt deletes the refund receipt.
func (c *Client) DeleteRefundReceipt(refundReceipt *RefundReceipt) error {
	if refundReceipt.ID == "" || refundReceipt.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("refundreceipt", refundReceipt, nil, map[string]string{"operation": "delete"})
}

// FindRefundReceipts gets the full list of RefundReceipts in the QuickBooks account.
func (c *Client) FindRefundReceipts() ([]RefundReceipt, error) {
	var resp struct {
		QueryResponse struct {
			RefundReceipts []RefundReceipt `json:"RefundReceipt"`
			MaxResults     int
			StartPosition  int
			TotalCount     int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM RefundReceipt", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no refund receipts could be found")
	}

	refundReceipts := make([]RefundReceipt, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM RefundReceipt ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.RefundReceipts == nil {
			return nil, errors.New("no refund receipts could be found")
		}

		refundReceipts = append(refundReceipts, resp.QueryResponse.RefundReceipts...)
	}

	return refundReceipts, nil
}

// FindRefundReceiptByID finds the refund receipt by the given id.
func (c *Client) FindRefundReceiptByID(id string) (*RefundReceipt, error) {
	var resp struct {
		RefundReceipt RefundReceipt
		Time          Date
	}

	if err := c.get("refundreceipt/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.RefundReceipt, nil
}

// QueryRefundReceipts accepts an SQL query and returns all refund receipts found using it.
func (c *Client) QueryRefundReceipts(query string) ([]RefundReceipt, error) {
	var resp struct {
		QueryResponse struct {
			RefundReceipts []RefundReceipt `json:"RefundReceipt"`
			StartPosition  int
			MaxResults     int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.RefundReceipts == nil {
		return nil, errors.New("could not find any refund receipts")
	}

	return resp.QueryResponse.RefundReceipts, nil
}

// UpdateRefundReceipt updates the refund receipt.
func (c *Client) UpdateRefundReceipt(refundReceipt *RefundReceipt) (*RefundReceipt, error) {
	if refundReceipt.ID == "" {
		return nil, errors.New("missing refund receipt id")
	}

	existingRefundReceipt, err := c.FindRefundReceiptByID(refundReceipt.ID)
	if err != nil {
		return nil, err
	}

	refundReceipt.SyncToken = existingRefundReceipt.SyncToken

	payload := struct {
		*RefundReceipt
		Sparse bool `json:"sparse"`
	}{
		RefundReceipt: refundReceipt,
		Sparse:        true,
	}

	var refundReceiptData struct {
		RefundReceipt RefundReceipt
		Time          Date
	}

	if err = c.post("refundreceipt", payload, &refundReceiptData, nil); err != nil {
		return nil, err
	}

	return &refundReceiptData.RefundReceipt, err
}

// VoidRefundReceipt voids the refund receipt.
func (c *Client) VoidRefundReceipt(refundReceipt *RefundReceipt) error {
	if refundReceipt.ID == "" {
		return errors.New("missing refund receipt id")
	}

	existingRefundReceipt, err := c.FindRefundReceiptByID(refundReceipt.ID)
	if err != nil {
		return err
	}

	refundReceipt.SyncToken = existingRefundReceipt.SyncToken

	return c.post("refundreceipt", refundReceipt, nil, map[string]string{"operation": "void"})
}
