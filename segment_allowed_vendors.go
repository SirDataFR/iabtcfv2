package iabtcf

import (
	"encoding/base64"
)

type AllowedVendors struct {
	SegmentType     int
	MaxVendorId     int
	IsRangeEncoding bool
	AllowedVendors  map[int]bool
	NumEntries      int
	RangeEntries    []*RangeEntry
}

func (a *AllowedVendors) IsVendorAllowed(id int) bool {
	if a.IsRangeEncoding {
		for _, entry := range a.RangeEntries {
			if entry.StartVendorID <= id && id <= entry.EndVendorID {
				return true
			}
		}
		return false
	}

	return a.AllowedVendors[id]
}

func (a *AllowedVendors) Encode() string {
	bitSize := 20

	if a.IsRangeEncoding {
		bitSize += 12
		entriesSize := len(a.RangeEntries)
		for _, entry := range a.RangeEntries {
			if entry.EndVendorID > entry.StartVendorID {
				entriesSize += 16 * 2
			} else {
				entriesSize += 16
			}
		}
		bitSize += entriesSize
	} else {
		bitSize += a.MaxVendorId
	}

	var e = NewTCEncoder(make([]byte, bitSize/8))
	if bitSize%8 != 0 {
		e = NewTCEncoder(make([]byte, bitSize/8+1))
	}

	e.WriteInt(a.SegmentType, 3)
	e.WriteInt(a.MaxVendorId, 16)
	e.WriteBool(a.IsRangeEncoding)
	if a.IsRangeEncoding {
		e.WriteInt(len(a.RangeEntries), 12)
		e.WriteRangeEntries(a.RangeEntries)
	} else {
		for i := 0; i < a.MaxVendorId; i++ {
			e.WriteBool(a.IsVendorAllowed(i + 1))
		}
	}

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
