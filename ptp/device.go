package ptp

import "encoding/binary"

type DataTypeCode uint16

// The most significant nibble (4 bits) is used to indicate the category of the code and whether the code value is
// standard or vendor-extended: 0101 = standard, 1101 = vendor-extended.
type DevicePropCode uint16

type DevicePropDescCode uint8
type DevicePropFormFlag uint8

const (
	DPC_Undefined DevicePropCode = 0x5000
	// DPC_BatteryLevel is a read-only property typically represented by a range of integers. The minimum field should
	// be set to the integer used for no power (example 0), and the maximum should be set to the integer used for full
	// power (example 100). The step field, or the individual thresholds in an enumerated list, are used to indicate
	// when the device intends to generate a DevicePropChanged event to let the opposing device know a threshold has
	// been reached, and therefore should be conservative (example 10). The value 0 may be realized in situations where
	// the device has alternate power provided by the ip or some other means.
	DPC_BatteryLevel DevicePropCode = 0x5001
	// DPC_FunctionalMode allows the functional mode of the device to be controlled. All devices are assumed to default
	// to a "standard mode." Alternate modes are typically used to indicate support for a reduced mode of operation
	// (e.g. sleep state) or an advanced mode or add-on that offers extended capabilities. The definition of
	// non-standard modes is device dependent. Any change in capability caused by a change in FunctionalMode shall be
	// evident by the DeviceInfoChanged event that is required to be sent by a device if its capabilities can change.
	// This property is described using the Enumeration form of the DevicePropDesc dataset. This property is also
	// exposed outside of sessions in the corresponding field in the DeviceInfo dataset.
	DPC_FunctionalMode DevicePropCode = 0x5002
	// DPC_ImageSize controls the height and width of the image that will be captured in pixels supported by the device.
	// This property takes the form of a Unicode, null terminated string that is parsed as follows: "WxH" where the W
	// represents the width and the H represents the height interpreted as unsigned integers. Example: width = 800,
	// height = 600, ImageSize string = "800x600" with a null-terminator on the end. This property may be expressed as
	// an enumerated list of allowed combinations, or if the individual width and height are linearly settable and
	// orthogonal to each other, they may be expressed as a range. For example, for a device that could set width from 1
	// to 640 and height from 1 to 480, the minimum in the range field would be "1x1" (nul lterminated), for a one-pixel
	// image, and the maximum would be "640x480" (null terminated), for the largest possible image. In this example, the
	// step would be "1x1" (null-terminated), indicating that the width and height are each incrementable to the
	// integer.
	// Changing this device property often causes fields in StorageInfo datasets to change, such as FreeSpaceInImages.
	// If this occurs, the device is required to issue a StorageInfoChanged event immediately after this property is
	// changed.
	DPC_ImageSize DevicePropCode = 0x5003
	// DPC_CompressionSetting is a property intended to be as close as is possible to being linear with respect to
	// perceived image quality over a broad range of scene content, and is represented by either a range or an
	// enumeration of integers. Low integers are used to represent low quality (i.e. maximum compression) while high
	// integers are used to represent high quality (i.e. minimum compression). No attempt is made in this standard to
	// assign specific values of this property with any absolute benchmark, so any available settings on a device are
	// relative to that device only and are therefore device-specific.
	DPC_CompressionSetting DevicePropCode = 0x5004
	// DPC_WhiteBalance is used to set how the device weights color channels. The device enumerates its supported values
	// for this property.
	DPC_WhiteBalance DevicePropCode = 0x5005
	// DPC_RGBGain takes the form of a Unicode, null-terminated string that is parsed as follows: "R:G:B" where the R
	// represents the red gain, the G represents the green gain, and the B represents the blue gain. For example, for an
	// RGB ratio of (red=4, green=2, blue=3), RGB string could be "4:2:3" (null-terminated) or "2000:1000:1500"
	// (null-terminated). The string parser for this property value should be able to support up to UINT16 integers for
	// R, G, and B. These values are relative to each other, and therefore may take on any integer value. This property
	// may be supported as an enumerated list of settings, or using a range. The minimum value would represent the
	// smallest numerical value (typically "1:1:1" null terminated). Using values of zero for a particular color channel
	// would mean that color channel would be dropped, so a value of "0:0:0" would result in images with all pixel
	// values being equal to zero. The maximum value would represent the largest value each field may be set to (up to
	// "65535:65535:65535" null-terminated), effectively determining the setting's granularity by an order of magnitude
	// per significant digit. The step value is typically "1:1:1". If a particular implementation desires the capability
	// to enforce minimum and/or maximum ratios, the green channel may be forced to a fixed value. An example of this
	// would be a minimum field of "1:1000:1", a maximum field of "20000:1000:20000" and a step field of "1:0:1".
	DPC_RGBGain DevicePropCode = 0x5006
	// DPC_FNumber allows the exposure program mode settings of the device, corresponding to the "Exposure Program" tag
	// within an EXIF or a TIFF/EP image file, to be constrained by a list of allowed exposure program mode settings
	// supported by the device.
	DPC_FNumber DevicePropCode = 0x5007
	// DPC_FocalLength represents the 35mm equivalent focal length. The values of this property correspond to the focal
	// length in millimeters multiplied by 100.
	DPC_FocalLength DevicePropCode = 0x5008
	// DPC_FocusDistance has unsigned integers as values corresponding to millimeters. A value of 0xFFFF corresponds to
	// a setting greater than 655 meters.
	DPC_FocusDistance DevicePropCode = 0x5009
	// DPC_FocusMode: the device enumerates the supported values of this property.
	DPC_FocusMode DevicePropCode = 0x500A
	// DPC_ExposureMeteringMode: the device enumerates the supported values of this property.
	DPC_ExposureMeteringMode DevicePropCode = 0x500B
	// DPC_FlashMode: the device enumerates the supported values of this property.
	DPC_FlashMode DevicePropCode = 0x500C
	// DPC_ExposureTime corresponds to the shutter speed. It has units of seconds scaled by 10,000. When the device is
	// in an automatic Exposure Program Mode, the setting of this property via SetDeviceProp may cause other properties
	// to change. Like all properties that cause other properties to change, the device is required to issue
	// DevicePropChanged events for the other properties that changed as the result of the initial change. This property
	// is typically only used by the device when the ProgramExposureMode is set to Manual or Shutter Priority.
	DPC_ExposureTime DevicePropCode = 0x500D
	// DPC_ExposureProgramMode allows the exposure program mode settings of the device, corresponding to the "Exposure
	// Program" tag within an EXIF or a TIFF/EP image file, to be constrained by a list of allowed exposure program mode
	// settings supported by the device.
	DPC_ExposureProgramMode DevicePropCode = 0x500E
	// DPC_ExposureIndex allows for the emulation of film speed settings on a Digital Camera. The settings correspond to
	// the ISO designations (ASA/DIN). Typically, a device supports discrete enumerated values but continuous control
	// over a range is possible. A value of 0xFFFF corresponds to Automatic ISO setting.
	DPC_ExposureIndex DevicePropCode = 0x500F
	// DPC_ExposureBiasCompensation allows for the adjustment of the set point of the digital camera's auto exposure
	// control. For example, a setting of 0 will not change the factory set auto exposure level. The units are in
	// "stops" scaled by a factor of 1000, in order to allow for fractional stop values. A setting of 2000 corresponds
	// to 2 stops more exposure (4X more energy on the sensor) yielding brighter images. A setting of -1000 corresponds
	// to one stop less exposure (1/2x the energy on the sensor) yielding darker images. The setting values are in APEX
	// units (Additive system of Photographic Exposure). This property may be expressed as an enumerated list or as a
	// range. This property is typically only used when the device has an ExposureProgramMode of Manual.
	DPC_ExposureBiasCompensation DevicePropCode = 0x5010
	// DPC_DateTime allows the current device date/time to be read and set. Date and time are represented in ISO
	// standard format as described in ISO 8601 from the most significant number to the least significant number. This
	// shall take the form of a Unicode string in the format "YYYYMMDDThhmmss.s" where YYYY is the year, MM is the month
	// 01-12, DD is the day of the month 01-31, T is a constant character, hh is the hours since midnight 00-23, mm is
	// the minutes 00-59 past the hour, and ss.s is the seconds past the minute, with the ".s" being optional tenths of
	// a second past the second. This string can optionally be appended with Z to indicate UTC, or +/-hhmm to indicate
	// the time is relative to a time zone. Appending neither indicates the time zone is unknown.
	// This property does not need to use a range or an enumeration, as the possible allowed time values are implicitly
	// specified by the definition of standard time and the format given in this and the ISO 8601 specifications.
	DPC_DateTime DevicePropCode = 0x5011
	// DPC_CaptureDelay describes the amount of time delay that should be inserted between the capture trigger and the
	// actual initiation of the data capture. This value shall be interpreted as milliseconds. This property is not
	// intended to be used to describe the time between frames for single-initiation multiple captures such as burst or
	// timelapse, which have separate interval properties. In those cases it would still serve as an initial delay
	// before the first image in the series was captured, independent of the time between frames. For no pre-capture
	// delay, this property should be set to zero.
	DPC_CaptureDelay DevicePropCode = 0x5012
	// DPC_StillCaptureMode allows for the specification of the type of still capture that is performed upon a still
	// capture initiation.
	DPC_StillCaptureMode DevicePropCode = 0x5013
	// DPC_Contrast controls the perceived contrast of captured images. This property may use an enumeration or range.
	// The minimum supported value is used to represent the least contrast, while the maximum value represents the most
	// contrast. Typically a value in the middle of the range would represent normal (default) contrast.
	DPC_Contrast DevicePropCode = 0x5014
	// DPC_Sharpness controls the perceived sharpness of captured images. This property may use an enumeration or range.
	// The minimum value is used to represent the least amount of sharpness, while the maximum value represents maximum
	// sharpness. Typically a value in the middle of the range would represent normal (default) sharpness.
	DPC_Sharpness DevicePropCode = 0x5015
	// DPC_DigitalZoom controls the effective zoom ratio of digital camera's acquired image scaled by a factor of 10. No
	// digital zoom (1X) corresponds to a value of 10, which is the standard scene size captured by the camera. A value
	// of 20 corresponds to a 2X zoom where 1/4 of the standard scene size is captured by the camera. This property may
	// be represented by an enumeration or a range. The minimum value should represent the minimum digital zoom
	// (typically 10), while the maximum value should represent the maximum digital zoom that the device allows.
	DPC_DigitalZoom DevicePropCode = 0x5016
	// DPC_EffectMode addresses special image acquisition modes of the camera.
	DPC_EffectMode DevicePropCode = 0x5017
	// DPC_BurstNumber controls the number of images that the device will attempt to capture upon initiation of a burst
	// operation.
	DPC_BurstNumber DevicePropCode = 0x5018
	// DPC_BurstInterval controls the time delay between captures upon initiation of a burst operation. This value is
	// expressed in whole milliseconds.
	DPC_BurstInterval DevicePropCode = 0x5019
	// DPC_TimelapseNumber controls the number of images that the device will attempt to capture upon initiation of a
	// time-lapse capture.
	DPC_TimelapseNumber DevicePropCode = 0x501a
	// DPC_TimelapseInterval controls the time delay between captures upon initiation of a time-lapse capture operation.
	// This value is expressed in milliseconds.
	DPC_TimelapseInterval DevicePropCode = 0x501b
	// DPC_FocusMeteringMode controls which automatic focus mechanism is used by the device. The device enumerates the
	// supported values of this property.
	DPC_FocusMeteringMode DevicePropCode = 0x501c
	// DPC_UploadURL is used to describe a standard Internet URL (Universal Resource Locator) that the receiving device
	// may use to upload images or objects to once they are acquired from the device.
	DPC_UploadURL DevicePropCode = 0x501d
	// DPC_Artist is used to contain the name of the owner/user of the device. This property is intended for use by the
	// device to populate the Artist field in any EXIF images that are captured with the device.
	DPC_Artist DevicePropCode = 0x501e
	// DPC_CopyrightInfo is used to contain the copyright notification. This property is intended for use by the device
	// to populate the Copyright field in any EXIF images that are captured with the device.
	DPC_CopyrightInfo DevicePropCode = 0x501F

	// DPD_Get indicates a read-only property.
	DPD_Get DevicePropDescCode = 0x00
	// DPD_GetSet indicates a read-write property.
	DPD_GetSet DevicePropDescCode = 0x01

	// DPF_FormFlag_None is for properties like DateTime. In this case the FORM field is not present.
	DPF_FormFlag_None DevicePropFormFlag = 0x00
	// DPF_FormFlag_Range indicates a range form.
	DPF_FormFlag_Range DevicePropFormFlag = 0x01
	// DPF_FormFlag_Enum indicates an enumeration form.
	DPF_FormFlag_Enum DevicePropFormFlag = 0x02

	// Undefined
	DTC_UNDEF DataTypeCode = 0x0000
	// Signed 8 bit integer
	DTC_INT8 DataTypeCode = 0x0001
	// Unsigned 8 bit integer
	DTC_UINT8 DataTypeCode = 0x0002
	// Signed 16 bit integer
	DTC_INT16 DataTypeCode = 0x0003
	// Unsigned 16 bit integer
	DTC_UINT16 DataTypeCode = 0x0004
	// Signed 32 bit integer
	DTC_INT32 DataTypeCode = 0x0005
	// Unsigned 32 bit integer
	DTC_UINT32 DataTypeCode = 0x0006
	// Signed 64 bit integer
	DTC_INT64 DataTypeCode = 0x0007
	// Unsigned 64 bit integer
	DTC_UINT64 DataTypeCode = 0x0008
	// Signed 128 bit integer
	DTC_INT128 DataTypeCode = 0x0009
	// Unsigned 128 bit integer
	DTC_UINT128 DataTypeCode = 0x000A
	// Array of Signed 8 bit integers
	DTC_AINT8 DataTypeCode = 0x4001
	// Array of Unsigned 8 bit integers
	DTC_AUINT8 DataTypeCode = 0x4002
	// Array of Signed 16 bit integers
	DTC_AINT16 DataTypeCode = 0x4003
	// Array of Unsigned 16 bit integers
	DTC_AUINT16 DataTypeCode = 0x4004
	// Array of Signed 32 bit integers
	DTC_AINT32 DataTypeCode = 0x4005
	// Array of Unsigned 32 bit integers
	DTC_AUINT32 DataTypeCode = 0x4006
	// Array of Signed 64 bit integers
	DTC_AINT64 DataTypeCode = 0x4007
	// Array of Unsigned 64 bit integers
	DTC_AUINT64 DataTypeCode = 0x4008
	// Array of Signed 128 bit integers
	DTC_AINT128 DataTypeCode = 0x4009
	// Array of Unsigned 128 bit integers
	DTC_AUINT128 DataTypeCode = 0x400A
	// Variable-length Unicode String
	DTC_STR DataTypeCode = 0xFFFF
)

