package iabtcfv2

import "strings"

const (
	coreStringType       = 0
	disclosedVendorsType = 1
	allowedVendorsType   = 2
	publicherTCType      = 3
)

type TCData struct {
	CoreString       *CoreString
	DisclosedVendors *DisclosedVendors
	AllowedVendors   *AllowedVendors
	PublisherTC      *PublisherTC
}

// Returns true if user has given consent to purpose id
func (t *TCData) IsPurposeAllowed(id int) bool {
	return t.CoreString.IsPurposeAllowed(id)
}

// Returns true if legitimate interest is established for purpose id
// and user didn't exercise their right to object
func (t *TCData) IsPurposeLIAllowed(id int) bool {
	return t.CoreString.IsPurposeLIAllowed(id)
}

// Returns true if user has given consent to vendor id processing their personal data
func (t *TCData) IsVendorAllowed(id int) bool {
	return t.CoreString.IsVendorAllowed(id)
}

// Returns true if transparency for vendor id's legitimate interest is established
// and user didn't exercise their right to object
func (t *TCData) IsVendorLIAllowed(id int) bool {
	return t.CoreString.IsVendorLIAllowed(id)
}

// Returns structure as a base64 raw url encoded string
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
