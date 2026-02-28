package quickbooks

type CompanyInfo struct {
	ID                        string `json:"Id"`
	SyncToken                 string
	Domain                    string   `json:"domain"`
	LegalAddr                 *Address `json:",omitempty"`
	SupportedLanguages        *string  `json:",omitempty"`
	CompanyName               string   `json:",omitempty"`
	Country                   *string  `json:",omitempty"`
	CompanyAddr               Address
	FiscalYearStartMonth      *string          `json:",omitempty"`
	CustomerCommunicationAddr *Address         `json:",omitempty"`
	PrimaryPhone              *TelephoneNumber `json:",omitempty"`
	LegalName                 *string          `json:",omitempty"`
	CompanyStartDate          *string          `json:",omitempty"`
	EmployerID                *string          `json:"EmployerId,omitempty"`
	Email                     *EmailAddress    `json:",omitempty"`
	WebAddr                   *WebSiteAddress  `json:",omitempty"`
	NameValue                 []NameValue      `json:",omitempty"`
	Metadata                  MetaData
}

// FindCompanyInfo returns the QuickBooks CompanyInfo object. This is a good
// test to check whether you're connected.
func (c *Client) FindCompanyInfo() (*CompanyInfo, error) {
	var resp struct {
		CompanyInfo CompanyInfo
		Time        Date
	}

	if err := c.get("companyinfo/"+c.realm, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.CompanyInfo, nil
}

// UpdateCompanyInfo updates the company info
func (c *Client) UpdateCompanyInfo(companyInfo *CompanyInfo) (*CompanyInfo, error) {
	existingCompanyInfo, err := c.FindCompanyInfo()
	if err != nil {
		return nil, err
	}

	companyInfo.ID = existingCompanyInfo.ID
	companyInfo.SyncToken = existingCompanyInfo.SyncToken

	payload := struct {
		*CompanyInfo
		Sparse bool `json:"sparse"`
	}{
		CompanyInfo: companyInfo,
		Sparse:      true,
	}

	var companyInfoData struct {
		CompanyInfo CompanyInfo
		Time        Date
	}

	if err = c.post("companyInfo", payload, &companyInfoData, nil); err != nil {
		return nil, err
	}

	return &companyInfoData.CompanyInfo, err
}
