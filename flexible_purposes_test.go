package iabtcfv2

import (
	"testing"
	"time"
)

// The user consents to purposes 1 to 11 but refuses vendors 53 and 916, whose legitimate interest is not objected
// to. Purpose 4 has no legitimate interest transparency, as TCF 2.2 requires. No legal basis allows 53 for
// purposes 1 and 4, nor 916 for purpose 4 on legitimate interest.
func flexibleCore(restrictions ...*PubRestriction) *CoreString {
	return testCore(func(c *CoreString) {
		for p := 1; p <= 11; p++ {
			c.PurposesConsent[p] = true
		}
		for _, p := range []int{2, 7, 8, 9, 10, 11} {
			c.PurposesLITransparency[p] = true
		}
		c.VendorsConsent[755] = true
		c.VendorsLITransparency[53] = true
		c.VendorsLITransparency[916] = true
		c.PubRestrictions = restrictions
	})
}

// A publisher restriction that includes other vendors only used to switch off the check of the default legal basis
// for every vendor: vendors 53 and 916 were then allowed without any legal basis.
func TestFlexiblePurposes_RestrictionOnAnotherVendorChangesNothing(t *testing.T) {
	for name, core := range map[string]*CoreString{
		"no restriction":                  flexibleCore(),
		"restrictions on vendor 755 only": flexibleCore(testRestriction(1, RestrictionTypeRequireConsent, 755), testRestriction(4, RestrictionTypeRequireConsent, 755)),
	} {
		d := testDecode(t, core, 53, 755, 916)
		if d.IsVendorAllowedForFlexiblePurposes(53, 1, 4) {
			t.Errorf("%s: vendor 53 allowed for purposes 1 and 4 without consent", name)
		}
		if d.IsVendorAllowedForFlexiblePurposesLI(916, 4) {
			t.Errorf("%s: vendor 916 allowed for purpose 4 on a legitimate interest the purpose does not have", name)
		}
		if !d.IsVendorAllowedForFlexiblePurposesLI(53, 2) {
			t.Errorf("%s: vendor 53 not allowed for purpose 2 on its legitimate interest", name)
		}
	}
}

// A restriction that includes the vendor sets the legal basis of the purpose.
func TestFlexiblePurposes_RestrictionOnTheVendorSetsTheBasis(t *testing.T) {
	cases := []struct {
		name        string
		restriction *PubRestriction
		consent     bool // IsVendorAllowedForFlexiblePurposes(53, 2)
		li          bool // IsVendorAllowedForFlexiblePurposesLI(53, 2)
	}{
		{"no restriction", nil, false, true},
		{"not allowed", testRestriction(2, RestrictionTypeNotAllowed, 53), false, false},
		{"require consent", testRestriction(2, RestrictionTypeRequireConsent, 53), false, false},
		{"require legitimate interest", testRestriction(2, RestrictionTypeRequireLI, 53), true, true},
	}
	for _, c := range cases {
		var core *CoreString
		if c.restriction == nil {
			core = flexibleCore()
		} else {
			core = flexibleCore(c.restriction)
		}
		d := testDecode(t, core, 53, 755, 916)
		if got := d.IsVendorAllowedForFlexiblePurposes(53, 2); got != c.consent {
			t.Errorf("%s: IsVendorAllowedForFlexiblePurposes(53, 2) = %v, want %v", c.name, got, c.consent)
		}
		if got := d.IsVendorAllowedForFlexiblePurposesLI(53, 2); got != c.li {
			t.Errorf("%s: IsVendorAllowedForFlexiblePurposesLI(53, 2) = %v, want %v", c.name, got, c.li)
		}
	}
	// Vendor 755 has consent: a restriction requiring consent on purpose 4 allows it on either default.
	d := testDecode(t, flexibleCore(testRestriction(4, RestrictionTypeRequireConsent, 755)), 53, 755, 916)
	if !d.IsVendorAllowedForFlexiblePurposes(755, 4) || !d.IsVendorAllowedForFlexiblePurposesLI(755, 4) {
		t.Error("vendor 755 not allowed for purpose 4 although the restriction requires the consent it has")
	}
}

