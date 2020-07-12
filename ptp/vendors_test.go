package ptp

import "testing"

func TestVendorStringToType(t *testing.T) {
	check := map[string]VendorExtension{
		"kodak":     VE_EastmanKodakCompany,
		"epson":     VE_SeikoEpson,
		"agilent":   VE_AgilentTechnologiesInc,
		"polaroid":  VE_PolaroidCorporation,
		"agfa":      VE_AgfaGevaert,
		"ms":        VE_MicrosoftCorporation,
		"equinox":   VE_EquinoxResearchLtd,
		"vq":        VE_ViewQuestTechnologies,
		"st":        VE_STMicroelectronics,
		"nikon":     VE_NikonCorporation,
		"canon":     VE_CanonInc,
		"fn":        VE_FotoNationInc,
		"pentax":    VE_PENTAXCorporation,
		"fuji":      VE_FujiPhotoFilmCoLtd,
		"ndd":       VE_NddMedicalTechnologies,
		"samsung":   VE_SamsungElectronicsCoLtd,
		"parrot":    VE_ParrotDronesSAS,
		"panasonic": VE_PanasonicCorporation,
		"generic":   0,
	}

	for code, want := range check {
		got := VendorStringToType(code)
		if got != want {
			t.Errorf("VendorStringToType() return = %#x, want %#x", got, want)
		}
	}
}
