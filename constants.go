package iabtcfv2

type SegmentType int

const (
	SegmentTypeUndefined        SegmentType = -1
	SegmentTypeCoreString       SegmentType = 0
	SegmentTypeDisclosedVendors SegmentType = 1
	SegmentTypePublisherTC      SegmentType = 3
)

type TcfVersion int

const (
	TcfVersionUndefined TcfVersion = -1
	TcfVersion1         TcfVersion = 1
	TcfVersion2         TcfVersion = 2
	TcfVersion23        TcfVersion = 2 // TCF v2.3 uses Version=2 in the TC String; policyVersion >= 5 determines v2.3 rules
)

// TcfPolicyVersion23 is the policy version introduced by TCF v2.3
const TcfPolicyVersion23 = 5

// tcfPolicyVersion22 is the policy version introduced by TCF v2.2, which deprecates legitimate interest for purposes
// 3 to 6.
const tcfPolicyVersion22 = 4

// TcfV23Deadline is the Unix timestamp for the TCF v2.3 policy enforcement deadline (Feb 28, 2026 00:00:00 UTC)
const TcfV23Deadline int64 = 1772236800

type RestrictionType int

const (
	RestrictionTypeNotAllowed     RestrictionType = 0
	RestrictionTypeRequireConsent RestrictionType = 1
	RestrictionTypeRequireLI      RestrictionType = 2
	RestrictionTypeUndefined      RestrictionType = 3
)

const (
	bitsBool = 1
	bitsChar = 6
	bitsTime = 36

	bitsSegmentType = 3

	bitsVersion                = 6
	bitsCreated                = bitsTime
	bitsLastUpdated            = bitsTime
	bitsCmpId                  = 12
	bitsCmpVersion             = 12
	bitsConsentScreen          = 6
	bitsConsentLanguage        = bitsChar * 2
	bitsVendorListVersion      = 12
	bitsTcfPolicyVersion       = 6
	bitsIsServiceSpecific      = bitsBool
	bitsUseNonStandardTexts    = bitsBool
	bitsSpecialFeatureOptIns   = 12
	bitsPurposesConsent        = 24
	bitsPurposesLITransparency = 24
	bitsPurposeOneTreatment    = bitsBool
	bitsPublisherCC            = bitsChar * 2

	bitsMaxVendorId     = 16
	bitsIsRangeEncoding = bitsBool
	bitsNumEntries      = 12
	bitsVendorId        = 16

	bitsNumPubRestrictions                  = 12
	bitsPubRestrictionsEntryPurposeId       = 6
	bitsPubRestrictionsEntryRestrictionType = 2

	bitsPubPurposesConsent        = 24
	bitsPubPurposesLITransparency = 24
	bitsNumCustomPurposes         = 6
)
