package quickbooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Purchase represents a QuickBooks Purchase object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, TotalAmt) are populated by the service.
type Purchase struct {
	ID            string         `json:"Id,omitempty"`
	SyncToken     string         `json:",omitempty"`
	MetaData      *MetaData      `json:",omitempty"`
	AccountRef    ReferenceType  `json:",omitempty"`
	PaymentType   string         `json:",omitempty"`
	Line          []Line         `json:",omitempty"`
	TxnDate       *Date          `json:",omitempty"`
	DocNumber     *string        `json:",omitempty"`
	PrivateNote   *string        `json:",omitempty"`
	TotalAmt      json.Number    `json:",omitempty"`
	EntityRef     *ReferenceType `json:",omitempty"`
	DepartmentRef *ReferenceType `json:",omitempty"`
	CurrencyRef   *ReferenceType `json:",omitempty"`
	ExchangeRate  json.Number    `json:",omitempty"`
	TxnTaxDetail  *TxnTaxDetail  `json:",omitempty"`
	Credit        *bool          `json:",omitempty"`
	PaymentMethodRef *ReferenceType `json:",omitempty"`
}

// PurchaseCreateInput contains the writable fields accepted when creating a Purchase.
// AccountRef, PaymentType, and Line are required; all other fields are optional.
type PurchaseCreateInput struct {
	AccountRef    ReferenceType  `json:",omitempty"`
	PaymentType   string         `json:",omitempty"`
	Line          []Line         `json:",omitempty"`
	TxnDate       *Date          `json:",omitempty"`
	DocNumber     *string        `json:",omitempty"`
	PrivateNote   *string        `json:",omitempty"`
	EntityRef     *ReferenceType `json:",omitempty"`
	DepartmentRef *ReferenceType `json:",omitempty"`
	CurrencyRef   *ReferenceType `json:",omitempty"`
	ExchangeRate  json.Number    `json:",omitempty"`
	TxnTaxDetail  *TxnTaxDetail  `json:",omitempty"`
	Credit        *bool          `json:",omitempty"`
	PaymentMethodRef *ReferenceType `json:",omitempty"`
}

// CreatePurchase creates the given Purchase on the QuickBooks server, returning
// the resulting Purchase object.
func (c *Client) CreatePurchase(input *PurchaseCreateInput) (*Purchase, error) {
	var resp struct {
		Purchase Purchase
		Time     Date
	}

	if err := c.post("purchase", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Purchase, nil
}

// DeletePurchase deletes the purchase.
func (c *Client) DeletePurchase(purchase *Purchase) error {
	if purchase.ID == "" || purchase.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("purchase", purchase, nil, map[string]string{"operation": "delete"})
}

// FindPurchases gets the full list of Purchases in the QuickBooks account.
func (c *Client) FindPurchases() ([]Purchase, error) {
	var resp struct {
		QueryResponse struct {
			Purchases     []Purchase `json:"Purchase"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Purchase", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no purchases could be found")
	}

	purchases := make([]Purchase, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Purchase ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Purchases == nil {
			return nil, errors.New("no purchases could be found")
		}

		purchases = append(purchases, resp.QueryResponse.Purchases...)
	}

	return purchases, nil
}

// FindPurchaseByID finds the purchase by the given id.
func (c *Client) FindPurchaseByID(id string) (*Purchase, error) {
	var resp struct {
		Purchase Purchase
		Time     Date
	}

	if err := c.get("purchase/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Purchase, nil
}

// QueryPurchases accepts an SQL query and returns all purchases found using it.
func (c *Client) QueryPurchases(query string) ([]Purchase, error) {
	var resp struct {
		QueryResponse struct {
			Purchases     []Purchase `json:"Purchase"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Purchases == nil {
		return nil, errors.New("could not find any purchases")
	}

	return resp.QueryResponse.Purchases, nil
}

// ListPurchases returns one page of Purchases ordered by Id.
// Pass an empty pageToken to start from the beginning.
// The returned nextPageToken is empty when there are no more results.
func (c *Client) ListPurchases(pageToken string, pageSize int) (*ListResponse[Purchase], error) {
	if pageSize <= 0 || pageSize > queryPageSize {
		pageSize = queryPageSize
	}

	startPosition := 1
	if pageToken != "" {
		var err error
		startPosition, err = strconv.Atoi(pageToken)
		if err != nil {
			return nil, fmt.Errorf("invalid page token: %v", err)
		}
	}

	var resp struct {
		QueryResponse struct {
			Purchases     []Purchase `json:"Purchase"`
			StartPosition int
			MaxResults    int
		}
	}

	query := "SELECT * FROM Purchase ORDERBY Id STARTPOSITION " + strconv.Itoa(startPosition) + " MAXRESULTS " + strconv.Itoa(pageSize)
	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	result := &ListResponse[Purchase]{Items: resp.QueryResponse.Purchases}
	if len(result.Items) == pageSize {
		result.NextPageToken = strconv.Itoa(startPosition + pageSize)
	}

	return result, nil
}

// UpdatePurchase updates the purchase.
func (c *Client) UpdatePurchase(purchase *Purchase) (*Purchase, error) {
	if purchase.ID == "" {
		return nil, errors.New("missing purchase id")
	}

	existingPurchase, err := c.FindPurchaseByID(purchase.ID)
	if err != nil {
		return nil, err
	}

	purchase.SyncToken = existingPurchase.SyncToken

	payload := struct {
		*Purchase
		Sparse bool `json:"sparse"`
	}{
		Purchase: purchase,
		Sparse:   true,
	}

	var purchaseData struct {
		Purchase Purchase
		Time     Date
	}

	if err = c.post("purchase", payload, &purchaseData, nil); err != nil {
		return nil, err
	}

	return &purchaseData.Purchase, err
}
