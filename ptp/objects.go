package ptp

import "time"

type AssociationDesc uint16
type AssociationType uint16
type ObjectFormatCode uint16
type ObjectHandle uint32

const (
	AD_Undefined            AssociationDesc = 0x0000
	AD_Unused               AssociationDesc = 0x0001
	AD_Reserved             AssociationDesc = 0x0002
	AD_DefaultPlaybackDelta AssociationDesc = 0x0003
	AD_ImagesPerRow         AssociationDesc = 0x0006

	AT_Undefined AssociationType = 0x0000
	// This association type is used to represent a folder that may hold any type of object, and is analogous to the
	// standard folder present in most filesystems. This association is typically used to represent a local grouping of
	// objects, with no other relationship implied.
	AT_GenericFolder AssociationType = 0x0001
	// This association type is the same as a folder but is used to hold image and data objects that have logical
	// groupings according to content, capture sessions, or any other unspecified user-determined grouping. These are
	// typically created by a user or automation technique. Some devices may wish to expose albums to the user but not
	// all generic folders. Devices that do not distinguish between albums and folders should only use the
	// AssociationType of GenericFolder. The AssociationDesc field is reserved for future definition by this standard,
	// and for devices that support this version of the specification, shall be set to 0x00000000.
	AT_Album AssociationType = 0x0002
	// Indicates that the data objects are part of a sequence of data captures that make up a set of time-ordered data
	// of the same subject. This association is used to represent time-lapse or burst sequences. The order is
	// interpreted to be sequential by capture time from first captured to last, and is indicated by the increasing
	// values of the SequenceNumber fields in the ObjectInfo dataset for each object. If known, the AssociationDesc acts
	// as a DefaultPlaybackDelta, and should be set to the desired time in milliseconds to delay between each object if
	// exposing them sequentially in real-time. If unknown, this value should be set to 0x00000000.
	AT_TimeSequence AssociationType = 0x0003
	// Indicates that the associated data objects make up a panoramic series of images that are arranged side-by-side,
	// in a horizontal fashion. The order of the sequence, from left to right when facing the subject, is indicated by
	// the increasing values of the SequenceNumber fields in each object. The AssociationDesc is unused, and should be
	// set to 0x00000000. For example, four images would have SequenceNumbers assigned as follows: | 1 | 2 | 3 | 4 |
	AT_HorizontalPanoramic AssociationType = 0x0004
	// Indicates that the associated data objects make up a panoramic series of images that are arranged bottom-to-top,
	// in a vertical fashion. The order of the sequence, from bottom to top when facing the subject, is indicated by the
	// increasing values of the SequenceNumber fields in each object. The AssociationDesc is unused, and should be set
	// to 0x00000000. For example, four images would have SequenceNumbers assigned as follows:
	// | 4 |
	// +---+
	// | 3 |
	// +---+
	// | 2 |
	// +---+
	// | 1 |
	AT_VerticalPanoramic AssociationType = 0x0005
	// Indicates that the associated data objects make up a two-dimensional panoramic series of images that are arranged
	// left-to-right and bottom-to-top in adjacent or overlapping horizontal strips. The order of the sequence, from
	// bottom-left to top-right when facing the subject, is indicated by the increasing values of the SequenceNumber
	// fields in each object. The AssociationDesc is used to indicate the number of images in each row. For example,
	// sixteen images arranged in a 4x4 2Dpanoramic would have SequenceNumbers assigned as follows:
	// | 13 | 14 | 15 | 16 |
	// +----+----+----+----+
	// |  9 | 10 | 11 | 12 |
	// +----+----+----+----+
	// |  5 |  6 |  7 |  8 |
	// +----+----+----+----+
	// |  1 |  2 |  3 |  4 |
	AT_2DPanoramic AssociationType = 0x0006
	// Indicates that the association represents one or more non-image objects being associated with an image object.
	// For example, an image capture that also stores independent audio text files that are temporally related to the
	// image capture may use this type of association to indicate the relationship. Optionally, if the individual
	// objects are ordered (e.g. multiple temporally-related sound files), the SequenceNumber field of each object's
	// ObjectInfo dataset should contain increasing integers. If the individual objects are unordered, the
	// AssociationDesc should be set to 0x00000000. The AssociationDesc field is unused, and should be set to
	// 0x00000000.
	AT_AncillaryData AssociationType = 0x0007

	// Undefined non-image object
	OFC_Undefined ObjectFormatCode = 0x3000
	// Association (e.g. folder)
	OFC_Association ObjectFormatCode = 0x3001
	// Device-model-specific script
	OFC_Script ObjectFormatCode = 0x3002
	// Device-model-specific binary executable
	OFC_Executable ObjectFormatCode = 0x3003
	// Text file
	OFC_Text ObjectFormatCode = 0x3004
	// HyperText Markup Language file (text)
	OFC_HTML ObjectFormatCode = 0x3005
	// Digital Print Order Format file (text)
	OFC_DPOF ObjectFormatCode = 0x3006
	// Audio clip
	OFC_AIFF ObjectFormatCode = 0x3007
	// Audio clip
	OFC_WAV ObjectFormatCode = 0x3008
	// Audio clip
	OFC_MP3 ObjectFormatCode = 0x3009
	// Video clip
	OFC_AVI ObjectFormatCode = 0x300A
	// Video clip
	OFC_MPEG ObjectFormatCode = 0x300B
	// Microsoft Advanced Streaming Format (video)
	OFC_ASF ObjectFormatCode = 0x300C
	// Unknown image object
	OFC_Unknown ObjectFormatCode = 0x3800
	// Exchangeable File Format, JEIDA standard
	OFC_EXIF_JPEG ObjectFormatCode = 0x3801
	// Tag Image File Format for Electronic Photography
	OFC_TIFF_EP ObjectFormatCode = 0x3802
	// Structured Storage Image Format
	OFC_FlashPix ObjectFormatCode = 0x3803
	// Microsoft Windows Bitmap file
	OFC_BMP ObjectFormatCode = 0x3804
	// Canon Camera Image File Format
	OFC_CIFF ObjectFormatCode = 0x3805
	// Graphics Interchange Format
	OFC_GIF ObjectFormatCode = 0x3807
	// JPEG File Interchange Format
	OFC_JFIF ObjectFormatCode = 0x3808
	// PhotoCD Image Pac
	OFC_PCD ObjectFormatCode = 0x3809
	// Quickdraw Image Format
	OFC_PICT ObjectFormatCode = 0x380A
	// Portable Network Graphics
	OFC_PNG ObjectFormatCode = 0x380B
	// Tag Image File Format
	OFC_TIFF ObjectFormatCode = 0x380D
	// Tag Image File Format for Information Technology (graphic arts)
	OFC_TIFF_IT ObjectFormatCode = 0x380E
	// JPEG2000 Baseline File Format
	OFC_JP2 ObjectFormatCode = 0x380F
	// JPEG2000 Extended File Format
	OFC_JPX ObjectFormatCode = 0x3810
)

