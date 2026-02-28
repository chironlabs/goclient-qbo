package quickbooks

import "time"

type CustomField struct {
	DefinitionID string `json:"DefinitionId,omitempty"`
	StringValue  string `json:"StringValue,omitempty"`
	Type         string `json:"Type,omitempty"`
	Name         string `json:"Name,omitempty"`
}

// Date represents a Quickbooks date
type Date struct {
	time.Time `json:",omitempty"`
}

// UnmarshalJSON removes time from parsed date
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	d.Time, err = time.Parse(format, string(b))
	if err != nil {
		d.Time, err = time.Parse(secondFormat, string(b))
	}

	return err
}

func (d Date) String() string {
	return d.Format(format)
}

// EmailAddress represents a QuickBooks email address.
type EmailAddress struct {
	Address *string `json:",omitempty"`
}

// EndpointURL specifies the endpoint to connect to
type EndpointURL string

const (
	// DiscoveryProductionEndpoint is for live apps.
	DiscoveryProductionEndpoint EndpointURL = "https://developer.api.intuit.com/.well-known/openid_configuration"
	// DiscoverySandboxEndpoint is for testing.
	DiscoverySandboxEndpoint EndpointURL = "https://developer.api.intuit.com/.well-known/openid_sandbox_configuration"
	// ProductionEndpoint is for live apps.
	ProductionEndpoint EndpointURL = "https://quickbooks.api.intuit.com"
	// SandboxEndpoint is for testing.
	SandboxEndpoint EndpointURL = "https://sandbox-quickbooks.api.intuit.com"

	format        = "2006-01-02T15:04:05-07:00"
	queryPageSize = 1000
	secondFormat  = "2006-01-02"
)

func (u EndpointURL) String() string {
	return string(u)
}

// MemoRef represents a QuickBooks MemoRef object.
type MemoRef struct {
	Value string `json:"value,omitempty"`
}

// MetaData is a timestamp of genesis and last change of a Quickbooks object
type MetaData struct {
	CreateTime      Date `json:",omitempty"`
	LastUpdatedTime Date `json:",omitempty"`
}

// Address represents a QuickBooks address.
type Address struct {
	ID string `json:"Id,omitempty"`
	// These lines are context-dependent! Read the QuickBooks API carefully.
	Line1   string  `json:",omitempty"`
	Line2   *string `json:",omitempty"`
	Line3   *string `json:",omitempty"`
	Line4   *string `json:",omitempty"`
	Line5   *string `json:",omitempty"`
	City    string  `json:",omitempty"`
	Country string  `json:",omitempty"`
	// A.K.A. State.
	CountrySubDivisionCode string `json:",omitempty"`
	PostalCode             string `json:",omitempty"`
	Lat                    string `json:",omitempty"`
	Long                   string `json:",omitempty"`
}

type NameValue struct {
	Value string `json:"value,omitempty"`
	Name  string `json:"name,omitempty"`
}

// ReferenceType represents a QuickBooks reference to another object.
type ReferenceType struct {
	NameValue
	Type string `json:"type,omitempty"`
}

// TelephoneNumber represents a QuickBooks phone number.
type TelephoneNumber struct {
	FreeFormNumber string `json:",omitempty"`
}

// WebSiteAddress represents a Quickbooks Website
type WebSiteAddress struct {
	URI string `json:",omitempty"`
}
