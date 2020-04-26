package ptp

type StorageType uint16
type FilesysytemType uint16
type AccessCapability uint16
type ProtectionStatus uint16
type StorageID uint32

const (
	ST_Undefined    StorageType = 0x0000
	ST_FixedROM     StorageType = 0x0001
	ST_RemovableROM StorageType = 0x0002
	ST_FixedRAM     StorageType = 0x0003
	ST_RemovableRAM StorageType = 0x0004

	FT_Undefined           FilesysytemType = 0x0000
	FT_GenericFlat         FilesysytemType = 0x0001
	FT_GenericHierarchical FilesysytemType = 0x0002
	FT_DCF                 FilesysytemType = 0x0003

	AC_ReadWrite           AccessCapability = 0x0000
	AC_ReadOnly_NoDeletion AccessCapability = 0x0001
	AC_ReadOnly_Deletion   AccessCapability = 0x0002

	PS_NoProtection ProtectionStatus = 0x0000
	PS_ReadOnly     ProtectionStatus = 0x0001
)
