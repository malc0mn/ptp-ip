package ptp

import "testing"

func TestOperationResponseCodeAsString(t *testing.T) {
	check := map[OperationResponseCode]string{
		RC_Undefined:                             "undefined response code",
		RC_OK:                                    "ok",
		RC_GeneralError:                          "general error occured",
		RC_SessionNotOpen:                        "session not open: open a session first",
		RC_InvalidTransactionID:                  "invalid transaction id",
		RC_OperationNotSupported:                 "operation not supported",
		RC_ParameterNotSupported:                 "paramter not supported",
		RC_IncompleteTransfer:                    "incomplete transfer",
		RC_InvalidStorageID:                      "invalid storage id",
		RC_InvalidObjectHandle:                   "invalid object handle",
		RC_DevicePropNotSupported:                "device property not supported",
		RC_InvalidObjectFormatCode:               "invalid object format code",
		RC_StoreFull:                             "store full",
		RC_ObjectWriteProtected:                  "object write protected",
		RC_StoreReadOnly:                         "store read only",
		RC_AccessDenied:                          "access denied",
		RC_NoThumbnailPresent:                    "no thumbnail present",
		RC_SelfTestFailed:                        "self test failed",
		RC_PartialDeletion:                       "partial deletion",
		RC_StoreNotAvailable:                     "store not available",
		RC_SpecificationByFormatUnsupported:      "specification by format unsupported",
		RC_NoValidObjectInfo:                     "no valid object info",
		RC_InvalidCodeFormat:                     "invalid code format",
		RC_UnknownVendorCode:                     "unknown vendor code",
		RC_CaptureAlreadyTerminated:              "capture already terminated",
		RC_DeviceBusy:                            "device busy",
		RC_InvalidParentObject:                   "invalid parent object",
		RC_InvalidDevicePropFormat:               "invalid device property format",
		RC_InvalidDevicePropValue:                "invalid device property value",
		RC_InvalidParameter:                      "invalid parameter",
		RC_SessionAlreadyOpen:                    "session already open",
		RC_TransactionCancelled:                  "transaction cancelled",
		RC_SpecificationofDestinationUnsupported: "specification of destination unsupported",
		OperationResponseCode(5082):              "unknown operation response code: 0x13da",
	}

	for code, want := range check {
		got := OperationResponseCodeAsString(code)
		if got != want {
			t.Errorf("OperationResponseCodeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestGetDeviceInfo(t *testing.T) {
	got := GetDeviceInfo()
	want := OC_GetDeviceInfo
	if got.OperationCode != want {
		t.Errorf("GetDeviceInfo() OperationCode = '%#x', want '%#x'", got.OperationCode, want)
	}
}

func TestOpenSession(t *testing.T) {
	got := OpenSession(1)
	wantCode := OC_OpenSession
	wantParam := uint32(1)
	if got.OperationCode != wantCode {
		t.Errorf("OpenSession() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam {
		t.Errorf("OpenSession() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestCloseSession(t *testing.T) {
	got := CloseSession()
	want := OC_CloseSession
	if got.OperationCode != want {
		t.Errorf("CloseSession() OperationCode = '%#x', want '%#x'", got.OperationCode, want)
	}
}

func TestGetStorageIDs(t *testing.T) {
	got := GetStorageIDs()
	want := OC_GetStorageIDs
	if got.OperationCode != want {
		t.Errorf("GetStorageIDs() OperationCode = '%#x', want '%#x'", got.OperationCode, want)
	}
}

func TestGetStorageInfo(t *testing.T) {
	got := GetStorageInfo(1)
	wantCode := OC_GetStorageInfo
	wantParam := uint32(1)
	if got.OperationCode != wantCode {
		t.Errorf("GetStorageInfo() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam {
		t.Errorf("GetStorageInfo() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestGetNumObjects(t *testing.T) {
	got := GetNumObjects(1, OFC_EXIF_JPEG, 3)
	wantCode := OC_GetNumObjects
	wantParam1 := uint32(1)
	wantParam2 := uint32(OFC_EXIF_JPEG)
	wantParam3 := uint32(3)
	if got.OperationCode != wantCode {
		t.Errorf("GetNumObjects() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("GetNumObjects() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if got.Parameter2 != wantParam2 {
		t.Errorf("GetNumObjects() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
	if got.Parameter3 != wantParam3 {
		t.Errorf("GetNumObjects() Parameter3 = '%#x', want '%#x'", got.Parameter3, wantParam3)
	}
}

func TestGetObjectHandles(t *testing.T) {
	got := GetObjectHandles(1, OFC_EXIF_JPEG, 3)
	wantCode := OC_GetObjectHandles
	wantParam1 := uint32(1)
	wantParam2 := uint32(OFC_EXIF_JPEG)
	wantParam3 := uint32(3)
	if got.OperationCode != wantCode {
		t.Errorf("GetNumObjects() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("GetNumObjects() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if got.Parameter2 != wantParam2 {
		t.Errorf("GetNumObjects() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
	if got.Parameter3 != wantParam3 {
		t.Errorf("GetNumObjects() Parameter3 = '%#x', want '%#x'", got.Parameter3, wantParam3)
	}
}

func TestGetObjectInfo(t *testing.T) {
	got := GetObjectInfo(1)
	wantCode := OC_GetObjectInfo
	wantParam := uint32(1)
	if got.OperationCode != wantCode {
		t.Errorf("GetObjectInfo() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam {
		t.Errorf("GetObjectInfo() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestGetObject(t *testing.T) {
	got := GetObject(1)
	wantCode := OC_GetObject
	wantParam := uint32(1)
	if got.OperationCode != wantCode {
		t.Errorf("GetObject() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam {
		t.Errorf("GetObject() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestGetThumb(t *testing.T) {
	got := GetThumb(1)
	wantCode := OC_GetThumb
	wantParam := uint32(1)
	if got.OperationCode != wantCode {
		t.Errorf("GetThumb() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam {
		t.Errorf("GetThumb() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestDeleteObject(t *testing.T) {
	got := DeleteObject(1, OFC_AVI)
	wantCode := OC_DeleteObject
	wantParam1 := uint32(1)
	wantParam2 := OFC_AVI
	if got.OperationCode != wantCode {
		t.Errorf("DeleteObject() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("DeleteObject() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if ObjectFormatCode(got.Parameter2) != wantParam2 {
		t.Errorf("DeleteObject() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
}

func TestSendObjectInfo(t *testing.T) {
	got := SendObjectInfo(1, 2)
	wantCode := OC_SendObjectInfo
	wantParam1 := uint32(1)
	wantParam2 := uint32(2)
	if got.OperationCode != wantCode {
		t.Errorf("SendObjectInfo() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("SendObjectInfo() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if got.Parameter2 != wantParam2 {
		t.Errorf("SendObjectInfo() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
}

func TestSendObject(t *testing.T) {
	got := SendObject()
	want := OC_SendObject
	if got.OperationCode != want {
		t.Errorf("SendObject() OperationCode = '%#x', want '%#x'", got.OperationCode, want)
	}
}

func TestInitiateCapture(t *testing.T) {
	got := InitiateCapture(1, OFC_EXIF_JPEG)
	wantCode := OC_InitiateCapture
	wantParam1 := uint32(1)
	wantParam2 := OFC_EXIF_JPEG
	if got.OperationCode != wantCode {
		t.Errorf("InitiateCapture() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("InitiateCapture() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if ObjectFormatCode(got.Parameter2) != wantParam2 {
		t.Errorf("InitiateCapture() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
}

func TestFormatStore(t *testing.T) {
	got := FormatStore(1, FT_GenericFlat)
	wantCode := OC_FormatStore
	wantParam1 := uint32(1)
	wantParam2 := FT_GenericFlat
	if got.OperationCode != wantCode {
		t.Errorf("FormatStore() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("FormatStore() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if FilesystemType(got.Parameter2) != wantParam2 {
		t.Errorf("FormatStore() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
}

func TestResetDevice(t *testing.T) {
	got := ResetDevice()
	want := OC_ResetDevice
	if got.OperationCode != want {
		t.Errorf("ResetDevice() OperationCode = '%#x', want '%#x'", got.OperationCode, want)
	}
}

func TestSelfTest(t *testing.T) {
	got := SelfTest(STT_Default)
	wantCode := OC_SelfTest
	wantParam := STT_Default
	if got.OperationCode != wantCode {
		t.Errorf("SelfTest() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if SelfTestType(got.Parameter1) != wantParam {
		t.Errorf("SelfTest() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestSetObjectProtection(t *testing.T) {
	got := SetObjectProtection(1, PS_ReadOnly)
	wantCode := OC_SetObjectProtection
	wantParam1 := uint32(1)
	wantParam2 := PS_ReadOnly
	if got.OperationCode != wantCode {
		t.Errorf("SetObjectProtection() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("SetObjectProtection() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if ProtectionStatus(got.Parameter2) != wantParam2 {
		t.Errorf("SetObjectProtection() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
}

func TestPowerDown(t *testing.T) {
	got := PowerDown()
	want := OC_PowerDown
	if got.OperationCode != want {
		t.Errorf("PowerDown() OperationCode = '%#x', want '%#x'", got.OperationCode, want)
	}
}

func TestGetDevicePropDesc(t *testing.T) {
	got := GetDevicePropDesc(DPC_FocusMeteringMode)
	wantCode := OC_GetDevicePropDesc
	wantParam := DPC_FocusMeteringMode
	if got.OperationCode != wantCode {
		t.Errorf("GetDevicePropDesc() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if DevicePropCode(got.Parameter1) != wantParam {
		t.Errorf("GetDevicePropDesc() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestGetDevicePropValue(t *testing.T) {
	got := GetDevicePropValue(DPC_BatteryLevel)
	wantCode := OC_GetDevicePropValue
	wantParam := DPC_BatteryLevel
	if got.OperationCode != wantCode {
		t.Errorf("GetDevicePropValue() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if DevicePropCode(got.Parameter1) != wantParam {
		t.Errorf("GetDevicePropValue() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestSetDevicePropValue(t *testing.T) {
	got := SetDevicePropValue(DPC_WhiteBalance, WB_Automatic)
	wantCode := OC_SetDevicePropValue
	wantParam := DPC_WhiteBalance
	if got.OperationCode != wantCode {
		t.Errorf("SetDevicePropValue() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if DevicePropCode(got.Parameter1) != wantParam {
		t.Errorf("SetDevicePropValue() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestResetDevicePropValue(t *testing.T) {
	got := ResetDevicePropValue(DPC_ExposureBiasCompensation)
	wantCode := OC_ResetDevicePropValue
	wantParam := DPC_ExposureBiasCompensation
	if got.OperationCode != wantCode {
		t.Errorf("ResetDevicePropValue() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if DevicePropCode(got.Parameter1) != wantParam {
		t.Errorf("ResetDevicePropValue() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestTerminateOpenCapture(t *testing.T) {
	got := TerminateOpenCapture(1)
	wantCode := OC_TerminateOpenCapture
	wantParam := uint32(1)
	if got.OperationCode != wantCode {
		t.Errorf("TerminateOpenCapture() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam {
		t.Errorf("TerminateOpenCapture() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam)
	}
}

func TestMoveObject(t *testing.T) {
	got := MoveObject(1, 5, 3)
	wantCode := OC_MoveObject
	wantParam1 := uint32(1)
	wantParam2 := uint32(5)
	wantParam3 := uint32(3)
	if got.OperationCode != wantCode {
		t.Errorf("MoveObject() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("MoveObject() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if got.Parameter2 != wantParam2 {
		t.Errorf("MoveObject() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
	if got.Parameter3 != wantParam3 {
		t.Errorf("MoveObject() Parameter3 = '%#x', want '%#x'", got.Parameter3, wantParam3)
	}
}

func TestCopyObject(t *testing.T) {
	got := CopyObject(1, 5, 3)
	wantCode := OC_CopyObject
	wantParam1 := uint32(1)
	wantParam2 := uint32(5)
	wantParam3 := uint32(3)
	if got.OperationCode != wantCode {
		t.Errorf("CopyObject() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("CopyObject() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if got.Parameter2 != wantParam2 {
		t.Errorf("CopyObject() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
	if got.Parameter3 != wantParam3 {
		t.Errorf("CopyObject() Parameter3 = '%#x', want '%#x'", got.Parameter3, wantParam3)
	}
}

func TestGetPartialObject(t *testing.T) {
	got := GetPartialObject(1, 5, 3500)
	wantCode := OC_GetPartialObject
	wantParam1 := uint32(1)
	wantParam2 := uint32(5)
	wantParam3 := uint32(3500)
	if got.OperationCode != wantCode {
		t.Errorf("GetPartialObject() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("GetPartialObject() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if got.Parameter2 != wantParam2 {
		t.Errorf("GetPartialObject() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
	if got.Parameter3 != wantParam3 {
		t.Errorf("GetPartialObject() Parameter3 = '%#x', want '%#x'", got.Parameter3, wantParam3)
	}
}

func TestInitiateOpenCapture(t *testing.T) {
	got := InitiateOpenCapture(1, OFC_EXIF_JPEG)
	wantCode := OC_InitiateOpenCapture
	wantParam1 := uint32(1)
	wantParam2 := OFC_EXIF_JPEG
	if got.OperationCode != wantCode {
		t.Errorf("InitiateOpenCapture() OperationCode = '%#x', want '%#x'", got.OperationCode, wantCode)
	}
	if got.Parameter1 != wantParam1 {
		t.Errorf("InitiateOpenCapture() Parameter1 = '%#x', want '%#x'", got.Parameter1, wantParam1)
	}
	if ObjectFormatCode(got.Parameter2) != wantParam2 {
		t.Errorf("InitiateOpenCapture() Parameter2 = '%#x', want '%#x'", got.Parameter2, wantParam2)
	}
}
