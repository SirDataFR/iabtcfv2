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

func NewTCEncoder(src []byte) *TCEncoder {
	return &TCEncoder{NewBits(src)}
}

func NewTCEncoderFromSize(bitSize int) *TCEncoder {
	if bitSize%8 != 0 {
		return NewTCEncoder(make([]byte, bitSize/8+1))
	}
	return NewTCEncoder(make([]byte, bitSize/8))
}

func (r *TCEncoder) ReadTime() time.Time {
	var ds = int64(r.ReadInt(bitsTime))
	return time.Unix(ds/decisecondsPerSecond, (ds%decisecondsPerSecond)*nanosecondsPerDecisecond).UTC()
}

func (r *TCEncoder) WriteTime(v time.Time) {
	r.WriteNumber(v.UnixNano()/nanosecondsPerDecisecond, bitsTime)
}

func (r *TCEncoder) ReadChars(n uint) string {
	var buf = make([]byte, 0, n/bitsChar)
	for i := uint(0); i < n/bitsChar; i++ {
		buf = append(buf, byte(r.ReadInt(bitsChar))+'A')
	}
	return string(buf)
}

func (r *TCEncoder) WriteChars(v string, n uint) {
	for i := uint(0); i < n/bitsChar; i++ {
		char := v[i]
		r.WriteInt(int(byte(char)-'A'), bitsChar)
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

func (b *Bits) WriteBools(getBool func(int) bool, n int) {
	for i := 1; i <= n; i++ {
		b.WriteBool(getBool(i))
	}
}

func (r *TCEncoder) WriteRangeEntries(entries []*RangeEntry) {
	r.WriteInt(len(entries), bitsNumEntries)
	for _, entry := range entries {
		if entry.EndVendorID > entry.StartVendorID {
			r.WriteBool(true)
			r.WriteInt(entry.StartVendorID, bitsVendorId)
			r.WriteInt(entry.EndVendorID, bitsVendorId)
		} else {
			r.WriteBool(false)
			r.WriteInt(entry.StartVendorID, bitsVendorId)
		}
	}
}

func (r *TCEncoder) ReadRangeEntries() (int, []*RangeEntry) {
	n := r.ReadInt(bitsNumEntries)
	var ret = make([]*RangeEntry, 0, n)
	for i := uint(0); i < uint(n); i++ {
		var isRange = r.ReadBool()
		var start, end int
		start = r.ReadInt(bitsVendorId)
		if isRange {
			end = r.ReadInt(bitsVendorId)
		} else {
			end = start
		}
		ret = append(ret, &RangeEntry{StartVendorID: start, EndVendorID: end})
	}
	return n, ret
}

func (r *TCEncoder) WritePubRestrictions(entries []*PubRestriction) {
	r.WriteInt(len(entries), bitsNumPubRestrictions)
	for _, entry := range entries {
		r.WriteInt(entry.PurposeId, bitsPubRestrictionsEntryPurposeId)
		r.WriteInt(int(entry.RestrictionType), bitsPubRestrictionsEntryRestrictionType)
		r.WriteRangeEntries(entry.RangeEntries)
	}
}

func (r *TCEncoder) ReadPubRestrictions() (int, []*PubRestriction) {
	n := r.ReadInt(bitsNumPubRestrictions)
	var ret = make([]*PubRestriction, 0, n)
	for i := uint(0); i < uint(n); i++ {
		var purposeId = r.ReadInt(bitsPubRestrictionsEntryPurposeId)
		var restrictionType = r.ReadInt(bitsPubRestrictionsEntryRestrictionType)
		numEntries, rangeEntries := r.ReadRangeEntries()
		ret = append(ret, &PubRestriction{PurposeId: purposeId,
			RestrictionType: RestrictionType(restrictionType),
			NumEntries:      numEntries,
			RangeEntries:    rangeEntries,
		})
	}
	return n, ret
}
