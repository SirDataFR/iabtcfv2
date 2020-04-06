# iab-tcf-v2
Go client library to read and encode IAB TCF V2.0 TC Strings.
####Installation

```
go get github.com/SirDataFR/iabtcfv2
```

The package defines a `TCData` structure with the four segments a TC String can contain:
- `CoreString`
- `DisclosedVendors`
- `AllowedVendors`
- `PublisherTC`

#### Decode a TC String

To decode a TC String, use the `Decode(tcString string)` function.

NOTE : A valid TC String must start with a *Core String* segment value.
```
var tcData, err = iabtcfv2.Decode("COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA.IF3EXySoGY2tho2YVFzBEIYwfJxyigMgShgQIsS0NQIeFLBoGPiAAHBGYJAQAGBAkkACBAQIsHGBMCQABgAgRiRCMQEGMDzNIBIBAggkbY0FACCVmnkHS3ZCY70-6u__QA.elAAAAAAAWA")
if err != nil {
  fmt.Printf("%v", err)
}
```

To decode a segment value of a TC String, use the appropriate function:
- `DecodeCoreString(coreString string)`
- `DecodeDisclosedVendors(disclosedVendors string)`
- `DecodeAllowedVendors(allowedVendors string)`
- `DecodePublisherTC(publisherTC string)`
```
var coreString, err = iabtcfv2.DecodeCoreString("COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA")
if err != nil {
  fmt.Printf("%v", err)
}
```

Use `GetSegmentType(segment string)` to read the segment type from a segment value:
- `0` = *Core String*
- `1` = *Disclosed Vendors*
- `2` = *Allowed Vendors*
- `3` = *Publisher TC*

Use `GetVersion(s string)` to verify the cookie version from a TC String or a *Core String* segment value. This function also supports TCF V1.1 consent strings:
- `1` = TCF V1.1
- `2` = TCF V2.0

#####Example
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
  fmt.Printf("%+v\n", tcData.AllowedVendors)
  fmt.Printf("%+v\n", tcData.PublisherTC)
}
```

#### Encode a TC String

To encode a TC String, use the `ToTCString()` function on the `TCData` structure.

To encode a single segment of a TC String, use the `Encode()` function on the appropriate segment.

#####Example
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
      TcfPolicyVersion: 2,
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

#### Read TC Data

To verify that a legal basis is established for a purpose or a vendor, use the functions on each structure.

##### CoreString
| Function                 | Parameter        | Description           |
| ------------------------ | :--------------: | --------------------- |
| IsSpecialFeatureAllowed  | int | Returns `true` if user has given consent to special feature id |
| IsPurposeAllowed         | int | Returns `true` if user has given consent to purpose id |
| IsPurposeLIAllowed       | int | Returns `true` if legitimate interest is established for purpose id and user didn't exercise their right to object | 
| IsVendorAllowed          | int | Returns `true` if user has given consent to vendor id processing their personal data |
| IsVendorLIAllowed        | int | Returns `true` if transparency for vendor id's legitimate interest is established and user didn't exercise their right to object |

NOTE: For convenience the `CoreString` functions are also available from the `TCData` structure.

##### DisclosedVendors
| Function                 | Parameter        | Description           |
| ------------------------ | :--------------: | --------------------- |
| IsVendorDisclosed        | int | Returns `true` if vendor id is disclosed for validating OOB signaling |

##### AllowedVendors
| Function                 | Parameter        | Description           |
| ------------------------ | :--------------: | --------------------- |
| IsVendorAllowed          | int | Returns `true` if vendor id is allowed for OOB signaling |

##### PublisherTC
| Function                 | Parameter        | Description           |
| ------------------------ | :--------------: | --------------------- |
| IsPurposeAllowed         | int | Returns `true` if user has given consent to standard purpose id |
| IsPurposeLIAllowed       | int | Returns `true` if legitimate interest is established for standard purpose id and user didn't exercise their right to object |
| IsCustomPurposeAllowed   | int | Returns `true` if user has given consent to custom purpose id |
| IsCustomPurposeLIAllowed | int | Returns `true` if legitimate interest is established for custom purpose id and user didn't exercise their right to object |

#####Example
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