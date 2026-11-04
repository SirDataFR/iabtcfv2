package iabtcfv2

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestDecode_RejectsOtherVersions(t *testing.T) {
	for _, version := range []int{1, 3} {
		core := testCore(func(c *CoreString) { c.Version = version; c.PurposesConsent[1] = true })
		if _, err := Decode(core.Encode() + "." + testDisclosedVendors(53)); err == nil || !strings.Contains(err.Error(), "version") {
			t.Errorf("version %d: expected an unsupported version error, got %v", version, err)
		}
	}
}

func TestDecode_RejectsUnknownSegmentTypes(t *testing.T) {
	core := testCore(func(c *CoreString) { c.PurposesConsent[1] = true }).Encode()
	unknown := (&DisclosedVendors{SegmentType: 7, DisclosedVendors: map[int]bool{53: true}}).Encode()
	for _, decode := range []func(string) (*TCData, error){Decode, DecodeLenient} {
		if _, err := decode(core + "." + testDisclosedVendors(53) + "." + unknown); err == nil || !strings.Contains(err.Error(), "segment type 7") {
			t.Errorf("expected an unsupported segment type error, got %v", err)
		}
	}
}

// An inverted range used to restrict nothing: the restriction was decoded, and silently ignored.
func TestDecode_RejectsInvalidRanges(t *testing.T) {
	core := testCore(func(c *CoreString) {
		c.PurposesConsent[1], c.PurposesConsent[4] = true, true
		c.VendorsConsent[53] = true
	}).Encode()
	valid := testWithRawRestriction(core, 4, int(RestrictionTypeNotAllowed), 40, 60) + "." + testDisclosedVendors(53)
	d, err := Decode(valid)
	if err != nil {
		t.Fatal(err)
	}
	if d.IsVendorAllowedForPurposes(53, 1, 4) {
		t.Fatal("vendor 53 allowed for purpose 4 despite the restriction")
	}
	for name, s := range map[string]string{
		"inverted range":             testWithRawRestriction(core, 4, int(RestrictionTypeNotAllowed), 60, 40),
		"range starting at vendor 0": testWithRawRestriction(core, 4, int(RestrictionTypeNotAllowed), 0, 60),
		"undefined restriction type": testWithRawRestriction(core, 4, int(RestrictionTypeUndefined), 53, 53),
	} {
		if _, err := Decode(s + "." + testDisclosedVendors(53)); err == nil {
			t.Errorf("%s: decoded without error", name)
		}
	}
	if _, err := Decode(testRangeBeyondMaxVendorId() + "." + testDisclosedVendors(53)); err == nil {
		t.Error("range of vendor consents beyond MaxVendorId decoded without error")
	}
}

// testRangeBeyondMaxVendorId returns a core string whose vendor consents hold the range 1-1000 under a MaxVendorId of
// 10. The encoder raises MaxVendorId to the end of the ranges, so the field is written over afterwards.
func testRangeBeyondMaxVendorId() string {
	core := testCore(func(c *CoreString) {
		c.PurposesConsent[1] = true
		c.IsRangeEncoding = true
		c.RangeEntries = []*RangeEntry{{StartVendorID: 1, EndVendorID: 1000}}
	}).Encode()
	b, _ := base64.RawURLEncoding.DecodeString(core)
	e := NewTCEncoder(b)
	e.Position = bitsVersion + bitsCreated + bitsLastUpdated + bitsCmpId + bitsCmpVersion + bitsConsentScreen +
		bitsConsentLanguage + bitsVendorListVersion + bitsTcfPolicyVersion + bitsIsServiceSpecific +
		bitsUseNonStandardTexts + bitsSpecialFeatureOptIns + bitsPurposesConsent + bitsPurposesLITransparency +
		bitsPurposeOneTreatment + bitsPublisherCC
	for i := bitsMaxVendorId - 1; i >= 0; i-- {
		e.WriteBool(10>>i&1 == 1)
	}
	return base64.RawURLEncoding.EncodeToString(e.Bytes)
}

