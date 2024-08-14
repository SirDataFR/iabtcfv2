package iabtcfv2

import (
	"encoding/base64"
	"time"
)

type CoreString struct {
	Version                int
	Created                time.Time
	LastUpdated            time.Time
	CmpId                  int
	CmpVersion             int
	ConsentScreen          int
	ConsentLanguage        string
	VendorListVersion      int
	TcfPolicyVersion       int
	IsServiceSpecific      bool
	UseNonStandardTexts    bool
	SpecialFeatureOptIns   map[int]bool
	PurposesConsent        map[int]bool
	PurposesLITransparency map[int]bool
	PurposeOneTreatment    bool
	PublisherCC            string
	MaxVendorId            int
	IsRangeEncoding        bool
	VendorsConsent         map[int]bool
	NumEntries             int
	RangeEntries           []*RangeEntry
	MaxVendorIdLI          int
	IsRangeEncodingLI      bool
	VendorsLITransparency  map[int]bool
	NumEntriesLI           int
	RangeEntriesLI         []*RangeEntry
	NumPubRestrictions     int
	PubRestrictions        []*PubRestriction
}

type PubRestriction struct {
	PurposeId       int
	RestrictionType RestrictionType
	NumEntries      int
	RangeEntries    []*RangeEntry
}

func (r *PubRestriction) getBitSize() (bitSize int) {
	bitSize += bitsPubRestrictionsEntryPurposeId
	bitSize += bitsPubRestrictionsEntryRestrictionType
	bitSize += bitsNumEntries
	for _, entry := range r.RangeEntries {
		bitSize += entry.getBitSize()
	}
	return bitSize
}

type RangeEntry struct {
	StartVendorID int
	EndVendorID   int
}

func (r *RangeEntry) getBitSize() (bitSize int) {
	bitSize += bitsIsRangeEncoding
	if r.EndVendorID > r.StartVendorID {
		bitSize += bitsVendorId * 2
	} else {
		bitSize += bitsVendorId
	}
	return bitSize
}

// Returns true if user has given consent to special feature id
func (c *CoreString) IsSpecialFeatureAllowed(id int) bool {
	return c.SpecialFeatureOptIns[id]
}

// Returns true if user has given consent to purpose id
func (c *CoreString) IsPurposeAllowed(id int) bool {
	return c.PurposesConsent[id]
}

// Returns true if legitimate interest is established for purpose id
// and user didn't exercise their right to object
func (c *CoreString) IsPurposeLIAllowed(id int) bool {
	return c.PurposesLITransparency[id]
}

// Returns true if user has given consent to vendor id processing their personal data
func (c *CoreString) IsVendorAllowed(id int) bool {
	if c.IsRangeEncoding {
		for _, entry := range c.RangeEntries {
			if entry.StartVendorID <= id && id <= entry.EndVendorID {
				return true
			}
		}
		return false
	}

	return c.VendorsConsent[id]
}

// Returns true if transparency for vendor id's legitimate interest is established
// and user didn't exercise their right to object
func (c *CoreString) IsVendorLIAllowed(id int) bool {
	if c.IsRangeEncodingLI {
		for _, entry := range c.RangeEntriesLI {
			if entry.StartVendorID <= id && id <= entry.EndVendorID {
				return true
			}
		}
		return false
	}

	return c.VendorsLITransparency[id]
}

// Returns true if user has given consent to vendor id processing all purposes ids
// and publisher hasn't set restrictions for them
func (c *CoreString) IsVendorAllowedForPurposes(id int, purposeIds ...int) bool {
	if !c.IsVendorAllowed(id) {
		return false
	}

	for _, p := range purposeIds {
		if !c.IsPurposeAllowed(p) {
			return false
		}
	}

	for _, p := range purposeIds {
		pr := c.GetPubRestrictionsForPurpose(p)
		for _, r := range pr {
			if (r.RestrictionType == RestrictionTypeNotAllowed || r.RestrictionType == RestrictionTypeRequireLI) && r.IsVendorIncluded(id) {
				return false
			}
		}
	}

	return true
}

