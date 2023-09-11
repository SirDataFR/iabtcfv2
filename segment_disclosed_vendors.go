package iabtcfv2

import (
	"encoding/base64"
)

type DisclosedVendors struct {
	SegmentType      int
	MaxVendorId      int
	IsRangeEncoding  bool
	DisclosedVendors map[int]bool
	NumEntries       int
	RangeEntries     []*RangeEntry
}

// Returns true if vendor id is disclosed for validating OOB signaling
func (d *DisclosedVendors) IsVendorDisclosed(id int) bool {
	if d.IsRangeEncoding {
		for _, entry := range d.RangeEntries {
			if entry.StartVendorID <= id && id <= entry.EndVendorID {
				return true
			}
		}
		return false
	}

	return d.DisclosedVendors[id]
}

// Returns structure as a base64 raw url encoded string
func (d *DisclosedVendors) Encode() string {
	var bitSize int
	bitSize += bitsSegmentType

	bitSize += bitsMaxVendorId
	bitSize += bitsIsRangeEncoding
	if d.IsRangeEncoding {
		bitSize += bitsNumEntries
		for _, entry := range d.RangeEntries {
			bitSize += entry.getBitSize()
		}
	} else {
		if d.MaxVendorId == 0 {
			for id, _ := range d.DisclosedVendors {
				if id > d.MaxVendorId {
					d.MaxVendorId = id
				}
			}
		}
		bitSize += d.MaxVendorId
	}

	e := newTCEncoderFromSize(bitSize)
	e.writeInt(d.SegmentType, bitsSegmentType)
	e.writeInt(d.MaxVendorId, bitsMaxVendorId)
	e.writeBool(d.IsRangeEncoding)
	if d.IsRangeEncoding {
		e.writeRangeEntries(d.RangeEntries)
	} else {
		e.writeBools(d.IsVendorDisclosed, d.MaxVendorId)
	}

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
