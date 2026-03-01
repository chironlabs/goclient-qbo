package quickbooks

import (
	"encoding/json"
	"errors"
)

// ExchangeRate represents a QuickBooks ExchangeRate object as returned by the API.
// ExchangeRates are read-only.
type ExchangeRate struct {
	SourceCurrencyCode string      `json:",omitempty"`
	TargetCurrencyCode *string     `json:",omitempty"`
	Rate               json.Number `json:",omitempty"`
	AsOfDate           *Date       `json:",omitempty"`
}

// FindExchangeRate returns the exchange rate for the given source currency code as of the given date.
// asOfDate should be formatted as "YYYY-MM-DD".
func (c *Client) FindExchangeRate(sourceCurrencyCode string, asOfDate string) (*ExchangeRate, error) {
	var resp struct {
		ExchangeRate ExchangeRate
		Time         Date
	}

	params := map[string]string{
		"sourcecurrencycode": sourceCurrencyCode,
	}
	if asOfDate != "" {
		params["asofdate"] = asOfDate
	}

	if err := c.get("exchangerate", &resp, params); err != nil {
		return nil, err
	}

	if resp.ExchangeRate.SourceCurrencyCode == "" {
		return nil, errors.New("exchange rate not found")
	}

	return &resp.ExchangeRate, nil
}

// QueryExchangeRates accepts an SQL query and returns all exchange rates found using it.
// Example: "SELECT * FROM ExchangeRate WHERE SourceCurrencyCode IN ('USD', 'EUR') AND AsOfDate = '2024-01-01'"
func (c *Client) QueryExchangeRates(query string) ([]ExchangeRate, error) {
	var resp struct {
		QueryResponse struct {
			ExchangeRates []ExchangeRate `json:"ExchangeRate"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.ExchangeRates == nil {
		return nil, errors.New("could not find any exchange rates")
	}

	return resp.QueryResponse.ExchangeRates, nil
}
