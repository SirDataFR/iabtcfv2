package iabtcfv2

import "encoding/base64"

type PublisherTC struct {
	SegmentType                  int
	PubPurposesConsent           map[int]bool
	PubPurposesLITransparency    map[int]bool
	NumCustomPurposes            int
	CustomPurposesConsent        map[int]bool
	CustomPurposesLITransparency map[int]bool
}

// Returns true if user has given consent to standard purpose id
func (p *PublisherTC) IsPurposeAllowed(id int) bool {
	return p.PubPurposesConsent[id]
}

// Returns true if legitimate interest is established for standard purpose id
// and user didn't exercise their right to object
func (p *PublisherTC) IsPurposeLIAllowed(id int) bool {
	return p.PubPurposesLITransparency[id]
}

// Returns true if user has given consent to custom purpose id
func (p *PublisherTC) IsCustomPurposeAllowed(id int) bool {
	return p.CustomPurposesConsent[id]
}

// Returns true if legitimate interest is established for custom purpose id
// and user didn't exercise their right to object
func (p *PublisherTC) IsCustomPurposeLIAllowed(id int) bool {
	return p.CustomPurposesLITransparency[id]
}

// Returns structure as a base64 raw url encoded string
func (p *PublisherTC) Encode() string {
	var bitSize int
	bitSize += bitsSegmentType

	bitSize += bitsPubPurposesConsent
	bitSize += bitsPubPurposesLITransparency
	bitSize += bitsNumCustomPurposes
	bitSize += p.NumCustomPurposes * 2

	e := newTCEncoderFromSize(bitSize)
	e.writeInt(p.SegmentType, bitsSegmentType)
	e.writeBools(p.IsPurposeAllowed, bitsPubPurposesConsent)
	e.writeBools(p.IsPurposeLIAllowed, bitsPubPurposesLITransparency)
	e.writeInt(p.NumCustomPurposes, bitsNumCustomPurposes)
	e.writeBools(p.IsCustomPurposeAllowed, p.NumCustomPurposes)
	e.writeBools(p.IsCustomPurposeLIAllowed, p.NumCustomPurposes)

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
