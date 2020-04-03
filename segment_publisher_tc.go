package iabtcf

import "encoding/base64"

type PublisherTC struct {
	SegmentType                  int
	PubPurposesConsent           map[int]bool
	PubPurposesLITransparency    map[int]bool
	NumCustomPurposes            int
	CustomPurposesConsent        map[int]bool
	CustomPurposesLITransparency map[int]bool
}

func (p *PublisherTC) IsPurposeAllowed(id int) bool {
	return p.PubPurposesConsent[id]
}

func (p *PublisherTC) IsPurposeLIAllowed(id int) bool {
	return p.PubPurposesLITransparency[id]
}

func (p *PublisherTC) IsCustomPurposeAllowed(id int) bool {
	return p.CustomPurposesConsent[id]
}

func (p *PublisherTC) IsCustomPurposeLIAllowed(id int) bool {
	return p.CustomPurposesLITransparency[id]
}

func (p *PublisherTC) Encode() string {
	bitSize := 57 + (p.NumCustomPurposes * 2)

	var e = NewTCEncoder(make([]byte, bitSize/8))
	if bitSize%8 != 0 {
		e = NewTCEncoder(make([]byte, bitSize/8+1))
	}

	e.WriteInt(p.SegmentType, 3)
	for i := 0; i < 24; i++ {
		e.WriteBool(p.IsPurposeAllowed(i + 1))
	}
	for i := 0; i < 24; i++ {
		e.WriteBool(p.IsPurposeLIAllowed(i + 1))
	}
	e.WriteInt(p.NumCustomPurposes, 6)
	for i := 0; i < p.NumCustomPurposes; i++ {
		e.WriteBool(p.IsCustomPurposeAllowed(i + 1))
	}
	for i := 0; i < p.NumCustomPurposes; i++ {
		e.WriteBool(p.IsCustomPurposeLIAllowed(i + 1))
	}

	return base64.RawURLEncoding.EncodeToString(e.bytes)
}
