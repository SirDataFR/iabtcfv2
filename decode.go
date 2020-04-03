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

	var e = NewTCEncoder(b)
	return e.ReadInt(3)
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

	var e = NewTCEncoder(b)

	c = &CoreString{}
	c.Version = e.ReadInt(6)
	c.Created = e.ReadTime()
	c.LastUpdated = e.ReadTime()
	c.CmpId = e.ReadInt(12)
	c.CmpVersion = e.ReadInt(12)
	c.ConsentScreen = e.ReadInt(6)
	c.ConsentLanguage = e.ReadIsoCode()
	c.VendorListVersion = e.ReadInt(12)
	c.TcfPolicyVersion = e.ReadInt(6)
	c.IsServiceSpecific = e.ReadBool()
	c.UseNonStandardStacks = e.ReadBool()
	c.SpecialFeatureOptIns = e.ReadBitField(12)
	c.PurposesConsent = e.ReadBitField(24)
	c.PurposesLITransparency = e.ReadBitField(24)
	c.PurposeOneTreatment = e.ReadBool()
	c.PublisherCC = e.ReadIsoCode()

	c.MaxVendorId = e.ReadInt(16)
	c.IsRangeEncoding = e.ReadBool()
	if c.IsRangeEncoding {
		c.NumEntries = e.ReadInt(12)
		c.RangeEntries = e.ReadRangeEntries(uint(c.NumEntries))
	} else {
		c.VendorsConsent = e.ReadBitField(uint(c.MaxVendorId))
	}

	c.MaxVendorIdLI = e.ReadInt(16)
	c.IsRangeEncodingLI = e.ReadBool()
	if c.IsRangeEncodingLI {
		c.NumEntriesLI = e.ReadInt(12)
		c.RangeEntriesLI = e.ReadRangeEntries(uint(c.NumEntriesLI))
	} else {
		c.VendorsLITransparency = e.ReadBitField(uint(c.MaxVendorIdLI))
	}

	c.NumPubRestrictions = e.ReadInt(12)
	c.PubRestrictions = e.ReadPubRestrictions(uint(c.NumPubRestrictions))

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

	var e = NewTCEncoder(b)

	d = &DisclosedVendors{}
	d.SegmentType = e.ReadInt(3)
	d.MaxVendorId = e.ReadInt(16)
	d.IsRangeEncoding = e.ReadBool()
	if d.IsRangeEncoding {
		d.NumEntries = e.ReadInt(12)
		d.RangeEntries = e.ReadRangeEntries(uint(d.NumEntries))
	} else {
		d.DisclosedVendors = e.ReadBitField(uint(d.MaxVendorId))
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

	var e = NewTCEncoder(b)

	a = &AllowedVendors{}
	a.SegmentType = e.ReadInt(3)
	a.MaxVendorId = e.ReadInt(16)
	a.IsRangeEncoding = e.ReadBool()
	if a.IsRangeEncoding {
		a.NumEntries = e.ReadInt(12)
		a.RangeEntries = e.ReadRangeEntries(uint(a.NumEntries))
	} else {
		a.AllowedVendors = e.ReadBitField(uint(a.MaxVendorId))
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

	var e = NewTCEncoder(b)

	p = &PublisherTC{}
	p.SegmentType = e.ReadInt(3)
	p.PubPurposesConsent = e.ReadBitField(24)
	p.PubPurposesLITransparency = e.ReadBitField(24)
	p.NumCustomPurposes = e.ReadInt(6)
	p.CustomPurposesConsent = e.ReadBitField(uint(p.NumCustomPurposes))
	p.CustomPurposesLITransparency = e.ReadBitField(uint(p.NumCustomPurposes))

	if p.SegmentType != 3 {
		err = fmt.Errorf("allowed vendors segment type must be 3")
		return nil, err
	}

	return p, nil
}
