package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Item represents a QuickBooks Item object as returned by the API (a product or service).
// Read-only fields (Id, SyncToken, MetaData, QtyOnHand) are populated by the service.
type Item struct {
	ID          string      `json:"Id,omitempty"`
	SyncToken   string      `json:",omitempty"`
	MetaData    *MetaData   `json:",omitempty"`
	Name        string
	SKU         *string     `json:"Sku,omitempty"`
	Description *string     `json:",omitempty"`
	Active      *bool       `json:",omitempty"`
	Taxable     *bool       `json:",omitempty"`
	SalesTaxIncluded    *bool          `json:",omitempty"`
	UnitPrice           json.Number    `json:",omitempty"`
	Type                string
	IncomeAccountRef    *ReferenceType `json:",omitempty"`
	ExpenseAccountRef   *ReferenceType `json:",omitempty"`
	PurchaseDesc        *string        `json:",omitempty"`
	PurchaseTaxIncluded *bool          `json:",omitempty"`
	PurchaseCost        json.Number    `json:",omitempty"`
	AssetAccountRef     *ReferenceType `json:",omitempty"`
	TrackQtyOnHand      *bool          `json:",omitempty"`
	QtyOnHand           json.Number    `json:",omitempty"`
	SalesTaxCodeRef     *ReferenceType `json:",omitempty"`
	PurchaseTaxCodeRef  *ReferenceType `json:",omitempty"`
}

// ItemCreateInput contains the writable fields accepted when creating an Item.
// Name and Type are required; account refs are conditionally required based on Type.
type ItemCreateInput struct {
	Name                string         `json:",omitempty"`
	Type                string         `json:",omitempty"`
	SKU                 *string        `json:"Sku,omitempty"`
	Description         *string        `json:",omitempty"`
	Active              *bool          `json:",omitempty"`
	Taxable             *bool          `json:",omitempty"`
	SalesTaxIncluded    *bool          `json:",omitempty"`
	UnitPrice           json.Number    `json:",omitempty"`
	IncomeAccountRef    *ReferenceType `json:",omitempty"`
	ExpenseAccountRef   *ReferenceType `json:",omitempty"`
	PurchaseDesc        *string        `json:",omitempty"`
	PurchaseTaxIncluded *bool          `json:",omitempty"`
	PurchaseCost        json.Number    `json:",omitempty"`
	AssetAccountRef     *ReferenceType `json:",omitempty"`
	TrackQtyOnHand      *bool          `json:",omitempty"`
	SalesTaxCodeRef     *ReferenceType `json:",omitempty"`
	PurchaseTaxCodeRef  *ReferenceType `json:",omitempty"`
}

// CreateItem creates the given Item on the QuickBooks server, returning
// the resulting Item object.
func (c *Client) CreateItem(input *ItemCreateInput) (*Item, error) {
	var resp struct {
		Item Item
		Time Date
	}

	if err := c.post("item", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Item, nil
}

// FindItems gets the full list of Items in the QuickBooks account.
func (c *Client) FindItems() ([]Item, error) {
	var resp struct {
		QueryResponse struct {
			Items         []Item `json:"Item"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Item", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no items could be found")
	}

	items := make([]Item, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Item ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Items == nil {
			return nil, errors.New("no items could be found")
		}

		items = append(items, resp.QueryResponse.Items...)
	}

	return items, nil
}

// FindItemByID returns an item with a given Id.
func (c *Client) FindItemByID(id string) (*Item, error) {
	var resp struct {
		Item Item
		Time Date
	}

	if err := c.get("item/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Item, nil
}

// QueryItems accepts an SQL query and returns all items found using it
func (c *Client) QueryItems(query string) ([]Item, error) {
	var resp struct {
		QueryResponse struct {
			Items         []Item `json:"Item"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Items == nil {
		return nil, errors.New("could not find any items")
	}

	return resp.QueryResponse.Items, nil
}

// DeleteItem deletes the item.
func (c *Client) DeleteItem(item *Item) error {
	if item.ID == "" || item.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("item", item, nil, map[string]string{"operation": "delete"})
}

// UpdateItem updates the item
func (c *Client) UpdateItem(item *Item) (*Item, error) {
	if item.ID == "" {
		return nil, errors.New("missing item id")
	}

	existingItem, err := c.FindItemByID(item.ID)
	if err != nil {
		return nil, err
	}

	item.SyncToken = existingItem.SyncToken

	payload := struct {
		*Item
		Sparse bool `json:"sparse"`
	}{
		Item:   item,
		Sparse: true,
	}

	var itemData struct {
		Item Item
		Time Date
	}

	if err = c.post("item", payload, &itemData, nil); err != nil {
		return nil, err
	}

	return &itemData.Item, err
}
