package ptp

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
	OFC_AVI ObjectFormatCode = 0x300a
	// Video clip
	OFC_MPEG ObjectFormatCode = 0x300b
	// Microsoft Advanced Streaming Format (video)
	OFC_ASF ObjectFormatCode = 0x300c
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
	OFC_PICT ObjectFormatCode = 0x380a
	// Portable Network Graphics
	OFC_PNG ObjectFormatCode = 0x380b
	// Tag Image File Format
	OFC_TIFF ObjectFormatCode = 0x380d
	// Tag Image File Format for Information Technology (graphic arts)
	OFC_TIFF_IT ObjectFormatCode = 0x380e
	// JPEG2000 Baseline File Format
	OFC_JP2 ObjectFormatCode = 0x380f
	// JPEG2000 Extended File Format
	OFC_JPX ObjectFormatCode = 0x3810
)
