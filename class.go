package quickbooks

type Class struct {
	ID                 string `json:"Id"`
	FullyQualifiedName string
	Name               string
	SyncToken          string
	ParentRef          *ReferenceType
	SubClass           *bool
	Active             *bool
	MetaData           *MetaData
}

func (c *Client) CreateClass(class Class) (*Class, error) {
	var resp struct {
		Class Class
		Time  Date
	}

	if err := c.post("class", class, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Class, nil
}
