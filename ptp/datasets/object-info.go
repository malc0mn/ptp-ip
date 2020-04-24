package ptp

import (
	"github.com/malc0mn/ptp-ip/ptp/consts"
	"time"
)

// This dataset is used to define the information about data objects in persistent store, as well as optional
// information if the data is known to be an image or an association object. It is required that these data items be
// accounted for in response to a GetObjectInfo operation. If the data is not known to be an image, or the image
// information is unavailable, the image-specific fields shall be set to zero. Objects of type Association are fully
// qualified by the ObjectInfo dataset.
type ObjectInfo struct {
	// The StorageID of the device's store in which the image resides.
	StorageID ptp.StorageID

	//  Indicates ObjectFormatCode of the object.
	ObjectFormat ptp.ObjectFormatCode

	// An optional field representing the write-protection status of the data object. Objects that are protected may not
	// be deleted as the result of any operations specified in this standard without first separately removing their
	// protection status in a separate transaction.
	// This protection field is distinctly different in scope than the AccessCapability field present in the StorageInfo
	// dataset. If an attempt to delete an object is made, success will only occur if the ProtectionStatus of the object
	// is 0x0000 and the AccessCapability of the store allows deletion. If a device does not support object protection,
	// this field should always be set to 0x0000, and the SetProtection operation should not be supported.
	ProtectionStatus ptp.ProtectionStatus

	// The size of the buffer needed to hold the entire binary object in bytes. This field may be used for memory
	// allocation purposes in object receivers by transport implementations.
	ObjectCompressedSize uint32

	// Indicates ObjectFormatCode of the thumbnail. In order for an object to be referred to as an image, it must be able to
	// produce a thumbnail as the response to a request. Therefore, this value should only be 0x00000000 for the case of
	// non-image objects.
	ThumbFormat ptp.ObjectFormatCode

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
	ParentObject ptp.ObjectHandle

	// A field that is only used for objects of type Association. This code indicates the type of association. If the
	// object is not an association, this field should be set to 0x0000.
	AssociationType ptp.AssociationType

	// This field is used to hold a descriptor parameter for the association, and may therefore only be non-zero if the
	// AssociationType is non-zero. The interpretation of this field is dependent upon the particular AssociationType,
	// and is only used for certain types of associations. If unused, this field should be set to 0x00000000.
	AssociationDesc ptp.AssociationDesc

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
