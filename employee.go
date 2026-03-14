package quickbooks

import (
	"errors"
	"fmt"
	"strconv"
)

// Employee represents a QuickBooks Employee object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type Employee struct {
	ID               string    `json:"Id,omitempty"`
	SyncToken        string    `json:",omitempty"`
	MetaData         *MetaData `json:",omitempty"`
	DisplayName      string    `json:",omitempty"`
	PrintOnCheckName *string   `json:",omitempty"`
	FamilyName       string    `json:",omitempty"`
	GivenName        string    `json:",omitempty"`
	Active           *bool     `json:",omitempty"`
	BillableTime     *bool     `json:",omitempty"`

	// PII fields are intentionally excluded from serialization.
	SSN          *string         `json:"-"`
	PrimaryAddr  *Address        `json:"-"`
	PrimaryPhone TelephoneNumber `json:"-"`
}

// EmployeeCreateInput contains the writable fields accepted when creating an Employee.
// At least one of GivenName, FamilyName, or DisplayName is required.
type EmployeeCreateInput struct {
	GivenName        string  `json:",omitempty"`
	FamilyName       string  `json:",omitempty"`
	DisplayName      string  `json:",omitempty"`
	PrintOnCheckName *string `json:",omitempty"`
	Active           *bool   `json:",omitempty"`
	BillableTime     *bool   `json:",omitempty"`
}

// CreateEmployee creates the given employee within QuickBooks
func (c *Client) CreateEmployee(input *EmployeeCreateInput) (*Employee, error) {
	var resp struct {
		Employee Employee
		Time     Date
	}

	if err := c.post("employee", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Employee, nil
}

// FindEmployees gets the full list of Employees in the QuickBooks account.
func (c *Client) FindEmployees() ([]Employee, error) {
	var resp struct {
		QueryResponse struct {
			Employees     []Employee `json:"Employee"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Employee", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no employees could be found")
	}

	employees := make([]Employee, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Employee ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Employees == nil {
			return nil, errors.New("no employees could be found")
		}

		employees = append(employees, resp.QueryResponse.Employees...)
	}

	return employees, nil
}

// FindEmployeeByID returns an employee with a given Id.
func (c *Client) FindEmployeeByID(id string) (*Employee, error) {
	var resp struct {
		Employee Employee
		Time     Date
	}

	if err := c.get("employee/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Employee, nil
}

// QueryEmployees accepts an SQL query and returns all employees found using it
func (c *Client) QueryEmployees(query string) ([]Employee, error) {
	var resp struct {
		QueryResponse struct {
			Employees     []Employee `json:"Employee"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Employees == nil {
		return nil, errors.New("could not find any employees")
	}

	return resp.QueryResponse.Employees, nil
}

// ListEmployees returns one page of Employees ordered by Id.
// Pass an empty pageToken to start from the beginning.
// The returned nextPageToken is empty when there are no more results.
func (c *Client) ListEmployees(pageToken string, pageSize int) (*ListResponse[Employee], error) {
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
			Employees     []Employee `json:"Employee"`
			StartPosition int
			MaxResults    int
		}
	}

	query := "SELECT * FROM Employee ORDERBY Id STARTPOSITION " + strconv.Itoa(startPosition) + " MAXRESULTS " + strconv.Itoa(pageSize)
	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	result := &ListResponse[Employee]{Items: resp.QueryResponse.Employees}
	if len(result.Items) == pageSize {
		result.NextPageToken = strconv.Itoa(startPosition + pageSize)
	}

	return result, nil
}

// DeleteEmployee deletes the employee.
func (c *Client) DeleteEmployee(employee *Employee) error {
	if employee.ID == "" || employee.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("employee", employee, nil, map[string]string{"operation": "delete"})
}

// UpdateEmployee updates the employee
func (c *Client) UpdateEmployee(employee *Employee) (*Employee, error) {
	if employee.ID == "" {
		return nil, errors.New("missing employee id")
	}

	existingEmployee, err := c.FindEmployeeByID(employee.ID)
	if err != nil {
		return nil, err
	}

	employee.SyncToken = existingEmployee.SyncToken

	payload := struct {
		*Employee
		Sparse bool `json:"sparse"`
	}{
		Employee: employee,
		Sparse:   true,
	}

	var employeeData struct {
		Employee Employee
		Time     Date
	}

	if err = c.post("employee", payload, &employeeData, nil); err != nil {
		return nil, err
	}

	return &employeeData.Employee, err
}