// This dataset is used to define the information about data objects in persistent store, as well as optional
// information if the data is known to be an image or an association object. It is required that these data items be
// accounted for in response to a GetObjectInfo operation. If the data is not known to be an image, or the image
// information is unavailable, the image-specific fields shall be set to zero. Objects of type Association are fully
// qualified by the ObjectInfo dataset.
type ObjectInfo struct {
	// The StorageID of the device's store in which the image resides.
	StorageID StorageID

	//  Indicates ObjectFormatCode of the object.
	ObjectFormat ObjectFormatCode

	// An optional field representing the write-protection status of the data object. Objects that are protected may not
	// be deleted as the result of any operations specified in this standard without first separately removing their
	// protection status in a separate transaction.
	// This protection field is distinctly different in scope than the AccessCapability field present in the StorageInfo
	// dataset. If an attempt to delete an object is made, success will only occur if the ProtectionStatus of the object
	// is 0x0000 and the AccessCapability of the store allows deletion. If a device does not support object protection,
	// this field should always be set to 0x0000, and the SetProtection operation should not be supported.
	ProtectionStatus ProtectionStatus

	// The size of the buffer needed to hold the entire binary object in bytes. This field may be used for memory
	// allocation purposes in object receivers by ip implementations.
	ObjectCompressedSize uint32

	// Indicates ObjectFormatCode of the thumbnail. In order for an object to be referred to as an image, it must be able to
	// produce a thumbnail as the response to a request. Therefore, this value should only be 0x00000000 for the case of
	// non-image objects.
	ThumbFormat ObjectFormatCode

	// The size of the buffer needed to hold the thumbnail. This field may be used for memory allocation purposes. In
	// order for an object to be referred to as an image, it must be able to produce a thumbnail as the response to a
	// request. Therefore, this value should only be 0x00000000 for the case of non-image objects.
	ThumbCompressedSize uint32

	// An optional field representing the width of the thumbnail in pixels. If this field is not supported or the object
	// is not an image, the value 0x00000000 shall be used.
	ThumbPixWidth uint32

	// An optional field representing the height of the thumbnail in pixels. If this field is not supported or the
	// object is not an image, the value 0x00000000 shall be used.
	ThumbPixHeight uint32

	// An optional field representing the width of the image in pixels. If the data is not known to be an image, this
	// field should be set to 0x00000000. The purpose of this field is to enable an application to provide the width
	// information to a user prior to transferring the image. If this field is not supported, the value 0x00000000 shall
	// be used.
	ImagePixWidth uint32

	// An optional field representing the height of the image in pixels. If the data is not known to be an image, this
	// field should be set to 0x00000000. The purpose of this field is to enable an application to provide the height
	// information to a user prior to transferring the image. If this field is not supported, the value 0x00000000 shall
	// be used.
	ImagePixHeight uint32

	// An optional field representing the total number of bits per pixel of the uncompressed image. If the data is not
	// known to be an image, this field should be set to 0x00000000. The purpose of this field is to enable an
	// application to provide the bit depth information to a user prior to transferring the image. This field does not
	// attempt to specify the number of bits assigned to particular color channels, but instead represents the total
	// number of bits used to describe one pixel. If this field is not supported, the value 0x00000000 shall be used.
	// This field should not be used for memory allocation purposes, but is strictly information that is typically
	// inside of an image object, that may affect whether or not a user wishes to transfer the image, and therefore is
	// exposed prior to object transfer in the ObjectInfo dataset.
	ImageBitDepth uint32

	// Indicates the handle of the object that is the parent of this object. The ParentObject must be of object type
	// Association. If the device does not support associations, or the object is in the “root” of the hierarchical
	// store, then this value should be set to 0x00000000.
	ParentObject ObjectHandle

	// A field that is only used for objects of type Association. This code indicates the type of association. If the
	// object is not an association, this field should be set to 0x0000.
	AssociationType AssociationType

	// This field is used to hold a descriptor parameter for the association, and may therefore only be non-zero if the
	// AssociationType is non-zero. The interpretation of this field is dependent upon the particular AssociationType,
	// and is only used for certain types of associations. If unused, this field should be set to 0x00000000.
	AssociationDesc AssociationDesc

	// This field is optional, and is only used if the object is a member of an association, and only if the association
	// is ordered. If the object is not a member of an ordered association, this value should be set to 0x00000000.
	// These numbers should be created consecutively. However, to be a valid sequence, they do not need to be
	// consecutive, but only monotonically increasing. Therefore, if a data object in the sequence is deleted, the
	// SequenceNumbers of the other objects in the ordered association do not need to be renumbered, and examination of
	// the sequential numbers will indicate a possibly deleted object by the missing sequence number.
	SequenceNumber uint32

	// An optional string representing filename information. This field should not include any filesystem path
	// information, but only the name of the file or directory itself. The interpretation of this string is dependent
	// upon the FilenameFormat field in the StorageInfo dataset that describes the logical storage area in which this
	// object is stored.
	Filename string

	// A static optional field representing the time that the data object was initially captured. This is not
	// necessarily the same as any date held in the ModificationDate field.
	CaptureDate time.Time

	// An optional field representing the time of last modification of the data object. This is not necessarily the same
	// as the CaptureDate field.
	ModificationDate time.Time

	// An optional string representing keywords associated with the image. Each keyword shall be separated by a space. A
	// keyword that consists of more than one word shall use underscore ( _ ) characters to separate individual words
	// within one keyword.
	Keywords string
}
