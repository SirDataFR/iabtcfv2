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

	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return TcfVersionUndefined, err
	}

	var e = newTCEncoder(b)
	return TcfVersion(e.readInt(6)), nil
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
	return SegmentType(e.readInt(3)), nil
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
	c.Version = e.readInt(6)
	c.Created = e.readTime()
	c.LastUpdated = e.readTime()
	c.CmpId = e.readInt(12)
	c.CmpVersion = e.readInt(12)
	c.ConsentScreen = e.readInt(6)
	c.ConsentLanguage = e.readIsoCode()
	c.VendorListVersion = e.readInt(12)
	c.TcfPolicyVersion = e.readInt(6)
	c.IsServiceSpecific = e.readBool()
	c.UseNonStandardStacks = e.readBool()
	c.SpecialFeatureOptIns = e.readBitField(12)
	c.PurposesConsent = e.readBitField(24)
	c.PurposesLITransparency = e.readBitField(24)
	c.PurposeOneTreatment = e.readBool()
	c.PublisherCC = e.readIsoCode()

	c.MaxVendorId = e.readInt(16)
	c.IsRangeEncoding = e.readBool()
	if c.IsRangeEncoding {
		c.NumEntries = e.readInt(12)
		c.RangeEntries = e.readRangeEntries(uint(c.NumEntries))
	} else {
		c.VendorsConsent = e.readBitField(uint(c.MaxVendorId))
	}

	c.MaxVendorIdLI = e.readInt(16)
	c.IsRangeEncodingLI = e.readBool()
	if c.IsRangeEncodingLI {
		c.NumEntriesLI = e.readInt(12)
		c.RangeEntriesLI = e.readRangeEntries(uint(c.NumEntriesLI))
	} else {
		c.VendorsLITransparency = e.readBitField(uint(c.MaxVendorIdLI))
	}

	c.NumPubRestrictions = e.readInt(12)
	c.PubRestrictions = e.readPubRestrictions(uint(c.NumPubRestrictions))

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
	d.SegmentType = e.readInt(3)
	d.MaxVendorId = e.readInt(16)
	d.IsRangeEncoding = e.readBool()
	if d.IsRangeEncoding {
		d.NumEntries = e.readInt(12)
		d.RangeEntries = e.readRangeEntries(uint(d.NumEntries))
	} else {
		d.DisclosedVendors = e.readBitField(uint(d.MaxVendorId))
	}

	if d.SegmentType != 1 {
		err = fmt.Errorf("disclosed vendors segment type must be 1")
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
	a.SegmentType = e.readInt(3)
	a.MaxVendorId = e.readInt(16)
	a.IsRangeEncoding = e.readBool()
	if a.IsRangeEncoding {
		a.NumEntries = e.readInt(12)
		a.RangeEntries = e.readRangeEntries(uint(a.NumEntries))
	} else {
		a.AllowedVendors = e.readBitField(uint(a.MaxVendorId))
	}

	if a.SegmentType != 2 {
		return nil, fmt.Errorf("allowed vendors segment type must be 2")
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
	p.SegmentType = e.readInt(3)
	p.PubPurposesConsent = e.readBitField(24)
	p.PubPurposesLITransparency = e.readBitField(24)
	p.NumCustomPurposes = e.readInt(6)
	p.CustomPurposesConsent = e.readBitField(uint(p.NumCustomPurposes))
	p.CustomPurposesLITransparency = e.readBitField(uint(p.NumCustomPurposes))

	if p.SegmentType != 3 {
		return nil, fmt.Errorf("allowed vendors segment type must be 3")
	}

	return p, nil
}
