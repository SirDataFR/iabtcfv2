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
	bitSize := 57 + (p.NumCustomPurposes * 2)

	var e = newTCEncoder(make([]byte, bitSize/8))
	if bitSize%8 != 0 {
		e = newTCEncoder(make([]byte, bitSize/8+1))
	}

	e.writeInt(p.SegmentType, 3)
	for i := 0; i < 24; i++ {
		e.writeBool(p.IsPurposeAllowed(i + 1))
	}
	for i := 0; i < 24; i++ {
		e.writeBool(p.IsPurposeLIAllowed(i + 1))
	}
	e.writeInt(p.NumCustomPurposes, 6)
	for i := 0; i < p.NumCustomPurposes; i++ {
		e.writeBool(p.IsCustomPurposeAllowed(i + 1))
	}
	for i := 0; i < p.NumCustomPurposes; i++ {
		e.writeBool(p.IsCustomPurposeLIAllowed(i + 1))
	}

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
