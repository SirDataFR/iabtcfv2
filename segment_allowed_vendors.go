package iabtcfv2

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

// Returns true if vendor id is allowed for OOB signaling
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

// Returns structure as a base64 raw url encoded string
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
		if a.MaxVendorId == 0 {
			for id, _ := range a.AllowedVendors {
				if id > a.MaxVendorId {
					a.MaxVendorId = id
				}
			}
		}
		bitSize += a.MaxVendorId
	}

	var e = newTCEncoder(make([]byte, bitSize/8))
	if bitSize%8 != 0 {
		e = newTCEncoder(make([]byte, bitSize/8+1))
	}

	e.writeInt(a.SegmentType, 3)
	e.writeInt(a.MaxVendorId, 16)
	e.writeBool(a.IsRangeEncoding)
	if a.IsRangeEncoding {
		e.writeInt(len(a.RangeEntries), 12)
		e.writeRangeEntries(a.RangeEntries)
	} else {
		for i := 0; i < a.MaxVendorId; i++ {
			e.writeBool(a.IsVendorAllowed(i + 1))
		}
	}

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