// prePolicy23Core returns a core string of policyVersion, created before the TCF v2.3 deadline. The user refuses
// purpose 1 and consents to purposes 2 to 6, and the legitimate interest bits of purposes 1 to 6 are set: against the
// specification for purpose 1, which never rests on legitimate interest, and from TCF v2.2 (policy version 4) for
// purposes 3 to 6. Vendor 53 has a legitimate interest, and no consent.
func prePolicy23Core(policyVersion int, restrictions ...*PubRestriction) *CoreString {
	return testCore(func(c *CoreString) {
		c.TcfPolicyVersion = policyVersion
		c.Created = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		c.LastUpdated = c.Created
		for p := 2; p <= 6; p++ {
			c.PurposesConsent[p] = true
		}
		for p := 1; p <= 6; p++ {
			c.PurposesLITransparency[p] = true
		}
		c.VendorsLITransparency[53] = true
		c.PubRestrictions = restrictions
	})
}

// Legitimate interest never allows purpose 1, nor purposes 3 to 6 from TCF v2.2, whatever the bits of the string: a
// restriction that requires it for these consent purposes means that the vendor should not process. Purpose 2 stays
// allowed on its legitimate interest, and purposes 3 to 6 on theirs before TCF v2.2.
func TestFlexiblePurposes_PurposesThatCannotBeLegitimateInterest(t *testing.T) {
	for name, core := range map[string]*CoreString{
		"no restriction": prePolicy23Core(4),
		"restrictions requiring legitimate interest": prePolicy23Core(4,
			testRestriction(1, RestrictionTypeRequireLI, 53), testRestriction(4, RestrictionTypeRequireLI, 53)),
	} {
		d := testDecode(t, core, 53)
		for _, p := range []int{1, 3, 4, 5, 6} {
			if d.IsVendorAllowedForFlexiblePurposes(53, p) || d.IsVendorAllowedForFlexiblePurposesLI(53, p) {
				t.Errorf("%s: vendor 53 allowed for purpose %d on legitimate interest", name, p)
			}
		}
		if !d.IsVendorAllowedForFlexiblePurposesLI(53, 2) {
			t.Errorf("%s: vendor 53 not allowed for purpose 2 on its legitimate interest", name)
		}
	}
	d := testDecode(t, prePolicy23Core(3), 53)
	if !d.IsVendorAllowedForFlexiblePurposesLI(53, 3, 4, 5, 6) || d.IsVendorAllowedForFlexiblePurposesLI(53, 1) {
		t.Error("TCF v2.1: purposes 3 to 6 rest on legitimate interest, purpose 1 never does")
	}
}

// A restriction includes the vendor through a range as through its id. When several restrictions on a purpose include
// the vendor, each sets a basis and all must allow it. A restriction type that the specification does not define
// allows nothing.
func TestFlexiblePurposes_RangedSeveralAndUndefinedRestrictions(t *testing.T) {
	ranged := &PubRestriction{PurposeId: 2, RestrictionType: RestrictionTypeRequireConsent, NumEntries: 1,
		RangeEntries: []*RangeEntry{{StartVendorID: 50, EndVendorID: 60}}}
	if testDecode(t, flexibleCore(ranged), 53, 755, 916).IsVendorAllowedForFlexiblePurposesLI(53, 2) {
		t.Error("range 50-60 requiring consent: vendor 53 allowed for purpose 2 on legitimate interest")
	}
	both := flexibleCore(testRestriction(2, RestrictionTypeRequireConsent, 53, 755), testRestriction(2, RestrictionTypeRequireLI, 53, 755))
	if d := testDecode(t, both, 53, 755, 916); d.IsVendorAllowedForFlexiblePurposesLI(53, 2) || d.IsVendorAllowedForFlexiblePurposes(755, 2) {
		t.Error("restrictions requiring consent and legitimate interest: a vendor allowed on one basis only")
	}
	// The decoder may refuse type 3: the restriction is set on the decoded structure.
	d := testDecode(t, flexibleCore(), 53, 755, 916)
	d.CoreString.PubRestrictions = []*PubRestriction{testRestriction(2, RestrictionTypeUndefined, 53)}
	if d.IsVendorAllowedForFlexiblePurposesLI(53, 2) {
		t.Error("undefined restriction type: vendor 53 allowed for purpose 2")
	}
}
