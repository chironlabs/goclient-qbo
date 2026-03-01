package quickbooks

// Class represents a QuickBooks Class object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, FullyQualifiedName) are populated by the service.
type Class struct {
	ID                 string         `json:"Id,omitempty"`
	SyncToken          string         `json:",omitempty"`
	MetaData           *MetaData      `json:",omitempty"`
	Name               string         `json:",omitempty"`
	FullyQualifiedName string         `json:",omitempty"`
	ParentRef          *ReferenceType `json:",omitempty"`
	SubClass           *bool          `json:",omitempty"`
	Active             *bool          `json:",omitempty"`
}

// ClassCreateInput contains the writable fields accepted when creating a Class.
// Name is required; all other fields are optional.
type ClassCreateInput struct {
	Name      string         `json:",omitempty"`
	ParentRef *ReferenceType `json:",omitempty"`
	SubClass  *bool          `json:",omitempty"`
	Active    *bool          `json:",omitempty"`
}

func (c *Client) CreateClass(input *ClassCreateInput) (*Class, error) {
	var resp struct {
		Class Class
		Time  Date
	}

	if err := c.post("class", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Class, nil
}
