# iab-tcf-v2
Go client library to read and encode IAB TCF v2.0/v2.2/v2.3 TC Strings.
### Installation

```
go get github.com/SirDataFR/iabtcfv2
```

The package defines a `TCData` structure with the three segments a TC String can contain:
- `CoreString`
- `DisclosedVendors`
- `PublisherTC`

### Decode a TC String

To decode a TC String, use `Decode(tcString string) (t *TCData, err error)`.

`Decode()` enforces TCF v2.3 compliance rules: TC Strings created after **February 28, 2026** must have `TcfPolicyVersion >= 5` and include a `DisclosedVendors` segment. Strings that do not meet these requirements are rejected with an error.

Use this function in **production** to validate user consent before processing personal data.

NOTE : A valid TC String must start with a *Core String* segment value.
```
var tcData, err = iabtcfv2.Decode("COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA.IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70-6u__QA.elAAAAAAAWA")
if err != nil {
  fmt.Printf("%v", err)
}
```

#### DecodeLenient — for debugging and analysis only

> ⚠️ **`DecodeLenient()` is NOT suitable for production consent checking.**
> It bypasses TCF v2.3 compliance rules and will accept TC Strings that are
> invalid under the current policy (e.g. missing `DisclosedVendors` segment,
> `TcfPolicyVersion < 5` after the deadline). Only use it for debugging,
> logging, or forensic analysis where you need to inspect the contents of a
> non-compliant TC String. **Never use it to decide whether to process user data.**

```
// DecodeLenient: extracts TC data regardless of v2.3 compliance.
// Use ONLY for debugging/analysis — never for consent decisions in production.
tcData, err := iabtcfv2.DecodeLenient(tcString)
if err != nil {
  // format is completely invalid — not a TC string at all
  fmt.Printf("%v", err)
}

// If you need to check compliance separately:
if err := tcData.Validate(); err != nil {
  fmt.Printf("Non-compliant: %v", err)
}
```

To decode a segment value of a TC String, use the appropriate function:
- `DecodeCoreString(coreString string) (c *CoreString, err error)`
- `DecodeDisclosedVendors(disclosedVendors string) (d *DisclosedVendors, err error)`
- `DecodePublisherTC(publisherTC string) (p *PublisherTC, err error)`
```
var coreString, err = iabtcfv2.DecodeCoreString("COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA")
if err != nil {
  fmt.Printf("%v", err)
}
```

Use `GetSegmentType(segment string) (segmentType SegmentType, err error)` to read the segment type from a segment value:
- `SegmentTypeUndefined` = undefined
- `SegmentTypeCoreString` = *Core String*
- `SegmentTypeDisclosedVendors` = *Disclosed Vendors*
- `SegmentTypePublisherTC` = *Publisher TC*