// Returns true if transparency for vendor id's legitimate interest is established for all purpose ids
// and publisher hasn't set restrictions for them
func (c *CoreString) IsVendorAllowedForPurposesLI(id int, purposeIds ...int) bool {
	if !c.IsVendorLIAllowed(id) {
		return false
	}

	for _, p := range purposeIds {
		if !c.IsPurposeLIAllowed(p) {
			return false
		}
	}

	for _, p := range purposeIds {
		pr := c.GetPubRestrictionsForPurpose(p)
		for _, r := range pr {
			if (r.RestrictionType == RestrictionTypeNotAllowed || r.RestrictionType == RestrictionTypeRequireConsent) && r.IsVendorIncluded(id) {
				return false
			}
		}
	}

	return true
}

// Returns true if user has given consent to vendor id processing all purposes ids
// or if transparency for its legitimate interest is established in accordance with publisher restrictions
func (c *CoreString) IsVendorAllowedForFlexiblePurposes(id int, purposeIds ...int) bool {
	if !c.IsVendorAllowed(id) && !c.IsVendorLIAllowed(id) {
		return false
	}

	for _, p := range purposeIds {
		if !c.IsPurposeAllowed(p) && !c.IsPurposeLIAllowed(p) {
			return false
		}

		pr := c.GetPubRestrictionsForPurpose(p)
		if len(pr) == 0 && (!c.IsVendorAllowed(id) || !c.IsPurposeAllowed(p)) {
			return false
		}

		for _, r := range pr {
			if !r.IsVendorIncluded(id) {
				continue
			}
			switch r.RestrictionType {
			case RestrictionTypeNotAllowed:
				return false
			case RestrictionTypeRequireConsent:
				if !c.IsVendorAllowed(id) || !c.IsPurposeAllowed(p) {
					return false
				}
			case RestrictionTypeRequireLI:
				if !c.IsVendorLIAllowed(id) || !c.IsPurposeLIAllowed(p) {
					return false
				}
			}
		}
	}

	return true
}

// Returns true if transparency for vendor id's legitimate interest is established for all purpose ids
// or if user has given consent in accordance with publisher restrictions
func (c *CoreString) IsVendorAllowedForFlexiblePurposesLI(id int, purposeIds ...int) bool {
	if !c.IsVendorAllowed(id) && !c.IsVendorLIAllowed(id) {
		return false
	}

	for _, p := range purposeIds {
		if !c.IsPurposeAllowed(p) && !c.IsPurposeLIAllowed(p) {
			return false
		}

		pr := c.GetPubRestrictionsForPurpose(p)
		if len(pr) == 0 && (!c.IsVendorLIAllowed(id) || !c.IsPurposeLIAllowed(p)) {
			return false
		}

		for _, r := range pr {
			if !r.IsVendorIncluded(id) {
				continue
			}
			switch r.RestrictionType {
			case RestrictionTypeNotAllowed:
				return false
			case RestrictionTypeRequireConsent:
				if !c.IsVendorAllowed(id) || !c.IsPurposeAllowed(p) {
					return false
				}
			case RestrictionTypeRequireLI:
				if !c.IsVendorLIAllowed(id) || !c.IsPurposeLIAllowed(p) {
					return false
				}
			}
		}
	}

	return true
}

// Returns a list of publisher restrictions applied to purpose id
func (c *CoreString) GetPubRestrictionsForPurpose(id int) []*PubRestriction {
	var pr []*PubRestriction
	for _, r := range c.PubRestrictions {
		if r.PurposeId == id {
			pr = append(pr, r)
		}
	}
	return pr
}

// Returns true if restriction is applied to vendor id
func (p *PubRestriction) IsVendorIncluded(id int) bool {
	for _, entry := range p.RangeEntries {
		if entry.StartVendorID <= id && id <= entry.EndVendorID {
			return true
		}
	}
	return false
}

