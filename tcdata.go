package iabtcfv2

import "strings"

type TCData struct {
	CoreString       *CoreString
	DisclosedVendors *DisclosedVendors
	AllowedVendors   *AllowedVendors
	PublisherTC      *PublisherTC
}

// Returns true if user has given consent to special feature id
func (t *TCData) IsSpecialFeatureAllowed(id int) bool {
	return t.CoreString.IsSpecialFeatureAllowed(id)
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

// Returns true if user has given consent to vendor id processing all purposes ids
// and publisher hasn't set restrictions for them
func (t *TCData) IsVendorAllowedForPurposes(id int, purposeIds ...int) bool {
	return t.CoreString.IsVendorAllowedForPurposes(id, purposeIds...)
}

// Returns true if transparency for vendor id's legitimate interest is established for all purpose ids
// and publisher hasn't set restrictions for them
func (t *TCData) IsVendorAllowedForPurposesLI(id int, purposeIds ...int) bool {
	return t.CoreString.IsVendorAllowedForPurposesLI(id, purposeIds...)
}

// Returns true if user has given consent to vendor id processing all purposes ids
// or if transparency for its legitimate interest is established in accordance with publisher restrictions
func (t *TCData) IsVendorAllowedForFlexiblePurposes(id int, purposeIds ...int) bool {
	return t.CoreString.IsVendorAllowedForFlexiblePurposes(id, purposeIds...)
}

// Returns true if transparency for vendor id's legitimate interest is established for all purpose ids
// or if user has given consent in accordance with publisher restrictions
func (t *TCData) IsVendorAllowedForFlexiblePurposesLI(id int, purposeIds ...int) bool {
	return t.CoreString.IsVendorAllowedForFlexiblePurposesLI(id, purposeIds...)
}

// Returns a list of publisher restrictions applied to purpose id
func (t *TCData) GetPubRestrictionsForPurpose(id int) []*PubRestriction {
	return t.CoreString.GetPubRestrictionsForPurpose(id)
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
