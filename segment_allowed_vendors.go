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
	var bitSize int
	bitSize += bitsSegmentType

	bitSize += bitsMaxVendorId
	bitSize += bitsIsRangeEncoding
	if a.IsRangeEncoding {
		bitSize += bitsNumEntries
		for _, entry := range a.RangeEntries {
			bitSize += entry.getBitSize()
		}
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

	e := newTCEncoderFromSize(bitSize)
	e.writeInt(a.SegmentType, bitsSegmentType)
	e.writeInt(a.MaxVendorId, bitsMaxVendorId)
	e.writeBool(a.IsRangeEncoding)
	if a.IsRangeEncoding {
		e.writeRangeEntries(a.RangeEntries)
	} else {
		e.writeBools(a.IsVendorAllowed, a.MaxVendorId)
	}

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
