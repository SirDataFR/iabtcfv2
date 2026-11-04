package iabtcfv2

import (
	"encoding/base64"
	"testing"
	"time"
)

// Helpers shared by the tests of flexible purposes and of decoding validation.

var testCreated = time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)

// testCore returns a valid core string, created after the TCF v2.3 deadline, that mutate completes.
func testCore(mutate func(c *CoreString)) *CoreString {
	c := &CoreString{
		Version: 2, Created: testCreated, LastUpdated: testCreated, CmpId: 92, CmpVersion: 1, ConsentScreen: 1,
		ConsentLanguage: "FR", VendorListVersion: 150, TcfPolicyVersion: 5, PublisherCC: "FR",
		SpecialFeatureOptIns: map[int]bool{}, PurposesConsent: map[int]bool{}, PurposesLITransparency: map[int]bool{},
		VendorsConsent: map[int]bool{}, VendorsLITransparency: map[int]bool{},
	}
	if mutate != nil {
		mutate(c)
	}
	return c
}

// testDisclosedVendors returns the Disclosed Vendors segment, mandatory since the TCF v2.3 deadline.
func testDisclosedVendors(ids ...int) string {
	m := map[int]bool{}
	for _, id := range ids {
		m[id] = true
	}
	return (&DisclosedVendors{SegmentType: int(SegmentTypeDisclosedVendors), DisclosedVendors: m}).Encode()
}

func testRestriction(purposeId int, restrictionType RestrictionType, vendorIds ...int) *PubRestriction {
	var entries []*RangeEntry
	for _, id := range vendorIds {
		entries = append(entries, &RangeEntry{StartVendorID: id, EndVendorID: id})
	}
	return &PubRestriction{PurposeId: purposeId, RestrictionType: restrictionType, NumEntries: len(entries), RangeEntries: entries}
}

func testDecode(t *testing.T, core *CoreString, disclosed ...int) *TCData {
	t.Helper()
	d, err := Decode(core.Encode() + "." + testDisclosedVendors(disclosed...))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// testWithRawRestriction replaces the publisher restrictions of an encoded core string by one restriction written
// bit by bit, so that the test can write values the encoder would not: an inverted range, a reserved type.
func testWithRawRestriction(core string, purposeId, restrictionType, start, end int) string {
	b, _ := base64.RawURLEncoding.DecodeString(core)
	e := NewTCEncoder(b)
	e.ReadInt(bitsVersion)
	e.ReadTime()
	e.ReadTime()
	e.ReadInt(bitsCmpId + bitsCmpVersion + bitsConsentScreen)
	e.ReadChars(bitsConsentLanguage)
	e.ReadInt(bitsVendorListVersion + bitsTcfPolicyVersion + 2)
	e.ReadBitField(bitsSpecialFeatureOptIns + bitsPurposesConsent + bitsPurposesLITransparency + 1)
	e.ReadChars(bitsPublisherCC)
	for i := 0; i < 2; i++ {
		if maxVendorId := e.ReadInt(bitsMaxVendorId); e.ReadBool() {
			e.ReadRangeEntries()
		} else {
			e.ReadBitField(uint(maxVendorId))
		}
	}
	position := e.Position
	src := NewTCEncoder(b)
	dst := NewTCEncoderFromSize(int(position) + bitsNumPubRestrictions + bitsPubRestrictionsEntryPurposeId +
		bitsPubRestrictionsEntryRestrictionType + bitsNumEntries + 1 + 2*bitsVendorId)
	for i := uint(0); i < position; i++ {
		dst.WriteBool(src.ReadBool())
	}
	dst.WriteInt(1, bitsNumPubRestrictions)
	dst.WriteInt(purposeId, bitsPubRestrictionsEntryPurposeId)
	dst.WriteInt(restrictionType, bitsPubRestrictionsEntryRestrictionType)
	dst.WriteInt(1, bitsNumEntries)
	dst.WriteBool(true)
	dst.WriteInt(start, bitsVendorId)
	dst.WriteInt(end, bitsVendorId)
	return base64.RawURLEncoding.EncodeToString(dst.Bytes)
}
