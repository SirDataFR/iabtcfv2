package iabtcfv2

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	version, err := GetVersion("BOr70tQOxPQw-BcAsCFRDEqAAAAu1rxyZn7kfUXiXSZxNuiGGp6h-Wd9CWUcKZYpMAnyhYZRfg_AQhQ4Eu0LRNNycgh45MoCCMoRQaiSkCABGgFcTpjTmxAUxoRLawAMBrwhWLEQeroyHcJzAAHN_QjACAA")
	if err != nil {
		t.Errorf("Version should be decoded without error: %s", err)
		return
	}

	if version != 1 {
		t.Errorf("Version should be 1")
	}

	version, err = GetVersion("COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA")
	if err != nil {
		t.Errorf("Version should be decoded without error: %s", err)
		return
	}

	if version != 2 {
		t.Errorf("Version should be 2")
	}
}

func TestGetSegmentType(t *testing.T) {
	str := "IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70"

	segType, err := GetSegmentType(str)
	if err != nil {
		t.Errorf("Segment type should be decoded without error: %s", err)
		return
	}

	if segType != 1 {
		t.Errorf("Segment type should be 1")
	}
}

func TestGetVersionAndSegmentTypeFail(t *testing.T) {
	str := "A"

	_, err := GetVersion(str)
	if err == nil {
		t.Errorf("Version should not be decoded")
		return
	}

	_, err = GetSegmentType(str)
	if err == nil {
		t.Errorf("Segment type should not be decoded")
	}
}

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

func TestDecodeMissingCore(t *testing.T) {
	str := "IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70-6u__QA.elAAAAAAAWA"

	_, err := Decode(str)
	if err == nil {
		t.Errorf("TC String should not be decoded: %s", err)
		return
	}
}

func TestDecodeWrongOrdered(t *testing.T) {
	str := "elAAAAAAAWA.COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA.IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70-6u__QA"

	_, err := Decode(str)
	if err == nil {
		t.Errorf("TC String should not be decoded: %s", err)
		return
	}
}

func TestDecodeDuplicateSegment(t *testing.T) {
	str := "COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA.COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA"

	_, err := Decode(str)
	if err == nil {
		t.Errorf("TC String should not be decoded: %s", err)
		return
	}
}

func TestDecodeCoreString(t *testing.T) {
	str := "COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA"

	segType, err := GetSegmentType(str)
	if err != nil {
		t.Errorf("Segment type should be decoded without error: %s", err)
		return
	}

	if segType != 0 {
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

	segType, err := GetSegmentType(str)
	if err != nil {
		t.Errorf("Segment type should be decoded without error: %s", err)
		return
	}

	if segType != 1 {
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

	segType, err := GetSegmentType(str)
	if err != nil {
		t.Errorf("Segment type should be decoded without error: %s", err)
		return
	}

	if segType != 2 {
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

	segType, err := GetSegmentType(str)
	if err != nil {
		t.Errorf("Segment type should be decoded without error: %s", err)
		return
	}

	if segType != 3 {
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