// DevicePropCodeAsString returns the DevicePropCode as string. When the DevicePropCode is unknown, it returns an empty
// string.
func DevicePropCodeAsString(code DevicePropCode) string {
	switch code {
	case DPC_BatteryLevel:
		return "battery level"
	case DPC_FunctionalMode:
		return "functional mode"
	case DPC_ImageSize:
		return "image size"
	case DPC_CompressionSetting:
		return "compression setting"
	case DPC_WhiteBalance:
		return "white balance"
	case DPC_RGBGain:
		return "RGB gain"
	case DPC_FNumber:
		return "F number"
	case DPC_FocalLength:
		return "focal length"
	case DPC_FocusDistance:
		return "focus distance"
	case DPC_FocusMode:
		return "focus mode"
	case DPC_ExposureMeteringMode:
		return "exposure metering mode"
	case DPC_FlashMode:
		return "flash mode"
	case DPC_ExposureTime:
		return "exposure time"
	case DPC_ExposureProgramMode:
		return "exposure program mode"
	case DPC_ExposureIndex:
		return "exposure index"
	case DPC_ExposureBiasCompensation:
		return "exposure bias compensation"
	case DPC_DateTime:
		return "date time"
	case DPC_CaptureDelay:
		return "capture delay"
	case DPC_StillCaptureMode:
		return "still capture mode"
	case DPC_Contrast:
		return "contrast"
	case DPC_Sharpness:
		return "sharpness"
	case DPC_DigitalZoom:
		return "digital zoom"
	case DPC_EffectMode:
		return "effect mode"
	case DPC_BurstNumber:
		return "burst number"
	case DPC_BurstInterval:
		return "burst interval"
	case DPC_TimelapseNumber:
		return "timelapse number"
	case DPC_TimelapseInterval:
		return "timelapse interval"
	case DPC_FocusMeteringMode:
		return "focus metering mode"
	case DPC_UploadURL:
		return "upload URL"
	case DPC_Artist:
		return "artist"
	case DPC_CopyrightInfo:
		return "copyright info"
	default:
		return ""
	}
}

