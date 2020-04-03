package iabtcf

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
	RestrictionType int
	NumEntries      int
	RangeEntries    []*RangeEntry
}

type RangeEntry struct {
	StartVendorID int
	EndVendorID   int
}

func (c *CoreString) IsSpecialFeatureAllowed(id int) bool {
	return c.SpecialFeatureOptIns[id]
}

func (c *CoreString) IsPurposeAllowed(id int) bool {
	return c.PurposesConsent[id]
}

func (c *CoreString) IsPurposeLIAllowed(id int) bool {
	return c.PurposesLITransparency[id]
}

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
		bitSize += c.MaxVendorIdLI
	}

	bitSize += 12
	for _, res := range c.PubRestrictions {
		entriesSize := 20
		for _, entry := range res.RangeEntries {
			if entry.EndVendorID > entry.StartVendorID {
				entriesSize += 16 * 2
			} else {
				entriesSize += 16
			}
		}
		bitSize += entriesSize
	}

	var e = NewTCEncoder(make([]byte, bitSize/8))
	if bitSize%8 != 0 {
		e = NewTCEncoder(make([]byte, bitSize/8+1))
	}

	e.WriteInt(c.Version, 6)
	e.WriteTime(c.Created)
	e.WriteTime(c.LastUpdated)
	e.WriteInt(c.CmpId, 12)
	e.WriteInt(c.CmpVersion, 12)
	e.WriteInt(c.ConsentScreen, 6)
	e.WriteIsoCode(c.ConsentLanguage)
	e.WriteInt(c.VendorListVersion, 12)
	e.WriteInt(c.TcfPolicyVersion, 6)
	e.WriteBool(c.IsServiceSpecific)
	e.WriteBool(c.UseNonStandardStacks)
	for i := 0; i < 12; i++ {
		e.WriteBool(c.IsSpecialFeatureAllowed(i + 1))
	}
	for i := 0; i < 24; i++ {
		e.WriteBool(c.IsPurposeAllowed(i + 1))
	}
	for i := 0; i < 24; i++ {
		e.WriteBool(c.IsPurposeLIAllowed(i + 1))
	}
	e.WriteBool(c.PurposeOneTreatment)
	e.WriteIsoCode(c.PublisherCC)

	e.WriteInt(c.MaxVendorId, 16)
	e.WriteBool(c.IsRangeEncoding)
	if c.IsRangeEncoding {
		e.WriteInt(len(c.RangeEntries), 12)
		e.WriteRangeEntries(c.RangeEntries)
	} else {
		for i := 0; i < c.MaxVendorId; i++ {
			e.WriteBool(c.IsVendorAllowed(i + 1))
		}
	}

	e.WriteInt(c.MaxVendorIdLI, 16)
	e.WriteBool(c.IsRangeEncodingLI)
	if c.IsRangeEncodingLI {
		e.WriteInt(len(c.RangeEntriesLI), 12)
		e.WriteRangeEntries(c.RangeEntriesLI)
	} else {
		for i := 0; i < c.MaxVendorIdLI; i++ {
			e.WriteBool(c.IsVendorLIAllowed(i + 1))
		}
	}

	e.WriteInt(len(c.PubRestrictions), 12)
	e.WritePubRestrictions(c.PubRestrictions)

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
