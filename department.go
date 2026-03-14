package quickbooks

import (
	"errors"
	"fmt"
	"strconv"
)

// Department represents a QuickBooks Department object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData, FullyQualifiedName) are populated by the service.
type Department struct {
	ID                 string         `json:"Id,omitempty"`
	SyncToken          string         `json:",omitempty"`
	MetaData           *MetaData      `json:",omitempty"`
	Name               string         `json:",omitempty"`
	FullyQualifiedName string         `json:",omitempty"`
	ParentRef          *ReferenceType `json:",omitempty"`
	SubDepartment      *bool          `json:",omitempty"`
	Active             *bool          `json:",omitempty"`
}

// DepartmentCreateInput contains the writable fields accepted when creating a Department.
// Name is required; all other fields are optional.
type DepartmentCreateInput struct {
	Name          string         `json:",omitempty"`
	ParentRef     *ReferenceType `json:",omitempty"`
	SubDepartment *bool          `json:",omitempty"`
	Active        *bool          `json:",omitempty"`
}

// CreateDepartment creates the given Department on the QuickBooks server, returning
// the resulting Department object.
func (c *Client) CreateDepartment(input *DepartmentCreateInput) (*Department, error) {
	var resp struct {
		Department Department
		Time       Date
	}

	if err := c.post("department", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Department, nil
}

// DeleteDepartment deletes the department.
func (c *Client) DeleteDepartment(department *Department) error {
	if department.ID == "" || department.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("department", department, nil, map[string]string{"operation": "delete"})
}

// FindDepartments gets the full list of Departments in the QuickBooks account.
func (c *Client) FindDepartments() ([]Department, error) {
	var resp struct {
		QueryResponse struct {
			Departments   []Department `json:"Department"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Department", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no departments could be found")
	}

	departments := make([]Department, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Department ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Departments == nil {
			return nil, errors.New("no departments could be found")
		}

		departments = append(departments, resp.QueryResponse.Departments...)
	}

	return departments, nil
}

// FindDepartmentByID returns a department with a given Id.
func (c *Client) FindDepartmentByID(id string) (*Department, error) {
	var resp struct {
		Department Department
		Time       Date
	}

	if err := c.get("department/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Department, nil
}

// QueryDepartments accepts an SQL query and returns all departments found using it.
func (c *Client) QueryDepartments(query string) ([]Department, error) {
	var resp struct {
		QueryResponse struct {
			Departments   []Department `json:"Department"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Departments == nil {
		return nil, errors.New("could not find any departments")
	}

	return resp.QueryResponse.Departments, nil
}

// ListDepartments returns one page of Departments ordered by Id.
// Pass an empty pageToken to start from the beginning.
// The returned nextPageToken is empty when there are no more results.
func (c *Client) ListDepartments(pageToken string, pageSize int) (*ListResponse[Department], error) {
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
			Departments   []Department `json:"Department"`
			StartPosition int
			MaxResults    int
		}
	}

	query := "SELECT * FROM Department ORDERBY Id STARTPOSITION " + strconv.Itoa(startPosition) + " MAXRESULTS " + strconv.Itoa(pageSize)
	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	result := &ListResponse[Department]{Items: resp.QueryResponse.Departments}
	if len(result.Items) == pageSize {
		result.NextPageToken = strconv.Itoa(startPosition + pageSize)
	}

	return result, nil
}

// UpdateDepartment updates the department.
func (c *Client) UpdateDepartment(department *Department) (*Department, error) {
	if department.ID == "" {
		return nil, errors.New("missing department id")
	}

	existingDepartment, err := c.FindDepartmentByID(department.ID)
	if err != nil {
		return nil, err
	}

	department.SyncToken = existingDepartment.SyncToken

	payload := struct {
		*Department
		Sparse bool `json:"sparse"`
	}{
		Department: department,
		Sparse:     true,
	}

	var departmentData struct {
		Department Department
		Time       Date
	}

	if err = c.post("department", payload, &departmentData, nil); err != nil {
		return nil, err
	}

	return &departmentData.Department, err
}
