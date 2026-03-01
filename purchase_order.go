package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// PurchaseOrder represents a QuickBooks PurchaseOrder object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, TotalAmt) are populated by the service.
type PurchaseOrder struct {
	ID            string         `json:"Id,omitempty"`
	SyncToken     string         `json:",omitempty"`
	MetaData      *MetaData      `json:",omitempty"`
	VendorRef     ReferenceType  `json:",omitempty"`
	APAccountRef  *ReferenceType `json:",omitempty"`
	Line          []Line         `json:",omitempty"`
	TxnDate       *Date          `json:",omitempty"`
	DocNumber     *string        `json:",omitempty"`
	PrivateNote   *string        `json:",omitempty"`
	Memo          *string        `json:",omitempty"`
	POStatus      *string        `json:",omitempty"`
	TotalAmt      json.Number    `json:",omitempty"`
	CurrencyRef   *ReferenceType `json:",omitempty"`
	ExchangeRate  json.Number    `json:",omitempty"`
	ShipAddr      *Address       `json:",omitempty"`
	VendorAddr    *Address       `json:",omitempty"`
	DepartmentRef *ReferenceType `json:",omitempty"`
	ShipMethodRef *ReferenceType `json:",omitempty"`
	TxnTaxDetail  *TxnTaxDetail  `json:",omitempty"`
	EmailStatus   *string        `json:",omitempty"`
	POEmail       *EmailAddress  `json:",omitempty"`
}

// PurchaseOrderCreateInput contains the writable fields accepted when creating a PurchaseOrder.
// VendorRef and Line are required; all other fields are optional.
type PurchaseOrderCreateInput struct {
	VendorRef     ReferenceType  `json:",omitempty"`
	APAccountRef  *ReferenceType `json:",omitempty"`
	Line          []Line         `json:",omitempty"`
	TxnDate       *Date          `json:",omitempty"`
	DocNumber     *string        `json:",omitempty"`
	PrivateNote   *string        `json:",omitempty"`
	Memo          *string        `json:",omitempty"`
	POStatus      *string        `json:",omitempty"`
	CurrencyRef   *ReferenceType `json:",omitempty"`
	ExchangeRate  json.Number    `json:",omitempty"`
	ShipAddr      *Address       `json:",omitempty"`
	VendorAddr    *Address       `json:",omitempty"`
	DepartmentRef *ReferenceType `json:",omitempty"`
	ShipMethodRef *ReferenceType `json:",omitempty"`
	TxnTaxDetail  *TxnTaxDetail  `json:",omitempty"`
	EmailStatus   *string        `json:",omitempty"`
	POEmail       *EmailAddress  `json:",omitempty"`
}

// CreatePurchaseOrder creates the given PurchaseOrder on the QuickBooks server, returning
// the resulting PurchaseOrder object.
func (c *Client) CreatePurchaseOrder(input *PurchaseOrderCreateInput) (*PurchaseOrder, error) {
	var resp struct {
		PurchaseOrder PurchaseOrder
		Time          Date
	}

	if err := c.post("purchaseorder", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.PurchaseOrder, nil
}

// DeletePurchaseOrder deletes the purchase order.
func (c *Client) DeletePurchaseOrder(purchaseOrder *PurchaseOrder) error {
	if purchaseOrder.ID == "" || purchaseOrder.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("purchaseorder", purchaseOrder, nil, map[string]string{"operation": "delete"})
}

// FindPurchaseOrders gets the full list of PurchaseOrders in the QuickBooks account.
func (c *Client) FindPurchaseOrders() ([]PurchaseOrder, error) {
	var resp struct {
		QueryResponse struct {
			PurchaseOrders []PurchaseOrder `json:"PurchaseOrder"`
			MaxResults     int
			StartPosition  int
			TotalCount     int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM PurchaseOrder", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no purchase orders could be found")
	}

	purchaseOrders := make([]PurchaseOrder, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM PurchaseOrder ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.PurchaseOrders == nil {
			return nil, errors.New("no purchase orders could be found")
		}

		purchaseOrders = append(purchaseOrders, resp.QueryResponse.PurchaseOrders...)
	}

	return purchaseOrders, nil
}

// FindPurchaseOrderByID finds the purchase order by the given id.
func (c *Client) FindPurchaseOrderByID(id string) (*PurchaseOrder, error) {
	var resp struct {
		PurchaseOrder PurchaseOrder
		Time          Date
	}

	if err := c.get("purchaseorder/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.PurchaseOrder, nil
}

// QueryPurchaseOrders accepts an SQL query and returns all purchase orders found using it.
func (c *Client) QueryPurchaseOrders(query string) ([]PurchaseOrder, error) {
	var resp struct {
		QueryResponse struct {
			PurchaseOrders []PurchaseOrder `json:"PurchaseOrder"`
			StartPosition  int
			MaxResults     int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.PurchaseOrders == nil {
		return nil, errors.New("could not find any purchase orders")
	}

	return resp.QueryResponse.PurchaseOrders, nil
}

// UpdatePurchaseOrder updates the purchase order.
func (c *Client) UpdatePurchaseOrder(purchaseOrder *PurchaseOrder) (*PurchaseOrder, error) {
	if purchaseOrder.ID == "" {
		return nil, errors.New("missing purchase order id")
	}

	existingPurchaseOrder, err := c.FindPurchaseOrderByID(purchaseOrder.ID)
	if err != nil {
		return nil, err
	}

	purchaseOrder.SyncToken = existingPurchaseOrder.SyncToken

	payload := struct {
		*PurchaseOrder
		Sparse bool `json:"sparse"`
	}{
		PurchaseOrder: purchaseOrder,
		Sparse:        true,
	}

	var purchaseOrderData struct {
		PurchaseOrder PurchaseOrder
		Time          Date
	}

	if err = c.post("purchaseorder", payload, &purchaseOrderData, nil); err != nil {
		return nil, err
	}

	return &purchaseOrderData.PurchaseOrder, err
}
