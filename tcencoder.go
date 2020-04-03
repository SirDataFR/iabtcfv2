package iabtcfv2

import (
	"time"
)

const (
	decisecondsPerSecond     = 10
	nanosecondsPerDecisecond = int64(time.Millisecond * 100)
)

type TCEncoder struct {
	*Bits
}

func newTCEncoder(src []byte) *TCEncoder {
	return &TCEncoder{newBits(src)}
}

func (r *TCEncoder) readTime() time.Time {
	var ds = int64(r.readInt(36))
	return time.Unix(ds/decisecondsPerSecond, (ds%decisecondsPerSecond)*nanosecondsPerDecisecond).UTC()
}

func (r *TCEncoder) writeTime(v time.Time) {
	r.writeNumber(v.UnixNano()/nanosecondsPerDecisecond, 36)
}

func (r *TCEncoder) readIsoCode() string {
	var buf = make([]byte, 0, 2)
	for i := uint(0); i < 2; i++ {
		buf = append(buf, byte(r.readInt(6))+'A')
	}
	return string(buf)
}

func (r *TCEncoder) writeIsoCode(v string) {
	for _, char := range v {
		r.writeInt(int(byte(char)-'A'), 6)
	}
}

func (r *TCEncoder) readBitField(n uint) map[int]bool {
	var m = make(map[int]bool)
	for i := uint(0); i < n; i++ {
		if r.readBool() {
			m[int(i)+1] = true
		}
	}
	return m
}

func (r *TCEncoder) writeRangeEntries(entries []*RangeEntry) {
	for _, entry := range entries {
		if entry.EndVendorID > entry.StartVendorID {
			r.writeBool(true)
			r.writeInt(entry.StartVendorID, 16)
			r.writeInt(entry.EndVendorID, 16)
		} else {
			r.writeBool(false)
			r.writeInt(entry.StartVendorID, 16)
		}
	}
}

func (r *TCEncoder) readRangeEntries(n uint) []*RangeEntry {
	var ret = make([]*RangeEntry, 0, n)
	for i := uint(0); i < n; i++ {
		var isRange = r.readBool()
		var start, end int
		start = r.readInt(16)
		if isRange {
			end = r.readInt(16)
		} else {
			end = start
		}
		ret = append(ret, &RangeEntry{StartVendorID: start, EndVendorID: end})
	}
	return ret
}

func (r *TCEncoder) writePubRestrictions(entries []*PubRestriction) {
	for _, entry := range entries {
		r.writeInt(entry.PurposeId, 6)
		r.writeInt(entry.RestrictionType, 2)
		r.writeInt(len(entry.RangeEntries), 12)
		r.writeRangeEntries(entry.RangeEntries)
	}
}

func (r *TCEncoder) readPubRestrictions(n uint) []*PubRestriction {
	var ret = make([]*PubRestriction, 0, n)
	for i := uint(0); i < n; i++ {
		var purposeId = r.readInt(6)
		var restrictionType = r.readInt(2)
		var numEntries = r.readInt(12)
		var rangeEntries = r.readRangeEntries(uint(numEntries))
		ret = append(ret, &PubRestriction{PurposeId: purposeId,
			RestrictionType: restrictionType,
			NumEntries:      numEntries,
			RangeEntries:    rangeEntries,
		})
	}
	return ret
}
