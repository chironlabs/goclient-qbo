package quickbooks

// Preferences represents a QuickBooks Preferences object as returned by the
// API. Only the currency preferences are mapped; the other sections are
// ignored when decoding. Read-only fields (Id, SyncToken, MetaData) are
// populated by the service.
type Preferences struct {
	ID            string         `json:"Id,omitempty"`
	SyncToken     string         `json:",omitempty"`
	CurrencyPrefs *CurrencyPrefs `json:",omitempty"`
	MetaData      *MetaData      `json:",omitempty"`
}

// CurrencyPrefs holds a company's currency preferences.
type CurrencyPrefs struct {
	MultiCurrencyEnabled *bool `json:",omitempty"`
	// HomeCurrency's Value is the ISO 4217 code of the company's home
	// currency ("USD").
	HomeCurrency *ReferenceType `json:",omitempty"`
}

// FindPreferences returns the company's preferences.
//
//	prefs, err := client.FindPreferences()
//	home := prefs.CurrencyPrefs.HomeCurrency.Value // "USD"
func (c *Client) FindPreferences() (*Preferences, error) {
	var resp struct {
		Preferences Preferences
		Time        Date
	}

	if err := c.get("preferences", &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Preferences, nil
}
