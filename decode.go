package iabtcf

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func DecodeSegmentType(s string) int {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0
	}

	var e = newTCEncoder(b)
	return e.readInt(3)
}

func Decode(s string) (t *TCData, err error) {
	t = &TCData{}
	for _, v := range strings.Split(s, ".") {
		segmentType := DecodeSegmentType(v)
		if segmentType == 1 {
			segment, err := DecodeDisclosedVendors(v)
			if err == nil {
				t.DisclosedVendors = segment
			}
		} else if segmentType == 2 {
			segment, err := DecodeAllowedVendors(v)
			if err == nil {
				t.AllowedVendors = segment
			}
		} else if segmentType == 3 {
			segment, err := DecodePubllisherTC(v)
			if err == nil {
				t.PublisherTC = segment
			}
		} else {
			segment, err := DecodeCoreString(v)
			if err == nil {
				t.CoreString = segment
			}
		}
	}

	if t.CoreString == nil {
		return nil, fmt.Errorf("invalid TC string")
	}

	return t, nil
}

func DecodeCoreString(s string) (c *CoreString, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(s)
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

func DecodeDisclosedVendors(s string) (d *DisclosedVendors, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(s)
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

func DecodeAllowedVendors(s string) (a *AllowedVendors, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(s)
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
		err = fmt.Errorf("allowed vendors segment type must be 2")
		return nil, err
	}

	return a, nil
}

func DecodePubllisherTC(s string) (p *PublisherTC, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	b, err := base64.RawURLEncoding.DecodeString(s)
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
		err = fmt.Errorf("allowed vendors segment type must be 3")
		return nil, err
	}

	return p, nil
}