// What the specification allows still decodes: a restriction on vendors past MaxVendorId (it does not depend on the
// vendor sections), a restriction without vendors, a range of one vendor, and ranges written by the encoder, which
// used to leave MaxVendorId at 0 and so wrote strings that the range check refused.
func TestDecode_AcceptsValidRanges(t *testing.T) {
	core := testCore(func(c *CoreString) {
		c.PurposesConsent[1], c.PurposesConsent[4] = true, true
		c.VendorsConsent[53] = true
		c.PubRestrictions = []*PubRestriction{
			testRestriction(2, RestrictionTypeRequireConsent, 916),
			{PurposeId: 3, RestrictionType: RestrictionTypeNotAllowed},
		}
	}).Encode()
	if _, err := Decode(core + "." + testDisclosedVendors(53)); err != nil {
		t.Errorf("restrictions past MaxVendorId and without vendors: %v", err)
	}
	if _, err := Decode(testWithRawRestriction(core, 4, int(RestrictionTypeNotAllowed), 53, 53) + "." + testDisclosedVendors(53)); err != nil {
		t.Errorf("range of one vendor: %v", err)
	}

	ranges := testCore(func(c *CoreString) {
		c.PurposesConsent[1] = true
		c.IsRangeEncoding, c.RangeEntries = true, []*RangeEntry{{StartVendorID: 1, EndVendorID: 10}}
		c.IsRangeEncodingLI, c.RangeEntriesLI = true, []*RangeEntry{{StartVendorID: 50, EndVendorID: 60}}
	})
	d, err := Decode(ranges.Encode() + "." + testDisclosedVendors(5))
	if err != nil {
		t.Fatalf("range-encoded vendors written by the encoder: %v", err)
	}
	if d.CoreString.MaxVendorId != 10 || d.CoreString.MaxVendorIdLI != 60 || !d.IsVendorAllowed(5) || !d.IsVendorLIAllowed(55) {
		t.Errorf("range-encoded vendors: MaxVendorId %d, MaxVendorIdLI %d, vendor 5 %v, vendor 55 LI %v",
			d.CoreString.MaxVendorId, d.CoreString.MaxVendorIdLI, d.IsVendorAllowed(5), d.IsVendorLIAllowed(55))
	}
}

// Decode used to leave out a Disclosed Vendors or Publisher TC segment that did not decode, as if the string had none.
// It now refuses the string; DecodeLenient still leaves the segment out.
func TestDecode_RejectsMalformedSegments(t *testing.T) {
	core := testCore(func(c *CoreString) { c.PurposesConsent[1] = true }).Encode()
	invalidDisclosed := (&DisclosedVendors{SegmentType: int(SegmentTypeDisclosedVendors), IsRangeEncoding: true,
		RangeEntries: []*RangeEntry{{StartVendorID: 0, EndVendorID: 60}}}).Encode()
	const truncatedPublisherTC = "YA" // segment type 3, then nothing
	for name, c := range map[string]struct{ s, want string }{
		"invalid Disclosed Vendors range": {core + "." + invalidDisclosed, "invalid Disclosed Vendors segment"},
		"truncated Publisher TC":          {core + "." + testDisclosedVendors(53) + "." + truncatedPublisherTC, "invalid Publisher TC segment"},
	} {
		if _, err := Decode(c.s); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: Decode expected %q, got %v", name, c.want, err)
		}
	}
	d, err := DecodeLenient(core + "." + testDisclosedVendors(53) + "." + truncatedPublisherTC)
	if err != nil || d.PublisherTC != nil || d.DisclosedVendors == nil {
		t.Errorf("truncated Publisher TC: DecodeLenient should leave out that segment only, got %v", err)
	}
}

// DecodeLenient reads the strings that Decode now refuses for their ranges or restriction types, as before: a
// withdrawal of consent is worth reading even from an off-spec string.
func TestDecodeLenient_SkipsTheChecksOfDecode(t *testing.T) {
	core := testCore(func(c *CoreString) {
		c.PurposesConsent[1], c.PurposesConsent[4] = true, true
		c.VendorsConsent[53] = true
	}).Encode()
	invalidDisclosed := (&DisclosedVendors{SegmentType: int(SegmentTypeDisclosedVendors), IsRangeEncoding: true,
		RangeEntries: []*RangeEntry{{StartVendorID: 0, EndVendorID: 60}}}).Encode()
	for name, s := range map[string]string{
		"inverted restriction range":      testWithRawRestriction(core, 4, int(RestrictionTypeNotAllowed), 60, 40) + "." + testDisclosedVendors(53),
		"undefined restriction type":      testWithRawRestriction(core, 4, int(RestrictionTypeUndefined), 53, 53) + "." + testDisclosedVendors(53),
		"vendor range beyond MaxVendorId": testRangeBeyondMaxVendorId() + "." + testDisclosedVendors(53),
		"invalid Disclosed Vendors range": core + "." + invalidDisclosed,
	} {
		if _, err := Decode(s); err == nil {
			t.Errorf("%s: Decode accepted the string", name)
		}
		d, err := DecodeLenient(s)
		if err != nil {
			t.Errorf("%s: DecodeLenient: %v", name, err)
			continue
		}
		if !d.IsPurposeAllowed(1) {
			t.Errorf("%s: DecodeLenient did not read the consent to purpose 1", name)
		}
	}
}

func TestDecoders_ReturnNothingWithAnError(t *testing.T) {
	core := testCore(func(c *CoreString) { c.PurposesConsent[1] = true }).Encode()
	if c, err := DecodeCoreString(core[:20]); err == nil || c != nil {
		t.Errorf("truncated core string: got %v, %v", c, err)
	}
	if segmentType, err := GetSegmentType(""); err == nil || segmentType != SegmentTypeUndefined {
		t.Errorf("empty segment: got type %v, %v", segmentType, err)
	}
	if version, err := GetVersion(""); err == nil || version != TcfVersionUndefined {
		t.Errorf("empty string: got version %v, %v", version, err)
	}
}
