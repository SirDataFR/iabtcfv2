package iabtcfv2

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Decodes a string and returns the TcfVersion
// It can also decode version from a TCF V1.1 consent string
// - TcfVersionUndefined = -1
// - TcfVersion1 = 1
// - TcfVersion2 = 2
func GetVersion(s string) (version TcfVersion, err error) {
	defer func() {
		if r := recover(); r != nil {
			version = TcfVersionUndefined
			err = fmt.Errorf("%v", r)
		}
	}()

	segments := strings.Split(s, ".")
	if len(segments) == 0 {
		return TcfVersionUndefined, err
	}

	b, err := base64.RawURLEncoding.DecodeString(segments[0])
	if err != nil {
		return TcfVersionUndefined, err
	}

	var e = NewTCEncoder(b)
	return TcfVersion(e.ReadInt(bitsVersion)), nil
}

// Decodes a segment value and returns the SegmentType
// - SegmentTypeUndefined = -1
// - SegmentTypeCoreString = 0
// - SegmentTypeDisclosedVendors = 1
// - SegmentTypePublisherTC = 3
func GetSegmentType(segment string) (segmentType SegmentType, err error) {
	defer func() {
		if r := recover(); r != nil {
			segmentType = SegmentTypeUndefined
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return SegmentTypeUndefined, err
	}

	var e = NewTCEncoder(b)
	return SegmentType(e.ReadInt(bitsSegmentType)), nil
}

// Decode a TC String and returns it as a TCData structure
// A valid TC String must start with a Core String segment.
// Subsequent segments (DisclosedVendors, PublisherTC) can appear in any order,
// identified by their SegmentType field.
// TCF v2.3 (policyVersion >= 5): DisclosedVendors segment is MANDATORY.
func Decode(tcString string) (t *TCData, err error) {
	t = &TCData{}
	mapSegments := map[SegmentType]bool{}
	segmentOrder := []SegmentType{}

	for i, v := range strings.Split(tcString, ".") {
		segmentType, err := GetSegmentType(v)
		if err != nil {
			return nil, err
		}

		switch segmentType {
		case SegmentTypeDisclosedVendors:
			if mapSegments[SegmentTypeDisclosedVendors] == true {
				return nil, fmt.Errorf("duplicate Disclosed Vendors segment")
			}
			segment, err := DecodeDisclosedVendors(v)
			if err != nil {
				return nil, fmt.Errorf("invalid Disclosed Vendors segment: %w", err)
			}
			t.DisclosedVendors = segment
			mapSegments[SegmentTypeDisclosedVendors] = true
			segmentOrder = append(segmentOrder, SegmentTypeDisclosedVendors)
		case SegmentTypePublisherTC:
			if mapSegments[SegmentTypePublisherTC] == true {
				return nil, fmt.Errorf("duplicate Publisher TC segment")
			}
			segment, err := DecodePublisherTC(v)
			if err != nil {
				return nil, fmt.Errorf("invalid Publisher TC segment: %w", err)
			}
			t.PublisherTC = segment
			mapSegments[SegmentTypePublisherTC] = true
			segmentOrder = append(segmentOrder, SegmentTypePublisherTC)
		case SegmentTypeCoreString:
			if mapSegments[SegmentTypeCoreString] == true {
				return nil, fmt.Errorf("duplicate Core String segment")
			}
			segment, err := DecodeCoreString(v)
			if err != nil {
				return nil, fmt.Errorf("invalid Core String segment: %w", err)
			}
			t.CoreString = segment
			if i == 0 {
				mapSegments[SegmentTypeCoreString] = true
				segmentOrder = append(segmentOrder, SegmentTypeCoreString)
			}
		default:
			return nil, fmt.Errorf("unsupported segment type %d", segmentType)
		}
	}

	if mapSegments[SegmentTypeCoreString] == false {
		return nil, fmt.Errorf("invalid TC string")
	}

	// TCF v2.3: After Feb 28, 2026, TC strings must have TcfPolicyVersion >= 5
	// AND a DisclosedVendors segment. Before March 1, 2026, DV is optional.
	if t.CoreString.Created.After(v23Deadline) {
		if t.CoreString.TcfPolicyVersion < TcfPolicyVersion23 {
			return nil, fmt.Errorf("TCF v2.3: TC Strings created after %s must have policyVersion >= %d", v23Deadline.Format("2006-01-02"), TcfPolicyVersion23)
		}
		if !mapSegments[SegmentTypeDisclosedVendors] {
			return nil, fmt.Errorf("TCF v2.3: DisclosedVendors segment is mandatory for TC Strings created after %s", v23Deadline.Format("2006-01-02"))
		}
		// Core String must be the first segment
		if segmentOrder[0] != SegmentTypeCoreString {
			return nil, fmt.Errorf("TCF v2.3: Core String must be the first segment")
		}
	}

	return t, nil
}

// DecodeLenient decodes a TC String without enforcing TCF v2.3 deadline/policyVersion checks.
// Accepts any well-formed TC string regardless of creation date or policy version. Unlike Decode, it reads the core
// string and the Disclosed Vendors segment without the checks of DecodeCoreString and DecodeDisclosedVendors
// (version, vendor ranges, restriction types), and leaves out a Disclosed Vendors or Publisher TC segment that does
// not decode: it is meant for the signals worth reading even from an off-spec string, such as a withdrawal.
func DecodeLenient(tcString string) (t *TCData, err error) {
	t = &TCData{}
	mapSegments := map[SegmentType]bool{}
	segmentOrder := []SegmentType{}

	for i, v := range strings.Split(tcString, ".") {
		segmentType, err := GetSegmentType(v)
		if err != nil {
			return nil, err
		}

		switch segmentType {
		case SegmentTypeDisclosedVendors:
			if mapSegments[SegmentTypeDisclosedVendors] {
				return nil, fmt.Errorf("duplicate Disclosed Vendors segment")
			}
			segment, err := decodeDisclosedVendors(v, false)
			if err == nil {
				t.DisclosedVendors = segment
				mapSegments[SegmentTypeDisclosedVendors] = true
				segmentOrder = append(segmentOrder, SegmentTypeDisclosedVendors)
			}
		case SegmentTypePublisherTC:
			if mapSegments[SegmentTypePublisherTC] {
				return nil, fmt.Errorf("duplicate Publisher TC segment")
			}
			segment, err := DecodePublisherTC(v)
			if err == nil {
				t.PublisherTC = segment
				mapSegments[SegmentTypePublisherTC] = true
				segmentOrder = append(segmentOrder, SegmentTypePublisherTC)
			}
		case SegmentTypeCoreString:
			if mapSegments[SegmentTypeCoreString] {
				return nil, fmt.Errorf("duplicate Core String segment")
			}
			segment, err := decodeCoreString(v, false)
			if err != nil {
				return nil, fmt.Errorf("invalid Core String segment: %w", err)
			}
			t.CoreString = segment
			if i == 0 {
				mapSegments[SegmentTypeCoreString] = true
				segmentOrder = append(segmentOrder, SegmentTypeCoreString)
			}
		default:
			return nil, fmt.Errorf("unsupported segment type %d", segmentType)
		}
	}

	if !mapSegments[SegmentTypeCoreString] {
		return nil, fmt.Errorf("invalid TC string")
	}

	// NO deadline/policyVersion checks — lenient mode
	_ = segmentOrder
	return t, nil
}

// Decodes a Core String value and returns it as a CoreString structure. It rejects what it would read wrong: another
// version than 2, a vendor range that starts before vendor 1 or ends before it starts, a range of the vendor consents
// or legitimate interests that goes beyond their MaxVendorId, and an undefined publisher restriction type.
func DecodeCoreString(coreString string) (c *CoreString, err error) {
	return decodeCoreString(coreString, true)
}

// decodeCoreString decodes a core string, with the checks of DecodeCoreString if strict is true.
func decodeCoreString(coreString string, strict bool) (c *CoreString, err error) {
	defer func() {
		if r := recover(); r != nil {
			c = nil
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(coreString)
	if err != nil {
		return nil, err
	}

	var e = NewTCEncoder(b)

	c = &CoreString{}
	c.Version = e.ReadInt(bitsVersion)
	if strict && c.Version != int(TcfVersion2) {
		panic(fmt.Sprintf("unsupported TC string version %d", c.Version))
	}
	c.Created = e.ReadTime()
	c.LastUpdated = e.ReadTime()
	c.CmpId = e.ReadInt(bitsCmpId)
	c.CmpVersion = e.ReadInt(bitsCmpVersion)
	c.ConsentScreen = e.ReadInt(bitsConsentScreen)
	c.ConsentLanguage = e.ReadChars(bitsConsentLanguage)
	c.VendorListVersion = e.ReadInt(bitsVendorListVersion)
	c.TcfPolicyVersion = e.ReadInt(bitsTcfPolicyVersion)
	c.IsServiceSpecific = e.ReadBool()
	c.UseNonStandardTexts = e.ReadBool()
	c.SpecialFeatureOptIns = e.ReadBitField(bitsSpecialFeatureOptIns)
	c.PurposesConsent = e.ReadBitField(bitsPurposesConsent)
	c.PurposesLITransparency = e.ReadBitField(bitsPurposesLITransparency)
	c.PurposeOneTreatment = e.ReadBool()
	c.PublisherCC = e.ReadChars(bitsPublisherCC)

	c.MaxVendorId = e.ReadInt(bitsMaxVendorId)
	c.IsRangeEncoding = e.ReadBool()
	if c.IsRangeEncoding {
		c.NumEntries, c.RangeEntries = e.ReadRangeEntries()
	} else {
		c.VendorsConsent = e.ReadBitField(uint(c.MaxVendorId))
	}

	c.MaxVendorIdLI = e.ReadInt(bitsMaxVendorId)
	c.IsRangeEncodingLI = e.ReadBool()
	if c.IsRangeEncodingLI {
		c.NumEntriesLI, c.RangeEntriesLI = e.ReadRangeEntries()
	} else {
		c.VendorsLITransparency = e.ReadBitField(uint(c.MaxVendorIdLI))
	}

	c.NumPubRestrictions, c.PubRestrictions = e.ReadPubRestrictions()

	if strict {
		checkRangesWithin(c.RangeEntries, c.MaxVendorId)
		checkRangesWithin(c.RangeEntriesLI, c.MaxVendorIdLI)
		for _, r := range c.PubRestrictions {
			if r.RestrictionType == RestrictionTypeUndefined {
				panic(fmt.Sprintf("undefined publisher restriction type %d for purpose %d", r.RestrictionType, r.PurposeId))
			}
			checkRanges(r.RangeEntries)
		}
	}

	return c, nil
}

// Decodes a Disclosed Vendors value and returns it as a DisclosedVendors structure. It rejects a vendor range that
// starts before vendor 1 or ends before it starts.
func DecodeDisclosedVendors(disclosedVendors string) (d *DisclosedVendors, err error) {
	return decodeDisclosedVendors(disclosedVendors, true)
}

// decodeDisclosedVendors decodes a Disclosed Vendors segment, with the checks of DecodeDisclosedVendors if strict is
// true.
func decodeDisclosedVendors(disclosedVendors string, strict bool) (d *DisclosedVendors, err error) {
	defer func() {
		if r := recover(); r != nil {
			d = nil
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(disclosedVendors)
	if err != nil {
		return nil, err
	}

	var e = NewTCEncoder(b)

	d = &DisclosedVendors{}
	d.SegmentType = e.ReadInt(bitsSegmentType)
	d.MaxVendorId = e.ReadInt(bitsMaxVendorId)
	d.IsRangeEncoding = e.ReadBool()
	if d.IsRangeEncoding {
		d.NumEntries, d.RangeEntries = e.ReadRangeEntries()
		if strict {
			checkRanges(d.RangeEntries)
		}
	} else {
		d.DisclosedVendors = e.ReadBitField(uint(d.MaxVendorId))
	}

	if d.SegmentType != int(SegmentTypeDisclosedVendors) {
		err = fmt.Errorf("disclosed vendors segment type must be %d", SegmentTypeDisclosedVendors)
		return nil, err
	}

	return d, nil
}

// Decodes a Publisher TC value and returns it as a PublisherTC structure
func DecodePublisherTC(publisherTC string) (p *PublisherTC, err error) {
	defer func() {
		if r := recover(); r != nil {
			p = nil
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(publisherTC)
	if err != nil {
		return nil, err
	}

	var e = NewTCEncoder(b)

	p = &PublisherTC{}
	p.SegmentType = e.ReadInt(bitsSegmentType)
	p.PubPurposesConsent = e.ReadBitField(bitsPubPurposesConsent)
	p.PubPurposesLITransparency = e.ReadBitField(bitsPubPurposesLITransparency)
	p.NumCustomPurposes = e.ReadInt(bitsNumCustomPurposes)
	p.CustomPurposesConsent = e.ReadBitField(uint(p.NumCustomPurposes))
	p.CustomPurposesLITransparency = e.ReadBitField(uint(p.NumCustomPurposes))

	if p.SegmentType != int(SegmentTypePublisherTC) {
		return nil, fmt.Errorf("publisher TC segment type must be %d", SegmentTypePublisherTC)
	}

	return p, nil
}

// checkRanges panics on a range of vendor ids that the decoder would read wrong. Vendor ids start at 1 and a range ends
// after it starts: a range that does not matches no vendor, and a restriction written with it silently restricts
// nothing.
func checkRanges(entries []*RangeEntry) {
	for _, entry := range entries {
		if entry.StartVendorID < 1 || entry.EndVendorID < entry.StartVendorID {
			panic(fmt.Sprintf("invalid vendor range %d-%d", entry.StartVendorID, entry.EndVendorID))
		}
	}
}

// checkRangesWithin also panics if a range goes beyond maxVendorId: the vendors past it are not in the vendor section
// that declares it.
func checkRangesWithin(entries []*RangeEntry, maxVendorId int) {
	checkRanges(entries)
	for _, entry := range entries {
		if entry.EndVendorID > maxVendorId {
			panic(fmt.Sprintf("vendor range %d-%d beyond max vendor id %d", entry.StartVendorID, entry.EndVendorID, maxVendorId))
		}
	}
}
