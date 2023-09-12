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

	var e = newTCEncoder(b)
	return TcfVersion(e.readInt(bitsVersion)), nil
}

// Decodes a segment value and returns the SegmentType
// - SegmentTypeUndefined = -1
// - SegmentTypeCoreString = 0
// - SegmentTypeDisclosedVendors = 1
// - SegmentTypeAllowedVendors = 2
// - SegmentTypePublisherTC = 3
func GetSegmentType(segment string) (segmentType SegmentType, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return SegmentTypeUndefined, err
	}

	var e = newTCEncoder(b)
	return SegmentType(e.readInt(bitsSegmentType)), nil
}

// Decode a TC String and returns it as a TCData structure
// A valid TC String must start with a Core String segment
// A TC String can optionally and arbitrarily ordered contain:
// - Disclosed Vendors
// - Allowed Vendors
// - Publisher TC
func Decode(tcString string) (t *TCData, err error) {
	t = &TCData{}
	mapSegments := map[SegmentType]bool{}
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
			if err == nil {
				t.DisclosedVendors = segment
				mapSegments[SegmentTypeDisclosedVendors] = true
			}
			break
		case SegmentTypeAllowedVendors:
			if mapSegments[SegmentTypeAllowedVendors] == true {
				return nil, fmt.Errorf("duplicate Allowed Vendors segment")
			}
			segment, err := DecodeAllowedVendors(v)
			if err == nil {
				t.AllowedVendors = segment
				mapSegments[SegmentTypeAllowedVendors] = true
			}
			break
		case SegmentTypePublisherTC:
			if mapSegments[SegmentTypePublisherTC] == true {
				return nil, fmt.Errorf("duplicate Publisher TC segment")
			}
			segment, err := DecodePublisherTC(v)
			if err == nil {
				t.PublisherTC = segment
				mapSegments[SegmentTypePublisherTC] = true
			}
			break
		default:
			if mapSegments[SegmentTypeCoreString] == true {
				return nil, fmt.Errorf("duplicate Core String segment")
			}
			segment, err := DecodeCoreString(v)
			if err == nil {
				t.CoreString = segment
				if i == 0 {
					mapSegments[SegmentTypeCoreString] = true
				}
			}
			break
		}
	}

	if mapSegments[SegmentTypeCoreString] == false {
		return nil, fmt.Errorf("invalid TC string")
	}

	return t, nil
}

// Decodes a Core String value and returns it as a CoreString structure
func DecodeCoreString(coreString string) (c *CoreString, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(coreString)
	if err != nil {
		return nil, err
	}

	var e = newTCEncoder(b)

	c = &CoreString{}
	c.Version = e.readInt(bitsVersion)
	c.Created = e.readTime()
	c.LastUpdated = e.readTime()
	c.CmpId = e.readInt(bitsCmpId)
	c.CmpVersion = e.readInt(bitsCmpVersion)
	c.ConsentScreen = e.readInt(bitsConsentScreen)
	c.ConsentLanguage = e.readChars(bitsConsentLanguage)
	c.VendorListVersion = e.readInt(bitsVendorListVersion)
	c.TcfPolicyVersion = e.readInt(bitsTcfPolicyVersion)
	c.IsServiceSpecific = e.readBool()
	c.UseNonStandardStacks = e.readBool()
	c.SpecialFeatureOptIns = e.readBitField(bitsSpecialFeatureOptIns)
	c.PurposesConsent = e.readBitField(bitsPurposesConsent)
	c.PurposesLITransparency = e.readBitField(bitsPurposesLITransparency)
	c.PurposeOneTreatment = e.readBool()
	c.PublisherCC = e.readChars(bitsPublisherCC)

	c.MaxVendorId = e.readInt(bitsMaxVendorId)
	c.IsRangeEncoding = e.readBool()
	if c.IsRangeEncoding {
		c.NumEntries, c.RangeEntries = e.readRangeEntries()
	} else {
		c.VendorsConsent = e.readBitField(uint(c.MaxVendorId))
	}

	c.MaxVendorIdLI = e.readInt(bitsMaxVendorId)
	c.IsRangeEncodingLI = e.readBool()
	if c.IsRangeEncodingLI {
		c.NumEntriesLI, c.RangeEntriesLI = e.readRangeEntries()
	} else {
		c.VendorsLITransparency = e.readBitField(uint(c.MaxVendorIdLI))
	}

	c.NumPubRestrictions, c.PubRestrictions = e.readPubRestrictions()

	return c, nil
}

// Decodes a Disclosed Vendors value and returns it as a DisclosedVendors structure
func DecodeDisclosedVendors(disclosedVendors string) (d *DisclosedVendors, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(disclosedVendors)
	if err != nil {
		return nil, err
	}

	var e = newTCEncoder(b)

	d = &DisclosedVendors{}
	d.SegmentType = e.readInt(bitsSegmentType)
	d.MaxVendorId = e.readInt(bitsMaxVendorId)
	d.IsRangeEncoding = e.readBool()
	if d.IsRangeEncoding {
		d.NumEntries, d.RangeEntries = e.readRangeEntries()
	} else {
		d.DisclosedVendors = e.readBitField(uint(d.MaxVendorId))
	}

	if d.SegmentType != int(SegmentTypeDisclosedVendors) {
		err = fmt.Errorf("disclosed vendors segment type must be %d", SegmentTypeDisclosedVendors)
		return nil, err
	}

	return d, nil
}

// Decodes a Allowed Vendors value and returns it as a AllowedVendors structure
func DecodeAllowedVendors(allowedVendors string) (a *AllowedVendors, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(allowedVendors)
	if err != nil {
		return nil, err
	}

	var e = newTCEncoder(b)

	a = &AllowedVendors{}
	a.SegmentType = e.readInt(bitsSegmentType)
	a.MaxVendorId = e.readInt(bitsMaxVendorId)
	a.IsRangeEncoding = e.readBool()
	if a.IsRangeEncoding {
		a.NumEntries, a.RangeEntries = e.readRangeEntries()
	} else {
		a.AllowedVendors = e.readBitField(uint(a.MaxVendorId))
	}

	if a.SegmentType != int(SegmentTypeAllowedVendors) {
		return nil, fmt.Errorf("allowed vendors segment type must be %d", SegmentTypeAllowedVendors)
	}

	return a, nil
}

// Decodes a Publisher TC value and returns it as a PublisherTC structure
func DecodePublisherTC(publisherTC string) (p *PublisherTC, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(publisherTC)
	if err != nil {
		return nil, err
	}

	var e = newTCEncoder(b)

	p = &PublisherTC{}
	p.SegmentType = e.readInt(bitsSegmentType)
	p.PubPurposesConsent = e.readBitField(bitsPubPurposesConsent)
	p.PubPurposesLITransparency = e.readBitField(bitsPubPurposesLITransparency)
	p.NumCustomPurposes = e.readInt(bitsNumCustomPurposes)
	p.CustomPurposesConsent = e.readBitField(uint(p.NumCustomPurposes))
	p.CustomPurposesLITransparency = e.readBitField(uint(p.NumCustomPurposes))

	if p.SegmentType != int(SegmentTypePublisherTC) {
		return nil, fmt.Errorf("publisher TC segment type must be %d", SegmentTypePublisherTC)
	}

	return p, nil
}
