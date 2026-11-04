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
func (c *CoreString) IsPurposeLIAllowed(id int) bool {
	// TCF v2.3: Purpose 1 cannot be established via legitimate interest
	if id == 1 && c.TcfPolicyVersion >= TcfPolicyVersion23 {
		return false
	}
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

// Returns true if vendor id is allowed for all purpose ids, each one on its legal basis: consent by default, or the
// basis that a publisher restriction which includes the vendor sets for the purpose. It is meant for the purposes
// that the vendor declares flexible in the Global Vendor List, with consent as default legal basis. Legitimate
// interest never allows purpose 1, nor purposes 3 to 6 from TCF v2.2.
func (c *CoreString) IsVendorAllowedForFlexiblePurposes(id int, purposeIds ...int) bool {
	return c.isVendorAllowedForFlexiblePurposes(id, RestrictionTypeRequireConsent, purposeIds)
}

// Returns true if vendor id is allowed for all purpose ids, each one on its legal basis: legitimate interest by
// default, or the basis that a publisher restriction which includes the vendor sets for the purpose. It is meant for
// the purposes that the vendor declares flexible in the Global Vendor List, with legitimate interest as default legal
// basis. Legitimate interest never allows purpose 1, nor purposes 3 to 6 from TCF v2.2.
func (c *CoreString) IsVendorAllowedForFlexiblePurposesLI(id int, purposeIds ...int) bool {
	return c.isVendorAllowedForFlexiblePurposes(id, RestrictionTypeRequireLI, purposeIds)
}

// isVendorAllowedForFlexiblePurposes checks each purpose on defaultBasis, unless publisher restrictions on the
// purpose include the vendor: each of them then sets the basis (not allowed, consent or legitimate interest). A
// restriction that includes other vendors only changes nothing for this one.
func (c *CoreString) isVendorAllowedForFlexiblePurposes(id int, defaultBasis RestrictionType, purposeIds []int) bool {
	if !c.IsVendorAllowed(id) && !c.IsVendorLIAllowed(id) {
		return false
	}
	for _, p := range purposeIds {
		restricted := false
		for _, r := range c.GetPubRestrictionsForPurpose(p) {
			if !r.IsVendorIncluded(id) {
				continue
			}
			restricted = true
			if !c.isAllowedOnBasis(id, p, r.RestrictionType) {
				return false
			}
		}
		if !restricted && !c.isAllowedOnBasis(id, p, defaultBasis) {
			return false
		}
	}
	return true
}

// isAllowedOnBasis returns true if vendor id and purpose purposeId are both allowed on the legal basis that basis
// names. A restriction type that allows nothing, or that the specification does not define, allows nothing.
func (c *CoreString) isAllowedOnBasis(id int, purposeId int, basis RestrictionType) bool {
	switch basis {
	case RestrictionTypeRequireConsent:
		return c.IsVendorAllowed(id) && c.IsPurposeAllowed(purposeId)
	case RestrictionTypeRequireLI:
		return c.canBeLegitimateInterest(purposeId) && c.IsVendorLIAllowed(id) && c.IsPurposeLIAllowed(purposeId)
	default:
		return false
	}
}

// canBeLegitimateInterest returns false for the purposes that the specification never allows on legitimate
// interest, whatever the bits of the string: purpose 1, always a consent purpose and never flexible, and purposes 3
// to 6 from TCF v2.2 (policy version 4), whose legitimate interest bits are required to be 0. A restriction that
// requires legitimate interest for them means that the vendor should not process.
func (c *CoreString) canBeLegitimateInterest(purposeId int) bool {
	if purposeId == 1 {
		return false
	}
	return c.TcfPolicyVersion < tcfPolicyVersion22 || purposeId < 3 || purposeId > 6
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
			// MaxVendorId is the largest vendor id that the section represents: a range past it does not decode.
			if entry.EndVendorID > c.MaxVendorId {
				c.MaxVendorId = entry.EndVendorID
			}
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
			if entry.EndVendorID > c.MaxVendorIdLI {
				c.MaxVendorIdLI = entry.EndVendorID
			}
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
