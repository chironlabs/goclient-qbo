package quickbooks

type CompanyInfo struct {
	Id                        string
	SyncToken                 string
	Domain                    string `json:"domain"`
	LegalAddr                 *Address
	SupportedLanguages        *string
	CompanyName               string
	Country                   *string
	CompanyAddr               Address
	ID                        string `json:"Id"`
	FiscalYearStartMonth      *string
	CustomerCommunicationAddr *Address
	PrimaryPhone              *struct {
		FreeFormNumber *string
	}
	LegalName        *string
	CompanyStartDate *string
	EmployerID       *string `json:"EmployerId"`
	Email            *struct {
		Address string
	}
	WebAddr *struct {
		URI *string `json:",omitempty"`
	}
	NameValue []struct {
		Name  string
		Value string
	} `json:",omitempty"`
	Metadata MetaData
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

	companyInfo.Id = existingCompanyInfo.Id
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