// Returns structure as a base64 raw url encoded string
func (c *CoreString) Encode() string {
	var bitSize int
	bitSize += bitsVersion
	bitSize += bitsCreated
	bitSize += bitsLastUpdated
	bitSize += bitsCmpId
	bitSize += bitsCmpVersion
	bitSize += bitsConsentScreen
	bitSize += bitsConsentLanguage
	bitSize += bitsVendorListVersion
	bitSize += bitsTcfPolicyVersion
	bitSize += bitsIsServiceSpecific
	bitSize += bitsUseNonStandardTexts
	bitSize += bitsSpecialFeatureOptIns
	bitSize += bitsPurposesConsent
	bitSize += bitsPurposesLITransparency
	bitSize += bitsPurposeOneTreatment
	bitSize += bitsPublisherCC

	bitSize += bitsMaxVendorId
	bitSize += bitsIsRangeEncoding
	if c.IsRangeEncoding {
		bitSize += bitsNumEntries
		for _, entry := range c.RangeEntries {
			bitSize += entry.getBitSize()
		}
	} else {
		if c.MaxVendorId == 0 {
			for id, _ := range c.VendorsConsent {
				if id > c.MaxVendorId {
					c.MaxVendorId = id
				}
			}
		}
		bitSize += c.MaxVendorId
	}

	bitSize += bitsMaxVendorId
	bitSize += bitsIsRangeEncoding
	if c.IsRangeEncodingLI {
		bitSize += bitsNumEntries
		for _, entry := range c.RangeEntriesLI {
			bitSize += entry.getBitSize()
		}
	} else {
		if c.MaxVendorIdLI == 0 {
			for id, _ := range c.VendorsLITransparency {
				if id > c.MaxVendorIdLI {
					c.MaxVendorIdLI = id
				}
			}
		}
		bitSize += c.MaxVendorIdLI
	}

	bitSize += bitsNumPubRestrictions
	for _, restriction := range c.PubRestrictions {
		bitSize += restriction.getBitSize()
	}

	e := NewTCEncoderFromSize(bitSize)
	e.WriteInt(c.Version, bitsVersion)
	e.WriteTime(c.Created)
	e.WriteTime(c.LastUpdated)
	e.WriteInt(c.CmpId, bitsCmpId)
	e.WriteInt(c.CmpVersion, bitsCmpVersion)
	e.WriteInt(c.ConsentScreen, bitsConsentScreen)
	e.WriteChars(c.ConsentLanguage, bitsConsentLanguage)
	e.WriteInt(c.VendorListVersion, bitsVendorListVersion)
	e.WriteInt(c.TcfPolicyVersion, bitsTcfPolicyVersion)
	e.WriteBool(c.IsServiceSpecific)
	e.WriteBool(c.UseNonStandardTexts)
	e.WriteBools(c.IsSpecialFeatureAllowed, bitsSpecialFeatureOptIns)
	e.WriteBools(c.IsPurposeAllowed, bitsPurposesConsent)
	e.WriteBools(c.IsPurposeLIAllowed, bitsPurposesLITransparency)
	e.WriteBool(c.PurposeOneTreatment)
	e.WriteChars(c.PublisherCC, bitsPublisherCC)

	e.WriteInt(c.MaxVendorId, bitsMaxVendorId)
	e.WriteBool(c.IsRangeEncoding)
	if c.IsRangeEncoding {
		e.WriteRangeEntries(c.RangeEntries)
	} else {
		e.WriteBools(c.IsVendorAllowed, c.MaxVendorId)
	}

	e.WriteInt(c.MaxVendorIdLI, bitsMaxVendorId)
	e.WriteBool(c.IsRangeEncodingLI)
	if c.IsRangeEncodingLI {
		e.WriteRangeEntries(c.RangeEntriesLI)
	} else {
		e.WriteBools(c.IsVendorLIAllowed, c.MaxVendorIdLI)
	}

	e.WritePubRestrictions(c.PubRestrictions)

	return base64.RawURLEncoding.EncodeToString(e.Bytes)
}
