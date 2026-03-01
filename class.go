package quickbooks

import (
	"errors"
	"strconv"
)

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

// CreateClass creates the given Class on the QuickBooks server, returning
// the resulting Class object.
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

// DeleteClass deletes the class.
func (c *Client) DeleteClass(class *Class) error {
	if class.ID == "" || class.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("class", class, nil, map[string]string{"operation": "delete"})
}

// FindClasses gets the full list of Classes in the QuickBooks account.
func (c *Client) FindClasses() ([]Class, error) {
	var resp struct {
		QueryResponse struct {
			Classes       []Class `json:"Class"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Class", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no classes could be found")
	}

	classes := make([]Class, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Class ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Classes == nil {
			return nil, errors.New("no classes could be found")
		}

		classes = append(classes, resp.QueryResponse.Classes...)
	}

	return classes, nil
}

// FindClassByID returns a class with a given Id.
func (c *Client) FindClassByID(id string) (*Class, error) {
	var resp struct {
		Class Class
		Time  Date
	}

	if err := c.get("class/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Class, nil
}

// QueryClasses accepts an SQL query and returns all classes found using it.
func (c *Client) QueryClasses(query string) ([]Class, error) {
	var resp struct {
		QueryResponse struct {
			Classes       []Class `json:"Class"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Classes == nil {
		return nil, errors.New("could not find any classes")
	}

	return resp.QueryResponse.Classes, nil
}

// UpdateClass updates the class.
func (c *Client) UpdateClass(class *Class) (*Class, error) {
	if class.ID == "" {
		return nil, errors.New("missing class id")
	}

	existingClass, err := c.FindClassByID(class.ID)
	if err != nil {
		return nil, err
	}

	class.SyncToken = existingClass.SyncToken

	payload := struct {
		*Class
		Sparse bool `json:"sparse"`
	}{
		Class:  class,
		Sparse: true,
	}

	var classData struct {
		Class Class
		Time  Date
	}

	if err = c.post("class", payload, &classData, nil); err != nil {
		return nil, err
	}

	return &classData.Class, err
}
