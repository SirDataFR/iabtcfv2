package iabtcfv2

import (
	"fmt"
	"strings"
	"time"
)

type TCData struct {
	CoreString       *CoreString
	DisclosedVendors *DisclosedVendors
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

// Returns true if vendor id is allowed for all purpose ids, each one on its legal basis: consent by default, or the
// basis that a publisher restriction which includes the vendor sets for the purpose. See
// CoreString.IsVendorAllowedForFlexiblePurposes.
func (t *TCData) IsVendorAllowedForFlexiblePurposes(id int, purposeIds ...int) bool {
	return t.CoreString.IsVendorAllowedForFlexiblePurposes(id, purposeIds...)
}

// Returns true if vendor id is allowed for all purpose ids, each one on its legal basis: legitimate interest by
// default, or the basis that a publisher restriction which includes the vendor sets for the purpose. See
// CoreString.IsVendorAllowedForFlexiblePurposesLI.
func (t *TCData) IsVendorAllowedForFlexiblePurposesLI(id int, purposeIds ...int) bool {
	return t.CoreString.IsVendorAllowedForFlexiblePurposesLI(id, purposeIds...)
}

// Returns a list of publisher restrictions applied to purpose id
func (t *TCData) GetPubRestrictionsForPurpose(id int) []*PubRestriction {
	return t.CoreString.GetPubRestrictionsForPurpose(id)
}

// IsVendorDisclosed returns true if the given vendor ID is in the DisclosedVendors segment.
func (t *TCData) IsVendorDisclosed(id int) bool {
	if t.DisclosedVendors == nil {
		return false
	}
	return t.DisclosedVendors.IsVendorDisclosed(id)
}

// IsV23 returns true if the TC string uses TCF v2.3 policy (TcfPolicyVersion >= 5).
func (t *TCData) IsV23() bool {
	if t.CoreString == nil {
		return false
	}
	return t.CoreString.TcfPolicyVersion >= TcfPolicyVersion23
}

// Validate checks the TCData for TCF v2.3 compliance.
// Returns an error if the string does not meet v2.3 requirements.
func (t *TCData) Validate() error {
	if t.CoreString == nil {
		return fmt.Errorf("core string is required")
	}

	// After the deadline: TcfPolicyVersion >= 5 AND DisclosedVendors required
	if t.CoreString.Created.After(v23Deadline) {
		if t.CoreString.TcfPolicyVersion < TcfPolicyVersion23 {
			return fmt.Errorf("TCF v2.3: TC Strings created after %s must have policyVersion >= %d", v23Deadline.Format("2006-01-02"), TcfPolicyVersion23)
		}
		if t.DisclosedVendors == nil {
			return fmt.Errorf("TCF v2.3: DisclosedVendors segment is mandatory for TC Strings created after %s", v23Deadline.Format("2006-01-02"))
		}
	}

	return nil
}

// v23Deadline is the IAB TCF v2.3 mandatory adoption deadline
var v23Deadline = time.Unix(TcfV23Deadline, 0).UTC()


// Returns structure as a base64 raw url encoded string
// Encoder produces canonical order: Core → DisclosedVendors → PublisherTC
// Note: Per spec, only CoreString must be first; other segments may appear in any order.
func (t *TCData) ToTCString() string {
	var segments []string

	if t.CoreString != nil {
		segments = append(segments, t.CoreString.Encode())
	}

	// TCF v2.3: DisclosedVendors is mandatory (presence) when policyVersion >= 5
	if t.CoreString != nil && t.CoreString.Version >= int(TcfVersion2) && t.CoreString.TcfPolicyVersion >= TcfPolicyVersion23 {
		if t.DisclosedVendors == nil {
			t.DisclosedVendors = &DisclosedVendors{SegmentType: int(SegmentTypeDisclosedVendors)}
		}
		segments = append(segments, t.DisclosedVendors.Encode())
	} else if t.DisclosedVendors != nil {
		// Pre-v2.3: DisclosedVendors is optional
		segments = append(segments, t.DisclosedVendors.Encode())
	}

	if t.PublisherTC != nil {
		segments = append(segments, t.PublisherTC.Encode())
	}

	return strings.Join(segments, ".")
}
