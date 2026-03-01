package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Transfer represents a QuickBooks Transfer object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, TotalAmt) are populated by the service.
type Transfer struct {
	ID             string         `json:"Id,omitempty"`
	SyncToken      string         `json:",omitempty"`
	MetaData       *MetaData      `json:",omitempty"`
	FromAccountRef ReferenceType  `json:",omitempty"`
	ToAccountRef   ReferenceType  `json:",omitempty"`
	Amount         json.Number    `json:",omitempty"`
	TxnDate        *Date          `json:",omitempty"`
	PrivateNote    *string        `json:",omitempty"`
	TxnTaxDetail   *TxnTaxDetail  `json:",omitempty"`
	CurrencyRef    *ReferenceType `json:",omitempty"`
	ExchangeRate   json.Number    `json:",omitempty"`
	TotalAmt       json.Number    `json:",omitempty"`
}

// TransferCreateInput contains the writable fields accepted when creating a Transfer.
// FromAccountRef, ToAccountRef, and Amount are required; all other fields are optional.
type TransferCreateInput struct {
	FromAccountRef ReferenceType  `json:",omitempty"`
	ToAccountRef   ReferenceType  `json:",omitempty"`
	Amount         json.Number    `json:",omitempty"`
	TxnDate        *Date          `json:",omitempty"`
	PrivateNote    *string        `json:",omitempty"`
	TxnTaxDetail   *TxnTaxDetail  `json:",omitempty"`
	CurrencyRef    *ReferenceType `json:",omitempty"`
	ExchangeRate   json.Number    `json:",omitempty"`
}

// CreateTransfer creates the given Transfer on the QuickBooks server, returning
// the resulting Transfer object.
func (c *Client) CreateTransfer(input *TransferCreateInput) (*Transfer, error) {
	var resp struct {
		Transfer Transfer
		Time     Date
	}

	if err := c.post("transfer", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Transfer, nil
}

// DeleteTransfer deletes the transfer.
func (c *Client) DeleteTransfer(transfer *Transfer) error {
	if transfer.ID == "" || transfer.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("transfer", transfer, nil, map[string]string{"operation": "delete"})
}

// FindTransfers gets the full list of Transfers in the QuickBooks account.
func (c *Client) FindTransfers() ([]Transfer, error) {
	var resp struct {
		QueryResponse struct {
			Transfers     []Transfer `json:"Transfer"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Transfer", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no transfers could be found")
	}

	transfers := make([]Transfer, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Transfer ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Transfers == nil {
			return nil, errors.New("no transfers could be found")
		}

		transfers = append(transfers, resp.QueryResponse.Transfers...)
	}

	return transfers, nil
}

// FindTransferByID finds the transfer by the given id.
func (c *Client) FindTransferByID(id string) (*Transfer, error) {
	var resp struct {
		Transfer Transfer
		Time     Date
	}

	if err := c.get("transfer/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Transfer, nil
}

// QueryTransfers accepts an SQL query and returns all transfers found using it.
func (c *Client) QueryTransfers(query string) ([]Transfer, error) {
	var resp struct {
		QueryResponse struct {
			Transfers     []Transfer `json:"Transfer"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Transfers == nil {
		return nil, errors.New("could not find any transfers")
	}

	return resp.QueryResponse.Transfers, nil
}

// UpdateTransfer updates the transfer.
func (c *Client) UpdateTransfer(transfer *Transfer) (*Transfer, error) {
	if transfer.ID == "" {
		return nil, errors.New("missing transfer id")
	}

	existingTransfer, err := c.FindTransferByID(transfer.ID)
	if err != nil {
		return nil, err
	}

	transfer.SyncToken = existingTransfer.SyncToken

	payload := struct {
		*Transfer
		Sparse bool `json:"sparse"`
	}{
		Transfer: transfer,
		Sparse:   true,
	}

	var transferData struct {
		Transfer Transfer
		Time     Date
	}

	if err = c.post("transfer", payload, &transferData, nil); err != nil {
		return nil, err
	}

	return &transferData.Transfer, err
}