You can find more information about segment types [here](https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/TCFv2/IAB%20Tech%20Lab%20-%20Consent%20string%20and%20vendor%20list%20formats%20v2.md#disclosed-vendors-oob).

Use `GetVersion(s string) (version TcfVersion, err error)` to read the cookie version from a TC String or a *Core String* segment value. This function also supports TCF v1.1 consent strings:
- `TcfVersionUndefined` = undefined
- `TcfVersion1` = TCF v1.1
- `TcfVersion2` = TCF v2.0/v2.2/v2.3
- `TcfVersion23` = TCF v2.3 (same as TcfVersion2; v2.3 rules determined by TcfPolicyVersion >= 5)

#### Example
```
package main

import (
  "fmt"
  "github.com/SirDataFR/iabtcfv2"
)

func main() {
  tcString := "COxSKBCOxSKCCBcABCENAgCMAPzAAEPAAAqIDaQBQAMgAgABqAR0A2gDaQAwAMgAgANoAAA"
  
  version, err := iabtcfv2.GetVersion(tcString)
  if err != nil {
    fmt.Printf("%v", err)
  }
  
  segmentType, err := iabtcfv2.GetSegmentType(tcString)
  if err != nil {
    fmt.Printf("%v", err)
  }
  
  var tcData, err = iabtcfv2.Decode(tcString)
  if err != nil {
    fmt.Printf("%v", err)
  }
  
  fmt.Printf("%+v\n", version)
  fmt.Printf("%+v\n", segmentType)
  fmt.Printf("%+v\n", tcData.CoreString)
  fmt.Printf("%+v\n", tcData.DisclosedVendors)
  fmt.Printf("%+v\n", tcData.PublisherTC)
}
```

### Encode a TC String

To encode a TC String, use `ToTCString() string` on the `TCData` structure.

To encode a single segment of a TC String, use `Encode() string` on the appropriate segment.

#### Example
```
package main

import (
  "fmt"
  "github.com/SirDataFR/iabtcfv2"
)

func main() {
  tcData := &TCData{
    CoreString: &CoreString{
      Version: 2,
      Created: time.Now(),
      LastUpdated: time.Now(),
      CmpId: 92,
      CmpVersion: 1,
      ConsentScreen: 1,
      ConsentLanguage: "EN",
      VendorListVersion: 32,
      TcfPolicyVersion: 5, // TCF v2.3 — DisclosedVendors segment will be mandatory
      PurposesConsent: map[int]bool{
        1:  true,
        2:  true,
        3:  true,
      },
    },
  }
  
  tcString := tcData.ToTCString()
  fmt.Printf("%v", tcString)
  
  segmentValue := tcData.CoreString.Encode()
  fmt.Printf("%v", segmentValue)
}
```

### Read TC Data

To verify that a legal basis is established for a purpose or a vendor, use the functions on each structure.

#### CoreString
| Function                 | Parameter        | Description           |
| ------------------------ | :--------------: | --------------------- |
| IsSpecialFeatureAllowed  | int | Returns `true` if user has given consent to special feature id |
| IsPurposeAllowed         | int | Returns `true` if user has given consent to purpose id |
| IsPurposeLIAllowed       | int | Returns `true` if legitimate interest is established for purpose id and user didn't exercise their right to object | 
| IsVendorAllowed          | int | Returns `true` if user has given consent to vendor id processing their personal data |
| IsVendorLIAllowed        | int | Returns `true` if transparency for vendor id's legitimate interest is established and user didn't exercise their right to object |
| IsVendorAllowedForPurposes | (int, ...int) | Returns `true` if user has given consent to vendor id processing all purposes ids and publisher hasn't set restrictions for them |
| IsVendorAllowedForPurposesLI | (int, ...int) | Returns `true` if transparency for vendor id's legitimate interest is established for all purpose ids and publisher hasn't set restrictions for them |
| IsVendorAllowedForFlexiblePurposes | (int, ...int) | Returns `true` if vendor id is allowed for all purpose ids, each one on its legal basis: consent by default, or the basis that a publisher restriction which includes the vendor sets for the purpose. For the purposes that the vendor declares flexible, with consent as default legal basis; legitimate interest never allows purpose 1, nor purposes 3 to 6 from TCF v2.2 |
| IsVendorAllowedForFlexiblePurposesLI | (int, ...int) | Returns `true` if vendor id is allowed for all purpose ids, each one on its legal basis: legitimate interest by default, or the basis that a publisher restriction which includes the vendor sets for the purpose. For the purposes that the vendor declares flexible, with legitimate interest as default legal basis; legitimate interest never allows purpose 1, nor purposes 3 to 6 from TCF v2.2 |
| GetPubRestrictionsForPurpose | int | Returns a list of publisher restrictions applied to purpose id |

NOTE: For convenience the `CoreString` functions are also available from the `TCData` structure.

#### DisclosedVendors
| Function                 | Parameter        | Description           |
| ------------------------ | :--------------: | --------------------- |
| IsVendorDisclosed        | int | Returns `true` if vendor id is disclosed for validating OOB signaling |

#### PublisherTC
| Function                 | Parameter        | Description           |
| ------------------------ | :--------------: | --------------------- |
| IsPurposeAllowed         | int | Returns `true` if user has given consent to standard purpose id |
| IsPurposeLIAllowed       | int | Returns `true` if legitimate interest is established for standard purpose id and user didn't exercise their right to object |
| IsCustomPurposeAllowed   | int | Returns `true` if user has given consent to custom purpose id |
| IsCustomPurposeLIAllowed | int | Returns `true` if legitimate interest is established for custom purpose id and user didn't exercise their right to object |

#### Example
```
package main

import (
  "fmt"
  "github.com/SirDataFR/iabtcfv2"
)

func main() {
  tcString := "COxSKBCOxSKCCBcABCENAgCMAPzAAEPAAAqIDaQBQAMgAgABqAR0A2gDaQAwAMgAgANoAAA.elAAAAAAAWA"
  
  var tcData, err = iabtcfv2.Decode(tcString)
  if err != nil {
    fmt.Printf("%v", err)
  }
  
  if tcData.IsPurposeAllowed(1) {
    fmt.Printf("user has given consent to purpose 1")
  }
  
  if tcData.IsVendorAllowed(53) {
    fmt.Printf("user has given consent to vendor 53")
  }
  
  if tcData.PublisherTC.IsPurposeAllowed(2) {
    fmt.Printf("user has given consent to publisher purpose 2")
  }
}
```

### TCF v2.3 Compliance

As of TCF v2.3 (policyVersion >= 5, mandatory from Feb 28, 2026), the following rules apply:

#### Decode() enforcement policy

| TC String creation date | TcfPolicyVersion | Has DisclosedVendors? | `Decode()` result |
|------------------------|------------------|----------------------|-------------------|
| Before March 1, 2026 | Any | Any | ✅ Valid |
| After Feb 28, 2026 | >= 5 | Yes | ✅ Valid |
| After Feb 28, 2026 | >= 5 | No | ❌ Invalid |
| After Feb 28, 2026 | < 5 | Any | ❌ Invalid |

**`DecodeLenient()`** bypasses these checks — see warning above.

#### What changed in v2.3

| Change | Description |
|--------|-------------|
| **DisclosedVendors mandatory** | Must appear immediately after Core String for TC Strings created after Feb 28, 2026 |
| **Deadline validation** | TC Strings created after 2026-02-28 without DV or with `TcfPolicyVersion < 5` are rejected by `Decode()` |
| **Purpose 1 LI restricted** | `IsPurposeLIAllowed(1)` returns `false` for v2.3 strings — Purpose 1 cannot be established via legitimate interest |
| **Segment order enforced** | Strict order: `[Core String].[DisclosedVendors].[PublisherTC]` |

#### Validate a TC String for v2.3 compliance

Use `Validate() error` on the `TCData` structure. `Validate()` applies the same deadline-based rules as `Decode()`:
```
var tcData, err = iabtcfv2.Decode(tcString)
if err != nil {
  fmt.Printf("Decode error: %v", err)
}
if err := tcData.Validate(); err != nil {
  fmt.Printf("Validation error: %v", err)
}
```

#### Check if a TC String follows v2.3 rules

Use `IsV23() bool` on the `TCData` structure:
```
if tcData.IsV23() {
  // TcfPolicyVersion >= 5 — v2.3 rules apply
}
```

#### Check if a vendor is disclosed

Use `IsVendorDisclosed(id int) bool` on the `TCData` structure:
```
if tcData.IsVendorDisclosed(42) {
  // vendor 42 is present in the DisclosedVendors segment
}
```

#### Encode a v2.3 compliant TC String

When `TcfPolicyVersion >= 5`, the `DisclosedVendors` segment is automatically included even if not explicitly set:
```
tcData := &TCData{
  CoreString: &CoreString{
    Version: 2,
    TcfPolicyVersion: 5, // v2.3 — DisclosedVendors will be auto-created
    // ... other fields
  },
  DisclosedVendors: nil, // Will be auto-created as empty DV segment
}
tcString := tcData.ToTCString() // Always includes DV segment for v2.3
```

#### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `TcfVersion23` | 2 | TCF v2.3 version (same as TcfVersion2; v2.3 rules determined by policyVersion) |
| `TcfPolicyVersion23` | 5 | Policy version introduced by TCF v2.3 |
| `TcfV23Deadline` | 1772236800 | Unix timestamp for TCF v2.3 enforcement deadline (Feb 28, 2026 00:00:00 UTC) |
