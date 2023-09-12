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

func newTCEncoderFromSize(bitSize int) *TCEncoder {
	if bitSize%8 != 0 {
		return newTCEncoder(make([]byte, bitSize/8+1))
	}
	return newTCEncoder(make([]byte, bitSize/8))
}

func (r *TCEncoder) readTime() time.Time {
	var ds = int64(r.readInt(bitsTime))
	return time.Unix(ds/decisecondsPerSecond, (ds%decisecondsPerSecond)*nanosecondsPerDecisecond).UTC()
}

func (r *TCEncoder) writeTime(v time.Time) {
	r.writeNumber(v.UnixNano()/nanosecondsPerDecisecond, bitsTime)
}

func (r *TCEncoder) readChars(n uint) string {
	var buf = make([]byte, 0, n/bitsChar)
	for i := uint(0); i < n/bitsChar; i++ {
		buf = append(buf, byte(r.readInt(bitsChar))+'A')
	}
	return string(buf)
}

func (r *TCEncoder) writeChars(v string, n uint) {
	for i := uint(0); i < n/bitsChar; i++ {
		char := v[i]
		r.writeInt(int(byte(char)-'A'), bitsChar)
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

func (b *Bits) writeBools(getBool func(int) bool, n int) {
	for i := 1; i <= n; i++ {
		b.writeBool(getBool(i))
	}
}

func (r *TCEncoder) writeRangeEntries(entries []*RangeEntry) {
	r.writeInt(len(entries), bitsNumEntries)
	for _, entry := range entries {
		if entry.EndVendorID > entry.StartVendorID {
			r.writeBool(true)
			r.writeInt(entry.StartVendorID, bitsVendorId)
			r.writeInt(entry.EndVendorID, bitsVendorId)
		} else {
			r.writeBool(false)
			r.writeInt(entry.StartVendorID, bitsVendorId)
		}
	}
}

func (r *TCEncoder) readRangeEntries() (int, []*RangeEntry) {
	n := r.readInt(bitsNumEntries)
	var ret = make([]*RangeEntry, 0, n)
	for i := uint(0); i < uint(n); i++ {
		var isRange = r.readBool()
		var start, end int
		start = r.readInt(bitsVendorId)
		if isRange {
			end = r.readInt(bitsVendorId)
		} else {
			end = start
		}
		ret = append(ret, &RangeEntry{StartVendorID: start, EndVendorID: end})
	}
	return n, ret
}

func (r *TCEncoder) writePubRestrictions(entries []*PubRestriction) {
	r.writeInt(len(entries), bitsNumPubRestrictions)
	for _, entry := range entries {
		r.writeInt(entry.PurposeId, bitsPubRestrictionsEntryPurposeId)
		r.writeInt(int(entry.RestrictionType), bitsPubRestrictionsEntryRestrictionType)
		r.writeRangeEntries(entry.RangeEntries)
	}
}

func (r *TCEncoder) readPubRestrictions() (int, []*PubRestriction) {
	n := r.readInt(bitsNumPubRestrictions)
	var ret = make([]*PubRestriction, 0, n)
	for i := uint(0); i < uint(n); i++ {
		var purposeId = r.readInt(bitsPubRestrictionsEntryPurposeId)
		var restrictionType = r.readInt(bitsPubRestrictionsEntryRestrictionType)
		numEntries, rangeEntries := r.readRangeEntries()
		ret = append(ret, &PubRestriction{PurposeId: purposeId,
			RestrictionType: RestrictionType(restrictionType),
			NumEntries:      numEntries,
			RangeEntries:    rangeEntries,
		})
	}
	return n, ret
}
