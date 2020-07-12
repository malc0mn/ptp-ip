package ptp

type VendorExtension uint32

const (
	VE_EastmanKodakCompany     VendorExtension = 0x00000001
	VE_SeikoEpson              VendorExtension = 0x00000002
	VE_AgilentTechnologiesInc  VendorExtension = 0x00000003
	VE_PolaroidCorporation     VendorExtension = 0x00000004
	VE_AgfaGevaert             VendorExtension = 0x00000005
	VE_MicrosoftCorporation    VendorExtension = 0x00000006
	VE_EquinoxResearchLtd      VendorExtension = 0x00000007
	VE_ViewQuestTechnologies   VendorExtension = 0x00000008
	VE_STMicroelectronics      VendorExtension = 0x00000009
	VE_NikonCorporation        VendorExtension = 0x0000000A
	VE_CanonInc                VendorExtension = 0x0000000B
	VE_FotoNationInc           VendorExtension = 0x0000000C
	VE_PENTAXCorporation       VendorExtension = 0x0000000D
	VE_FujiPhotoFilmCoLtd      VendorExtension = 0x0000000E
	VE_NddMedicalTechnologies  VendorExtension = 0x00000012
	VE_SamsungElectronicsCoLtd VendorExtension = 0x0000001A
	VE_ParrotDronesSAS         VendorExtension = 0x0000001B
	VE_PanasonicCorporation    VendorExtension = 0x0000001C
)

func VendorStringToType(vendor string) VendorExtension {
	switch vendor {
	case "kodak":
		return VE_EastmanKodakCompany
	case "epson":
		return VE_SeikoEpson
	case "agilent":
		return VE_AgilentTechnologiesInc
	case "polaroid":
		return VE_PolaroidCorporation
	case "agfa":
		return VE_AgfaGevaert
	case "ms":
		return VE_MicrosoftCorporation
	case "equinox":
		return VE_EquinoxResearchLtd
	case "vq":
		return VE_ViewQuestTechnologies
	case "st":
		return VE_STMicroelectronics
	case "nikon":
		return VE_NikonCorporation
	case "canon":
		return VE_CanonInc
	case "fn":
		return VE_FotoNationInc
	case "pentax":
		return VE_PENTAXCorporation
	case "fuji":
		return VE_FujiPhotoFilmCoLtd
	case "ndd":
		return VE_NddMedicalTechnologies
	case "samsung":
		return VE_SamsungElectronicsCoLtd
	case "parrot":
		return VE_ParrotDronesSAS
	case "panasonic":
		return VE_PanasonicCorporation
	default:
		return 0
	}
}