type DevicePropDesc struct {
	// DevicePropertyCode is a specific DevicePropCode
	DevicePropertyCode DevicePropCode
	// DataType identifies the DatatypeCode of the property
	DataType DataTypeCode
	// GetSet indicates whether the property is read-only (Get) or read-write (Get/Set).
	GetSet DevicePropDescCode
	// FactoryDefaultValue identifies the value of the factory default setting for the property.
	FactoryDefaultValue []byte
	// CurrentValue identifies the current value of the property.
	CurrentValue []byte
	// FormFlag indicates the format of the next field.
	FormFlag DevicePropFormFlag
	// Form is the Enumeration-Form or the Range-Form, or is absent if FormFlag = 0.
	Form interface{}
}

func (dpd *DevicePropDesc) SizeOfValueInBytes() int {
	switch dpd.DataType {
	case DTC_INT8, DTC_UINT8:
		return 1
	case DTC_INT16, DTC_UINT16:
		return 2
	case DTC_INT32, DTC_UINT32:
		return 4
	case DTC_INT64, DTC_UINT64:
		return 8
	default:
		return 0
	}
}

func (dpd *DevicePropDesc) CurrentValueAsInt64() int64 {
	v := dpd.CurrentValue
	if dpd.SizeOfValueInBytes() < 8 {
		pad := make([]byte, 8-dpd.SizeOfValueInBytes())
		v = append(v, pad...)
	}

	// Converting between uint64 and int64 does not change the sign bit, only the way it is interpreted.
	return int64(binary.LittleEndian.Uint64(v))
}

