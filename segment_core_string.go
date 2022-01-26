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
	UseNonStandardStacks   bool
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

type RangeEntry struct {
	StartVendorID int
	EndVendorID   int
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
	bitSize := 230

	if c.IsRangeEncoding {
		bitSize += 12
		entriesSize := len(c.RangeEntries)
		for _, entry := range c.RangeEntries {
			if entry.EndVendorID > entry.StartVendorID {
				entriesSize += 16 * 2
			} else {
				entriesSize += 16
			}
		}
		bitSize += +entriesSize
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

	bitSize += 16
	if c.IsRangeEncodingLI {
		bitSize += 12
		entriesSize := len(c.RangeEntriesLI)
		for _, entry := range c.RangeEntriesLI {
			if entry.EndVendorID > entry.StartVendorID {
				entriesSize += 16 * 2
			} else {
				entriesSize += 16
			}
		}
		bitSize += entriesSize
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

	bitSize += 12
	for _, res := range c.PubRestrictions {
		entriesSize := 20
		for _, entry := range res.RangeEntries {
			entriesSize++
			if entry.EndVendorID > entry.StartVendorID {
				entriesSize += 16 * 2
			} else {
				entriesSize += 16
			}
		}
		bitSize += entriesSize
	}

	var e = newTCEncoder(make([]byte, bitSize/8))
	if bitSize%8 != 0 {
		e = newTCEncoder(make([]byte, bitSize/8+1))
	}

	e.writeInt(c.Version, 6)
	e.writeTime(c.Created)
	e.writeTime(c.LastUpdated)
	e.writeInt(c.CmpId, 12)
	e.writeInt(c.CmpVersion, 12)
	e.writeInt(c.ConsentScreen, 6)
	e.writeIsoCode(c.ConsentLanguage)
	e.writeInt(c.VendorListVersion, 12)
	e.writeInt(c.TcfPolicyVersion, 6)
	e.writeBool(c.IsServiceSpecific)
	e.writeBool(c.UseNonStandardStacks)
	for i := 0; i < 12; i++ {
		e.writeBool(c.IsSpecialFeatureAllowed(i + 1))
	}
	for i := 0; i < 24; i++ {
		e.writeBool(c.IsPurposeAllowed(i + 1))
	}
	for i := 0; i < 24; i++ {
		e.writeBool(c.IsPurposeLIAllowed(i + 1))
	}
	e.writeBool(c.PurposeOneTreatment)
	e.writeIsoCode(c.PublisherCC)

	e.writeInt(c.MaxVendorId, 16)
	e.writeBool(c.IsRangeEncoding)
	if c.IsRangeEncoding {
		e.writeInt(len(c.RangeEntries), 12)
		e.writeRangeEntries(c.RangeEntries)
	} else {
		for i := 0; i < c.MaxVendorId; i++ {
			e.writeBool(c.IsVendorAllowed(i + 1))
		}
	}

	e.writeInt(c.MaxVendorIdLI, 16)
	e.writeBool(c.IsRangeEncodingLI)
	if c.IsRangeEncodingLI {
		e.writeInt(len(c.RangeEntriesLI), 12)
		e.writeRangeEntries(c.RangeEntriesLI)
	} else {
		for i := 0; i < c.MaxVendorIdLI; i++ {
			e.writeBool(c.IsVendorLIAllowed(i + 1))
		}
	}

	e.writeInt(len(c.PubRestrictions), 12)
	e.writePubRestrictions(c.PubRestrictions)

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
