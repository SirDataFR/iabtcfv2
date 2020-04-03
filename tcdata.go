package iabtcfv2

import "strings"

type TCData struct {
	CoreString       *CoreString
	DisclosedVendors *DisclosedVendors
	AllowedVendors   *AllowedVendors
	PublisherTC      *PublisherTC
}

func (t *TCData) IsPurposeAllowed(id int) bool {
	return t.CoreString.IsPurposeAllowed(id)
}

func (t *TCData) IsPurposeLIAllowed(id int) bool {
	return t.CoreString.IsPurposeLIAllowed(id)
}

func (t *TCData) IsVendorAllowed(id int) bool {
	return t.CoreString.IsVendorAllowed(id)
}

func (t *TCData) IsVendorLIAllowed(id int) bool {
	return t.CoreString.IsVendorLIAllowed(id)
}

func (t *TCData) ToTCString() string {
	var segments []string

	if t.CoreString != nil {
		segments = append(segments, t.CoreString.Encode())
	}
	if t.DisclosedVendors != nil {
		segments = append(segments, t.DisclosedVendors.Encode())
	}
	if t.AllowedVendors != nil {
		segments = append(segments, t.AllowedVendors.Encode())
	}
	if t.PublisherTC != nil {
		segments = append(segments, t.PublisherTC.Encode())
	}

	return strings.Join(segments, ".")
}
