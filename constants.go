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
)

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
