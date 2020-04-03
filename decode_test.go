package iabtcfv2

import (
	"testing"
)

func TestDecode(t *testing.T) {
	str := "COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA.IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70-6u__QA.elAAAAAAAWA"

	data, err := Decode(str)
	if err != nil {
		t.Errorf("TC String should be decoded without error: %s", err)
		return
	}

	result := data.ToTCString()
	if result == "" {
		t.Errorf("Encode() should be produce a string")
		return
	}

	if result != str {
		t.Errorf("Encode() should produce the same string: in = %s, out = %s", str, result)
	}
}

func TestDecodeInvalid(t *testing.T) {
	str := "IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70-6u__QA.elAAAAAAAWA"

	_, err := Decode(str)
	if err == nil {
		t.Errorf("TC String should not be decoded: %s", err)
		return
	}
}

func TestDecodeCoreString(t *testing.T) {
	str := "COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA"

	if DecodeSegmentType(str) != 0 {
		t.Errorf("Segment type should be 0")
		return
	}

	segment, err := DecodeCoreString(str)
	if err != nil {
		t.Errorf("Segment should be decoded without error: %s", err)
		return
	}

	result := segment.Encode()
	if result == "" {
		t.Errorf("Encode() should be produce a string")
		return
	}

	if result != str {
		t.Errorf("Encode() should produce the same string: in = %s, out = %s", str, result)
	}
}

func TestDecodeDisclosedVendors(t *testing.T) {
	str := "IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70-6u__QA"

	if DecodeSegmentType(str) != 1 {
		t.Errorf("Segment type should be 1")
		return
	}

	segment, err := DecodeDisclosedVendors(str)
	if err != nil {
		t.Errorf("Segment should be decoded without error: %s", err)
		return
	}

	if segment.IsVendorDisclosed(1) {
		t.Errorf("Vendor 1 should not be disclosed")
		return
	}

	if !segment.IsVendorDisclosed(9) {
		t.Errorf("Vendor 9 should be disclosed")
		return
	}

	result := segment.Encode()
	if result == "" {
		t.Errorf("Encode() should be produce a string")
		return
	}

	if result != str {
		t.Errorf("Encode() should produce the same string: in = %s, out = %s", str, result)
	}
}

func TestDecodeAllowedVendors(t *testing.T) {
	str := "QF3QAgABAA1A"

	if DecodeSegmentType(str) != 2 {
		t.Errorf("Segment type should be 2")
		return
	}

	segment, err := DecodeAllowedVendors(str)
	if err != nil {
		t.Errorf("Segment should be decoded without error: %s", err)
		return
	}

	if segment.IsVendorAllowed(10) {
		t.Errorf("Vendor 10 should not be disclosed")
		return
	}

	if !segment.IsVendorAllowed(53) {
		t.Errorf("Vendor 53 should be disclosed")
		return
	}

	result := segment.Encode()
	if result == "" {
		t.Errorf("Encode() should be produce a string")
		return
	}

	if result != str {
		t.Errorf("Encode() should produce the same string: in = %s, out = %s", str, result)
	}
}

func TestDecodePublisherTC(t *testing.T) {
	str := "elAAAAAAAWA"

	if DecodeSegmentType(str) != 3 {
		t.Errorf("Segment type should be 3")
		return
	}

	segment, err := DecodePublisherTC(str)
	if err != nil {
		t.Errorf("Segment should be decoded without error: %s", err)
		return
	}

	if !segment.IsPurposeAllowed(1) {
		t.Errorf("Purpose 1 should be allowed")
		return
	}

	if segment.NumCustomPurposes != 2 {
		t.Errorf("NumCustomPurposes should be 2")
	}

	result := segment.Encode()
	if result == "" {
		t.Errorf("Encode() should be produce a string")
		return
	}

	if result != str {
		t.Errorf("Encode() should produce the same string: in = %s, out = %s", str, result)
	}
}
