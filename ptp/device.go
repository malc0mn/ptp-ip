package ptp

type DataTypeCode uint16

const(
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

type DevicePropDesc struct {
	// A specific DevicePropCode
	DevicePropertyCode DevicePropCode
	// This field identifies the DatatypeCode of the property
	DataType DataTypeCode
	// This field indicates whether the property is read-only (Get) or read-write (Get/Set).
	GetSet DevicePropDescCode
	// This field identifies the value of the factory default setting for the property.
	FactoryDefaultValue interface{}
	// This field identifies the current value of the property.
	CurrentValue interface{}
	// This field indicates the format of the next field.
	FormFlag DevicePropFormFlag
	// This dataset is the Enumeration-Form or the Range-Form, or is absent if Form Flag = 0
	Form RangeForm
}

type RangeForm struct {
	// Minimum value of property supported by the device.
	MinimumValue interface{}
	// Maximum value of property supported by the device.
	MaximumValue interface{}
	// A particular vendor's device shall support all values of a property defined by MinimumValue + N x StepSize which
	// is less than or equal to MaximumValue where N=0 to a vendor defined maximum
	StepSize interface{}
}

type EnumerationForm struct {
	// This field indicates the number of values of size DTS of the particular property supported by the device.
	NumberOfValues int
	SupportedValues []interface{}
}

// This dataset is used to hold the description information for a device. The Initiator can obtain this dataset from the
// Responder without opening a session with the device. This dataset holds data that describes the device and its
// capabilities. This information is only static if the device capabilities cannot change during a session, which would
// be indicated by a change in the FunctionalMode value in the dataset. For example, if the device goes into a sleep
// mode in which it can still respond to GetDeviceInfo requests, the data in this dataset should reflect the
// capabilities of the device while it is in that mode only (including any operations and properties needed to change
// the FunctionalMode, if this is allowed remotely). If the power state or the capabilities of the device changes (due
// to a FunctionalMode change), a DeviceInfoChanged event shall be issued to all sessions in order to indicate how its
// capabilities have changed.
type DeviceInfo struct {
	// Highest version of the standard that the device can support. This represents the standard version expressed in
	// hundredths (e.g. 1.32 would be stored as 132).
	StandardVersion uint16

	// Provides the context for interpretation of any vendor extensions used by this device. If no extensions are
	// supported, this field shall be set to 0x00000000. If vendor-specific codes of any type are used, this field is
	// mandatory, and should not be set to 0x00000000. These IDs are assigned by PIMA, as described in Clause 9.5.
	VendorExtensionID uint32

	// The vendor-specific version number of extensions that are supported. This shall be expressed in hundredths (e.g.
	// 1.32 would be stored as 132).
	VendorExtensionVersion uint16

	// An optional string used to hold a human-readable description of the VendorExtensionID. This field should only be
	// used for informational purposes, and not as the context for the interpretation of vendor-extensions.
	VendorExtensionDesc string

	// An optional field used to hold the functional mode. This field controls whether the device is in an alternate
	// mode that provides a different set of capabilities (i.e. supported operations, events, etc.) If the device only
	// supports one mode, this value should always be zero.
	// The functional mode information is held by the device as a device property. In order to change the functional
	// mode of the device remotely, a session needs to be opened with the device, and the SetDeviceProp operation needs
	// to be used.
	FunctionalMode FunctionalMode

	// This field is an array of OperationCodes representing operations that the device is currently supporting, given
	// the FunctionalMode indicated.
	OperationsSupported []OperationCode

	// This field is an array of EventCodes representing the events that are currently generated by the device in
	// appropriate situations, given the FunctionalMode indicated.
	EventsSupported []EventCode

	// This field is an array of DevicePropCodes representing DeviceProperties that are currently exposed for reading
	// and/or modification, given the FunctionalMode indicated.
	DevicePropertiesSupported []DevicePropCode

	// The list of data formats in ObjectFormatCode form that the device can create using an InitiateCapture operation
	// and/or an InitiateOpenCapture operation, given the FunctionalMode indicated. These are typically image object
	// formats, but can include any object format that can be fully captured using a single trigger mechanism, or an
	// initiate/terminate mechanism. All image object formats that a device can capture data in shall be listed prior to
	// any non-image object formats, and shall be in preferential order such that the default capture format is first.
	CaptureFormats []ObjectFormatCode

	// The list of image formats in ObjectFormatCode form that the device supports in order of highest preference to
	// lowest preference. Support for an image format refers to the ability to interpret image file contents according
	// to that format's specifications, for display and/or manipulation purposes. For image output devices, this field
	// represents the image formats that the output device is capable of outputting. This field does not describe any
	// device format-translation capabilities.
	ImageFormats []ObjectFormatCode

	// An optional human-readable string used to hold the Responder's manufacturer.
	Manufacturer string

	// An optional human-readable string used to communicate the Responder's model name.
	Model string

	// An optional string used to communicate the Responder's firmware or software version in a vendor-specific way.
	DeviceVersion string

	// An optional string used to communicate the Responder's serial number, which is defined as a unique value among
	// all devices sharing identical Model and Device Version fields. If unique serial numbers are not supported, this
	// field shall be set to the empty string. The presence of a non-null string in the SerialNumber field for one
	// device infers that this field is non-zero and unique among all devices of that model and version.
	SerialNumber string
}
