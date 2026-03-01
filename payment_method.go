package quickbooks

import (
	"errors"
	"strconv"
)

// PaymentMethod represents a QuickBooks PaymentMethod object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type PaymentMethod struct {
	ID        string    `json:"Id,omitempty"`
	SyncToken string    `json:",omitempty"`
	MetaData  *MetaData `json:",omitempty"`
	Name      string    `json:",omitempty"`
	Type      *string   `json:",omitempty"`
	Active    *bool     `json:",omitempty"`
}

// PaymentMethodCreateInput contains the writable fields accepted when creating a PaymentMethod.
// Name is required; all other fields are optional.
type PaymentMethodCreateInput struct {
	Name   string  `json:",omitempty"`
	Type   *string `json:",omitempty"`
	Active *bool   `json:",omitempty"`
}

// CreatePaymentMethod creates the given PaymentMethod on the QuickBooks server, returning
// the resulting PaymentMethod object.
func (c *Client) CreatePaymentMethod(input *PaymentMethodCreateInput) (*PaymentMethod, error) {
	var resp struct {
		PaymentMethod PaymentMethod
		Time          Date
	}

	if err := c.post("paymentmethod", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.PaymentMethod, nil
}

// DeletePaymentMethod deletes the payment method.
func (c *Client) DeletePaymentMethod(paymentMethod *PaymentMethod) error {
	if paymentMethod.ID == "" || paymentMethod.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("paymentmethod", paymentMethod, nil, map[string]string{"operation": "delete"})
}

// FindPaymentMethods gets the full list of PaymentMethods in the QuickBooks account.
func (c *Client) FindPaymentMethods() ([]PaymentMethod, error) {
	var resp struct {
		QueryResponse struct {
			PaymentMethods []PaymentMethod `json:"PaymentMethod"`
			MaxResults     int
			StartPosition  int
			TotalCount     int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM PaymentMethod", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no payment methods could be found")
	}

	paymentMethods := make([]PaymentMethod, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM PaymentMethod ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.PaymentMethods == nil {
			return nil, errors.New("no payment methods could be found")
		}

		paymentMethods = append(paymentMethods, resp.QueryResponse.PaymentMethods...)
	}

	return paymentMethods, nil
}

// FindPaymentMethodByID returns a payment method with a given Id.
func (c *Client) FindPaymentMethodByID(id string) (*PaymentMethod, error) {
	var resp struct {
		PaymentMethod PaymentMethod
		Time          Date
	}

	if err := c.get("paymentmethod/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.PaymentMethod, nil
}

// QueryPaymentMethods accepts an SQL query and returns all payment methods found using it.
func (c *Client) QueryPaymentMethods(query string) ([]PaymentMethod, error) {
	var resp struct {
		QueryResponse struct {
			PaymentMethods []PaymentMethod `json:"PaymentMethod"`
			StartPosition  int
			MaxResults     int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.PaymentMethods == nil {
		return nil, errors.New("could not find any payment methods")
	}

	return resp.QueryResponse.PaymentMethods, nil
}

// UpdatePaymentMethod updates the payment method.
func (c *Client) UpdatePaymentMethod(paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	if paymentMethod.ID == "" {
		return nil, errors.New("missing payment method id")
	}

	existingPaymentMethod, err := c.FindPaymentMethodByID(paymentMethod.ID)
	if err != nil {
		return nil, err
	}

	paymentMethod.SyncToken = existingPaymentMethod.SyncToken

	payload := struct {
		*PaymentMethod
		Sparse bool `json:"sparse"`
	}{
		PaymentMethod: paymentMethod,
		Sparse:        true,
	}

	var paymentMethodData struct {
		PaymentMethod PaymentMethod
		Time          Date
	}

	if err = c.post("paymentmethod", payload, &paymentMethodData, nil); err != nil {
		return nil, err
	}

	return &paymentMethodData.PaymentMethod, err
}
