package ptp

type DevicePropCode uint16
type DevicePropDescCode uint16

const (
	DPC_Undefined DevicePropCode = 0x5000
	// Battery level is a read-only property typically represented by a range of integers. The minimum field should be
	// set to the integer used for no power (example 0), and the maximum should be set to the integer used for full
	// power (example 100). The step field, or the individual thresholds in an enumerated list, are used to indicate
	// when the device intends to generate a DevicePropChanged event to let the opposing device know a threshold has
	// been reached, and therefore should be conservative (example 10). The value 0 may be realized in situations where
	// the device has alternate power provided by the transport or some other means.
	DPC_BatteryLevel DevicePropCode = 0x5001
	// Allows the functional mode of the device to be controlled. All devices are assumed to default to a "standard
	// mode." Alternate modes are typically used to indicate support for a reduced mode of operation (e.g. sleep state)
	// or an advanced mode or add-on that offers extended capabilities. The definition of non-standard modes is
	// devicedependent. Any change in capability caused by a change in FunctionalMode shall be evident by the
	// DeviceInfoChanged event that is required to be sent by a device if its capabilities can change. This property is
	// described using the Enumeration form of the DevicePropDesc dataset. This property is also exposed outside of
	// sessions in the corresponding field in the DeviceInfo dataset.
	DPC_FunctionalMode DevicePropCode = 0x5002
	// This property controls the height and width of the image that will be captured in pixels supported by the device.
	// This property takes the form of a Unicode, nullterminated string that is parsed as follows: "WxH" where the W
	// represents the width and the H represents the height interpreted as unsigned integers. Example: width = 800,
	// height = 600, ImageSize string = "800x600" with a null-terminator on the end. This property may be expressed as
	// an enumerated list of allowed combinations, or if the individual width and height are linearly settable and
	// orthogonal to each other, they may be expressed as a range. For example, for a device that could set width from 1
	// to 640 and height from 1 to 480, the minimum in the range field would be "1x1" (nullterminated), for a one-pixel
	// image, and the maximum would be "640x480" (nullterminated), for the largest possible image. In this example, the
	// step would be "1x1" (null-terminated), indicating that the width and height are each incrementable to the
	// integer.
	// Changing this device property often causes fields in StorageInfo datasets to change, such as FreeSpaceInImages.
	// If this occurs, the device is required to issue a StorageInfoChanged event immediately after this property is
	// changed.
	DPC_ImageSize DevicePropCode = 0x5003
	// Compression setting is a property intended to be as close as is possible to being linear with respect to
	// perceived image quality over a broad range of scene content, and is represented by either a range or an
	// enumeration of integers. Low integers are used to represent low quality (i.e. maximum compression) while high
	// integers are used to represent high quality (i.e. minimum compression). No attempt is made in this standard to
	// assign specific values of this property with any absolute benchmark, so any available settings on a device are
	// relative to that device only and are therefore device-specific.
	DPC_CompressionSetting DevicePropCode = 0x5004
	// This property is used to set how the device weights color channels. The device enumerates its supported values
	// for this property.
	DPC_WiteBalance DevicePropCode = 0x5005
	// This property takes the form of a Unicode, null-terminated string that is parsed as follows: "R:G:B" where the R
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
	// This property allows the exposure program mode settings of the device, corresponding to the "Exposure Program"
	// tag within an EXIF or a TIFF/EP image file, to be constrained by a list of allowed exposure program mode settings
	// supported by the device.
	DPC_FNumber DevicePropCode = 0x5007
	// This property represents the 35mm equivalent focal length. The values of this property correspond to the focal
	// length in millimeters multiplied by 100.
	DPC_FocalLength DevicePropCode = 0x5008
	// The values of this property are unsigned integers with the values corresponding to millimeters. A value of 0xFFFF
	// corresponds to a setting greater than 655 meters.
	DPC_FocusDistance DevicePropCode = 0x5009
	// The device enumerates the supported values of this property.
	DPC_FocusMode DevicePropCode = 0x500a
	// The device enumerates the supported values of this property.
	DPC_ExposureMeteringMode DevicePropCode = 0x500b
	// The device enumerates the supported values of this property.
	DPC_FlashMode DevicePropCode = 0x500c
	// This property corresponds to the shutter speed. It has units of seconds scaled by 10,000. When the device is in
	// an automatic Exposure Program Mode, the setting of this property via SetDeviceProp may cause other properties to
	// change. Like all properties that cause other properties to change, the device is required to issue
	// DevicePropChanged events for the other properties that changed as the result of the initial change. This property
	// is typically only used by the device when the ProgramExposureMode is set to Manual or Shutter Priority.
	DPC_ExposureTime DevicePropCode = 0x500d
	// This property allows the exposure program mode settings of the device, corresponding to the "Exposure Program"
	// tag within an EXIF or a TIFF/EP image file, to be constrained by a list of allowed exposure program mode settings
	// supported by the device.
	DPC_ExposureProgramMode DevicePropCode = 0x500e
	// This property allows for the emulation of film speed settings on a Digital Camera. The settings correspond to the
	// ISO designations (ASA/DIN). Typically, a device supports discrete enumerated values but continuous control over a
	// range is possible. A value of 0xFFFF corresponds to Automatic ISO setting.
	DPC_ExposureIndex DevicePropCode = 0x500f
	// This property allows for the adjustment of the set point of the digital camera's auto exposure control. For
	// example, a setting of 0 will not change the factory set auto exposure level. The units are in "stops" scaled by a
	// factor of 1000, in order to allow for fractional stop values. A setting of 2000 corresponds to 2 stops more
	// exposure (4X more energy on the sensor) yielding brighter images. A setting of -1000 corresponds to one stop less
	// exposure (1/2x the energy on the sensor) yielding darker images. The setting values are in APEX units (Additive
	// system of Photographic Exposure). This property may be expressed as an enumerated list or as a range. This
	// property is typically only used when the device has an ExposureProgramMode of Manual.
	DPC_ExposureBiasCompensation DevicePropCode = 0x5010
	// This property allows the current device date/time to be read and set. Date and time are represented in ISO
	// standard format as described in ISO 8601 from the most significant number to the least significant number. This
	// shall take the form of a Unicode string in the format "YYYYMMDDThhmmss.s" where YYYY is the year, MM is the month
	// 01-12, DD is the day of the month 01-31, T is a constant character, hh is the hours since midnight 00-23, mm is
	// the minutes 00-59 past the hour, and ss.s is the seconds past the minute, with the ".s" being optional tenths of
	// a second past the second. This string can optionally be appended with Z to indicate UTC, or +/-hhmm to indicate
	// the time is relative to a time zone. Appending neither indicates the time zone is unknown.
	// This property does not need to use a range or an enumeration, as the possible allowed time values are implicitly
	// specified by the definition of standard time and the format given in this and the ISO 8601 specifications.
	DPC_DateTime DevicePropCode = 0x5011
	// This value describes the amount of time delay that should be inserted between the capture trigger and the actual
	// initiation of the data capture. This value shall be interpreted as milliseconds. This property is not intended to
	// be used to describe the time between frames for single-initiation multiple captures such as burst or time-lapse,
	// which have separate interval properties outlined in Clauses 13.4.25 and 13.4.27. In those cases it would still
	// serve as an initial delay before the first image in the series was captured, independent of the time between
	// frames. For no pre-capture delay, this property should be set to zero.
	DPC_CaptureDelay DevicePropCode = 0x5012
	// This property allows for the specification of the type of still capture that is performed upon a still capture
	// initiation.
	DPC_StillCaptureMode DevicePropCode = 0x5013
	// This property controls the perceived contrast of captured images. This property may use an enumeration or range.
	// The minimum supported value is used to represent the least contrast, while the maximum value represents the most
	// contrast. Typically a value in the middle of the range would represent normal (default) contrast.
	DPC_Contrast DevicePropCode = 0x5014
	// This property controls the perceived sharpness of captured images. This property may use an enumeration or range.
	// The minimum value is used to represent the least amount of sharpness, while the maximum value represents maximum
	// sharpness. Typically a value in the middle of the range would represent normal (default) sharpness.
	DPC_Sharpness DevicePropCode = 0x5015
	// This property controls the effective zoom ratio of digital camera's acquired image scaled by a factor of 10. No
	// digital zoom (1X) corresponds to a value of 10, which is the standard scene size captured by the camera. A value
	// of 20 corresponds to a 2X zoom where 1/4 of the standard scene size is captured by the camera. This property may
	// be represented by an enumeration or a range. The minimum value should represent the minimum digital zoom
	// (typically 10), while the maximum value should represent the maximum digital zoom that the device allows.
	DPC_DigitalZoom DevicePropCode = 0x5016
	// This property addresses special image acquisition modes of the camera.
	DPC_EffectMode DevicePropCode = 0x5017
	// This property controls the number of images that the device will attempt to capture upon initiation of a burst
	// operation.
	DPC_BurstNumber DevicePropCode = 0x5018
	// This property controls the time delay between captures upon initiation of a burst operation. This value is
	// expressed in whole milliseconds.
	DPC_BurstInterval DevicePropCode = 0x5019
	// This property controls the number of images that the device will attempt to capture upon initiation of a
	// time-lapse capture.
	DPC_TimelapseNumber DevicePropCode = 0x501a
	// This property controls the time delay between captures upon initiation of a time-lapse capture operation. This
	// value is expressed in milliseconds.
	DPC_TimelapseInterval DevicePropCode = 0x501b
	// This property controls which automatic focus mechanism is used by the device. The device enumerates the supported
	// values of this property.
	DPC_FocusMeteringMode DevicePropCode = 0x501c
	// This property is used to describe a standard Internet URL (Universal Resource Locator) that the receiving device
	// may use to upload images or objects to once they are acquired from the device.
	DPC_UploadURL DevicePropCode = 0x501d
	// This property is used to contain the name of the owner/user of the device. This property is intended for use by
	// the device to populate the Artist field in any EXIF images that are captured with the device.
	DPC_Artist DevicePropCode = 0x501e
	// This property is used to contain the copyright notification. This property is intended for use by the device to
	// populate the Copyright field in any EXIF images that are captured with the device.
	DPC_CopyrightInfo DevicePropCode = 0x501F

	// Indicates a read-only property.
	DPD_Get DevicePropDescCode = 0x00
	// Indicates a read-write property.
	DPD_GetSet DevicePropDescCode = 0x01
	// This is for properties like DateTime. In this case the FORM field is not present.
	DPD_FormFlag_None DevicePropDescCode = 0x00
	// Range form
	DPD_FormFlag_Range DevicePropDescCode = 0x01
	// Enumeration form
	DPD_FormFlag_Enum DevicePropDescCode = 0x02
)
