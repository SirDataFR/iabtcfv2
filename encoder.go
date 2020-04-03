package iabtcf

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

func NewTCEncoder(src []byte) *TCEncoder {
	return &TCEncoder{NewBits(src)}
}

func (r *TCEncoder) ReadTime() time.Time {
	var ds = int64(r.ReadInt(36))
	return time.Unix(ds/decisecondsPerSecond, (ds%decisecondsPerSecond)*nanosecondsPerDecisecond).UTC()
}

func (r *TCEncoder) WriteTime(v time.Time) {
	r.WriteNumber(v.UnixNano()/nanosecondsPerDecisecond, 36)
}

func (r *TCEncoder) ReadIsoCode() string {
	var buf = make([]byte, 0, 2)
	for i := uint(0); i < 2; i++ {
		buf = append(buf, byte(r.ReadInt(6))+'A')
	}
	return string(buf)
}

func (r *TCEncoder) WriteIsoCode(v string) {
	for _, char := range v {
		r.WriteInt(int(byte(char)-'A'), 6)
	}
}

func (r *TCEncoder) ReadBitField(n uint) map[int]bool {
	var m = make(map[int]bool)
	for i := uint(0); i < n; i++ {
		if r.ReadBool() {
			m[int(i)+1] = true
		}
	}
	return m
}

func (r *TCEncoder) WriteRangeEntries(entries []*RangeEntry) {
	for _, entry := range entries {
		if entry.EndVendorID > entry.StartVendorID {
			r.WriteBool(true)
			r.WriteInt(entry.StartVendorID, 16)
			r.WriteInt(entry.EndVendorID, 16)
		} else {
			r.WriteBool(false)
			r.WriteInt(entry.StartVendorID, 16)
		}
	}
}

func (r *TCEncoder) ReadRangeEntries(n uint) []*RangeEntry {
	var ret = make([]*RangeEntry, 0, n)
	for i := uint(0); i < n; i++ {
		var isRange = r.ReadBool()
		var start, end int
		start = r.ReadInt(16)
		if isRange {
			end = r.ReadInt(16)
		} else {
			end = start
		}
		ret = append(ret, &RangeEntry{StartVendorID: start, EndVendorID: end})
	}
	return ret
}

func (r *TCEncoder) WritePubRestrictions(entries []*PubRestriction) {
	for _, entry := range entries {
		r.WriteInt(entry.PurposeId, 6)
		r.WriteInt(entry.RestrictionType, 2)
		r.WriteInt(len(entry.RangeEntries), 12)
		r.WriteRangeEntries(entry.RangeEntries)
	}
}

func (r *TCEncoder) ReadPubRestrictions(n uint) []*PubRestriction {
	var ret = make([]*PubRestriction, 0, n)
	for i := uint(0); i < n; i++ {
		var purposeId = r.ReadInt(6)
		var restrictionType = r.ReadInt(2)
		var numEntries = r.ReadInt(12)
		var rangeEntries = r.ReadRangeEntries(uint(numEntries))
		ret = append(ret, &PubRestriction{PurposeId: purposeId,
			RestrictionType: restrictionType,
			NumEntries:      numEntries,
			RangeEntries:    rangeEntries,
		})
	}
	return ret
}
