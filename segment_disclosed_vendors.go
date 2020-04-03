package iabtcf

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

func (d *DisclosedVendors) Encode() string {
	bitSize := 20

	if d.IsRangeEncoding {
		bitSize += 12
		entriesSize := len(d.RangeEntries)
		for _, entry := range d.RangeEntries {
			if entry.EndVendorID > entry.StartVendorID {
				entriesSize += 16 * 2
			} else {
				entriesSize += 16
			}
		}
		bitSize += entriesSize
	} else {
		bitSize += d.MaxVendorId
	}

	var e = NewTCEncoder(make([]byte, bitSize/8))
	if bitSize%8 != 0 {
		e = NewTCEncoder(make([]byte, bitSize/8+1))
	}

	e.WriteInt(d.SegmentType, 3)
	e.WriteInt(d.MaxVendorId, 16)
	e.WriteBool(d.IsRangeEncoding)
	if d.IsRangeEncoding {
		e.WriteInt(len(d.RangeEntries), 12)
		e.WriteRangeEntries(d.RangeEntries)
	} else {
		for i := 0; i < d.MaxVendorId; i++ {
			e.WriteBool(d.IsVendorDisclosed(i + 1))
		}
	}

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
