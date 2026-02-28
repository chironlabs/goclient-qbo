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
