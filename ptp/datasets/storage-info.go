package ptp

import ptp "github.com/malc0mn/ptp-ip/ptp/consts"

// This dataset is used to hold the state information for a storage device.
type StorageInfo struct {
	// The code that identifies the type of storage, particularly whether the store is inherently random-access or
	// read-only memory, and whether it is fixed or removable media.
	StorageType ptp.StorageType

	// This optional code indicates the type of filesystem present on the device. This field may be used to determine
	// the filenaming convention used by the storage device, as well as to determine whether support for a hierarchical
	// system is present. If the storage device is DCF-conformant, it shall indicate so here.
	// All values having bit 31 set to zero are reserved for future use. If a proprietary implementation wishes to
	// extend the interpretation of this field, bit 31 should be set to 1.
	FilesystemType ptp.FilesysytemType

	// This field indicates whether the store is read-write or read-only. If the store is read-only, deletion may or may
	// not be allowed. The allowed values are described in the following table. Read-Write is only valid if the
	// StorageType is nonROM, as described in the StorageType field above.
	// All values having bit 15 set to zero are reserved for future use. If a proprietary implementation wishes to
	// extend the interpretation of this field, bit 15 should be set to 1.
	AccessCapability ptp.AccessCapability

	// This is an optional field that indicates the total storage capacity of the store in bytes. If this field is
	// unused, it should report 0xFFFFFFFF.
	MaxCapacity uint64

	// The amount of free space that is available in the store in bytes. If this value is not useful for the device, it
	// may set this field to 0xFFFFFFFF and rely upon the FreeSpaceInImages field instead.
	FreeSpaceInBytes uint64

	// The number of images that may still be captured into this store according to the current image capture settings
	// of the device. If the device does not implement this capability, this field should be set to 0xFFFFFFFF. This
	// field may be used for devices that do not report FreeSpaceInBytes, or the two fields may be used in combination.
	FreeSpaceInImages uint32

	// An optional field that may be used for a human-readable text description of the storage device. This should be
	// used for storage-type specific information as opposed to volume-specific information. Examples would be "Type I
	// Compact Flash" or "3.5-inch 1.44 MB Floppy". If unused, this field should be set to the empty string.
	StorageDescription string

	// An optional field that may be used to hold the volume label of the storage device, if such a label exists and is
	// known. If unused, this field should be set to the empty string.
	VolumeLabel string
}