type RangeForm struct {
	// MinimumValue is the minimum value of the property supported by the device.
	MinimumValue []byte
	// MaximumValue is the maximum value of the property supported by the device.
	MaximumValue []byte
	// StepSize indicates the size of steps that can be used to increment or decrement a value.
	// A particular vendor's device shall support all values of a property defined by MinimumValue + N x StepSize which
	// is less than or equal to MaximumValue where N=0 to a vendor defined maximum.
	StepSize []byte
}

type EnumerationForm struct {
	// NumberOfValues indicates the number of values of size DTS of the particular property supported by the device.
	NumberOfValues  int
	// SupportedValues holds the list of supported values.
	SupportedValues [][]byte
}

// DeviceInfo is used to hold the description information for a device. The Initiator can obtain this dataset from the
// Responder without opening a session with the device. This dataset holds data that describes the device and its
// capabilities. This information is only static if the device capabilities cannot change during a session, which would
// be indicated by a change in the FunctionalMode value in the dataset. For example, if the device goes into a sleep
// mode in which it can still respond to GetDeviceInfo requests, the data in this dataset should reflect the
// capabilities of the device while it is in that mode only (including any operations and properties needed to change
// the FunctionalMode, if this is allowed remotely). If the power state or the capabilities of the device changes (due
// to a FunctionalMode change), a DeviceInfoChanged event shall be issued to all sessions in order to indicate how its
// capabilities have changed.
type DeviceInfo struct {
	// StandardVersion is the highest version of the standard that the device can support. This represents the standard
	// version expressed in hundredths (e.g. 1.32 would be stored as 132).
	StandardVersion uint16

	// VendorExtensionID provides the context for interpretation of any vendor extensions used by this device. If no
	// extensions are supported, this field shall be set to 0x00000000. If vendor-specific codes of any type are used,
	// this field is mandatory, and should not be set to 0x00000000. These IDs are assigned by PIMA.
	VendorExtensionID uint32

	// VendorExtensionVersion is the vendor-specific version number of extensions that are supported. This shall be
	// expressed in hundredths (e.g. 1.32 would be stored as 132).
	VendorExtensionVersion uint16

	// VendorExtensionDesc is an optional string used to hold a human-readable description of the VendorExtensionID.
	// This field should only be used for informational purposes, and not as the context for the interpretation of
	// vendor-extensions.
	VendorExtensionDesc string

	// FunctionalMode is an optional field used to hold the functional mode. This field controls whether the device is
	// in an alternate mode that provides a different set of capabilities (i.e. supported operations, events, etc.) If
	// the device only supports one mode, this value should always be zero.
	// The functional mode information is held by the device as a device property. In order to change the functional
	// mode of the device remotely, a session needs to be opened with the device, and the SetDeviceProp operation needs
	// to be used.
	FunctionalMode FunctionalMode

	// OperationsSupported is an array of OperationCodes representing operations that the device is currently
	// supporting, given the FunctionalMode indicated.
	OperationsSupported []OperationCode

	// EventsSupported is an array of EventCodes representing the events that are currently generated by the device in
	// appropriate situations, given the FunctionalMode indicated.
	EventsSupported []EventCode

	// DevicePropertiesSupported is an array of DevicePropCodes representing DeviceProperties that are currently exposed
	// for reading and/or modification, given the FunctionalMode indicated.
	DevicePropertiesSupported []DevicePropCode

	// CaptureFormats is the list of data formats in ObjectFormatCode form that the device can create using an
	// InitiateCapture operation and/or an InitiateOpenCapture operation, given the FunctionalMode indicated. These are
	// typically image object formats, but can include any object format that can be fully captured using a single
	// trigger mechanism, or an initiate/terminate mechanism. All image object formats that a device can capture data in
	// shall be listed prior to any non-image object formats, and shall be in preferential order such that the default
	// capture format is first.
	CaptureFormats []ObjectFormatCode

	// ImageFormats is the list of image formats in ObjectFormatCode form that the device supports in order of highest
	// preference to lowest preference. Support for an image format refers to the ability to interpret image file
	// contents according to that format's specifications, for display and/or manipulation purposes. For image output
	// devices, this field represents the image formats that the output device is capable of outputting. This field does
	// not describe any device format-translation capabilities.
	ImageFormats []ObjectFormatCode

	// Manufacturer is an optional human-readable string used to hold the Responder's manufacturer.
	Manufacturer string

	// Model is an optional human-readable string used to communicate the Responder's model name.
	Model string

	// DeviceVersion is an optional string used to communicate the Responder's firmware or software version in a
	// vendor-specific way.
	DeviceVersion string

	// SerialNumber an optional string used to communicate the Responder's serial number, which is defined as a unique
	// value among all devices sharing identical Model and Device Version fields. If unique serial numbers are not
	// supported, this field shall be set to the empty string. The presence of a non-null string in the SerialNumber
	// field for one device infers that this field is non-zero and unique among all devices of that model and version.
	SerialNumber string
}
