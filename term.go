package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Term represents a QuickBooks Term object as returned by the API.
// Read-only fields (Id, SyncToken, MetaData) are populated by the service.
type Term struct {
	ID                 string      `json:"Id,omitempty"`
	SyncToken          string      `json:",omitempty"`
	MetaData           *MetaData   `json:",omitempty"`
	Name               string      `json:",omitempty"`
	Active             *bool       `json:",omitempty"`
	Type               *string     `json:",omitempty"`
	DueDays            *int        `json:",omitempty"`
	DiscountDays       *int        `json:",omitempty"`
	DiscountPercent    json.Number `json:",omitempty"`
	DayOfMonthDue      *int        `json:",omitempty"`
	DueNextMonthDays   *int        `json:",omitempty"`
	DiscountDayOfMonth *int        `json:",omitempty"`
}

// TermCreateInput contains the writable fields accepted when creating a Term.
// Name is required; all other fields are optional.
type TermCreateInput struct {
	Name               string      `json:",omitempty"`
	Active             *bool       `json:",omitempty"`
	Type               *string     `json:",omitempty"`
	DueDays            *int        `json:",omitempty"`
	DiscountDays       *int        `json:",omitempty"`
	DiscountPercent    json.Number `json:",omitempty"`
	DayOfMonthDue      *int        `json:",omitempty"`
	DueNextMonthDays   *int        `json:",omitempty"`
	DiscountDayOfMonth *int        `json:",omitempty"`
}

// CreateTerm creates the given Term on the QuickBooks server, returning
// the resulting Term object.
func (c *Client) CreateTerm(input *TermCreateInput) (*Term, error) {
	var resp struct {
		Term Term
		Time Date
	}

	if err := c.post("term", input, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Term, nil
}

// DeleteTerm deletes the term.
func (c *Client) DeleteTerm(term *Term) error {
	if term.ID == "" || term.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("term", term, nil, map[string]string{"operation": "delete"})
}

// FindTerms gets the full list of Terms in the QuickBooks account.
func (c *Client) FindTerms() ([]Term, error) {
	var resp struct {
		QueryResponse struct {
			Terms         []Term `json:"Term"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Term", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no terms could be found")
	}

	terms := make([]Term, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Term ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Terms == nil {
			return nil, errors.New("no terms could be found")
		}

		terms = append(terms, resp.QueryResponse.Terms...)
	}

	return terms, nil
}

// FindTermByID returns a term with a given Id.
func (c *Client) FindTermByID(id string) (*Term, error) {
	var resp struct {
		Term Term
		Time Date
	}

	if err := c.get("term/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Term, nil
}

// QueryTerms accepts an SQL query and returns all terms found using it.
func (c *Client) QueryTerms(query string) ([]Term, error) {
	var resp struct {
		QueryResponse struct {
			Terms         []Term `json:"Term"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Terms == nil {
		return nil, errors.New("could not find any terms")
	}

	return resp.QueryResponse.Terms, nil
}

// UpdateTerm updates the term.
func (c *Client) UpdateTerm(term *Term) (*Term, error) {
	if term.ID == "" {
		return nil, errors.New("missing term id")
	}

	existingTerm, err := c.FindTermByID(term.ID)
	if err != nil {
		return nil, err
	}

	term.SyncToken = existingTerm.SyncToken

	payload := struct {
		*Term
		Sparse bool `json:"sparse"`
	}{
		Term:   term,
		Sparse: true,
	}

	var termData struct {
		Term Term
		Time Date
	}

	if err = c.post("term", payload, &termData, nil); err != nil {
		return nil, err
	}

	return &termData.Term, err
}
