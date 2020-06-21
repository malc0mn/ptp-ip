package ptp

type EffectMode uint16
type ExposureMeteringMode uint16
type ExposureProgramMode uint16
type FlashMode uint16
type FocusMeteringMode uint16
type FocusMode uint16
type FunctionalMode uint16
type SelfTestType uint16
type StillCaptureMode uint16
type WhiteBalance uint16

const (
	EMM_Undefined             ExposureMeteringMode = 0x0000
	EMM_Avarage               ExposureMeteringMode = 0x0001
	EMM_CenterWeightedAvarage ExposureMeteringMode = 0x0002
	EMM_MultiSpot             ExposureMeteringMode = 0x0003
	EMM_CenterSpot            ExposureMeteringMode = 0x0004

	EPM_Undefined        ExposureProgramMode = 0x0000
	EPM_Manual           ExposureProgramMode = 0x0001
	EPM_Automatic        ExposureProgramMode = 0x0002
	EPM_AperturePriority ExposureProgramMode = 0x0003
	EPM_ShutterPriority  ExposureProgramMode = 0x0004
	EPM_ProgramCreative  ExposureProgramMode = 0x0005
	EPM_ProgramAction    ExposureProgramMode = 0x0006
	EPM_Portrait         ExposureProgramMode = 0x0007

	FCM_Undefined      FocusMode = 0x0000
	FCM_Manual         FocusMode = 0x0001
	FCM_Automatic      FocusMode = 0x0002
	FCM_AutomaticMacro FocusMode = 0x0003

	FLM_Undefined    FlashMode = 0x0000
	FLM_AutoFlash    FlashMode = 0x0001
	FLM_FlashOff     FlashMode = 0x0002
	FLM_FillFlash    FlashMode = 0x0003
	FLM_RedEyeAuto   FlashMode = 0x0004
	FLM_RedEyeFill   FlashMode = 0x0005
	FLM_ExternalSync FlashMode = 0x0006

	FMM_Undefined  FocusMeteringMode = 0x0000
	FMM_CenterSpot FocusMeteringMode = 0x0001
	FMM_MultiSpot  FocusMeteringMode = 0x0002

	FUM_StandardMode FunctionalMode = 0x0000
	FUM_SleepState   FunctionalMode = 0x0001

	FXM_Undefined  EffectMode = 0x0000
	FXM_Standard   EffectMode = 0x0001
	FXM_BlackWhite EffectMode = 0x0002
	FXM_Sepia      EffectMode = 0x0003

	SCM_Undefined StillCaptureMode = 0x0000
	SCM_Normal    StillCaptureMode = 0x0001
	SCM_Burst     StillCaptureMode = 0x0002
	SCM_Timelapse StillCaptureMode = 0x0003

	// STT_Default is the default device-specific self-test.
	STT_Default SelfTestType = 0x0000

	WB_Undefined WhiteBalance = 0x0000
	// WB_Manual indicates the white balance is set directly using the RGB Gain property and is static until changed.
	WB_Manual WhiteBalance = 0x0001
	// WB_Automatic indicates the device attempts to set the white balance using some kind of automatic mechanism.
	WB_Automatic WhiteBalance = 0x0002
	// WB_OnePushAutomatic indicates the user must press the capture button while pointing the device at a white field,
	// at which time the device determines the white balance setting.
	WB_OnePushAutomatic WhiteBalance = 0x0003
	// WB_Daylight indicates the device attempts to set the white balance to a value that is appropriate for use in
	// daylight conditions.
	WB_Daylight    WhiteBalance = 0x0004
	// WB_Fluorescent indicates the device attempts to set the white balance to a value that is appropriate for use in
	// with fluorescent lighting conditions.
	WB_Fluorescent WhiteBalance = 0x0005
	// WB_Tungsten indicates the device attempts to set the white balance to a value that is appropriate for use in
	// conditions with a tungsten light source.
	WB_Tungsten WhiteBalance = 0x0006
	// WB_Flash indicates the device attempts to set the white balance to a value that is appropriate for flash
	// conditions.
	WB_Flash WhiteBalance = 0x0007
)
