package viewfinder

import (
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"math"
	"strconv"
	"strings"
)

// NewFujiXT1Viewfinder returns a new Fuji X-T1 viewfinder containing a Widget list mimicking the real viewfinder.
// The image is needed for the widgets to calibrate their origin so they can render in their own designated place.
func NewFujiXT1Viewfinder(img *image.RGBA) *Viewfinder {
	return &Viewfinder{
		Widgets: map[ptp.DevicePropCode]*Widget{
			ptp.DPC_BatteryLevel:             NewFujiBatteryLevelWidget(img),
			ip.DPC_Fuji_CapturesRemaining:    NewFujiCapturesRemainingWidget(img),
			ptp.DPC_ExposureBiasCompensation: NewFujiExposureBiasCompensationWidget(img),
			ptp.DPC_ExposureProgramMode:      NewFujiExposureProgramModeWidget(img),
			ip.DPC_Fuji_ExposureIndex:        NewFujiISOWidget(img),
			ip.DPC_Fuji_FilmSimulation:       NewFujiFilmSimulationWidget(img),
			ptp.DPC_FNumber:                  NewFujiFNumberWidget(img),
			ip.DPC_Fuji_ImageAspectRatio:     NewFujiImageSizeWidget(img),
			ptp.DPC_WhiteBalance:             NewFujiWhiteBalanceWidget(img),
		},
	}
}

func NewFujiBatteryLevelWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.1)
	y := img.Bounds().Max.Y - 8

	w := NewWhiteGlyphWidget(img, int(x), y)
	w.Draw = drawFujiBattery3Bars

	return w
}

func drawFujiBattery3Bars(w *Widget, val int64) {
	w.ResetToOrigin()

	var lvl string
	switch ip.FujiBatteryLevel(val) {
	case ip.BAT_Fuji_3bOne:
		w.SetColour(255, 0, 0) // red
		lvl = "baU"
	case ip.BAT_Fuji_3bTwo:
		lvl = "bCT"
	case ip.BAT_Fuji_3bFull:
		lvl = "BAT"
	}

	w.DrawString(lvl)
}

func NewFujiCapturesRemainingWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.3)
	y := 18

	w := NewWhiteFontWidget(img, int(x), y)
	w.Draw = drawFujiCapturesRemaining

	return w
}

func drawFujiCapturesRemaining(w *Widget, val int64) {
	w.ResetToOrigin()

	w.DrawString(strconv.FormatInt(val, 10))
}

func NewFujiExposureBiasCompensationWidget(img *image.RGBA) *Widget {
	// Make sure the center point of our bias widget is in the center of the image.
	offset := VFGlyphs6x13.Width * len(getBias()) / 2

	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.5) - float64(offset)
	y := img.Bounds().Max.Y - 10

	w := NewWhiteGlyphWidget(img, int(x), y)
	w.Draw = drawFujiExposureBiasCompensation

	return w
}

func getBias() []rune {
	return []rune("6..5..4..0..1..2..3")
}

func drawFujiExposureBiasCompensation(w *Widget, val int64) {
	w.ResetToOrigin()
	w.ResetColour()

	zero := 9  // don't forget: zero indexed!
	stops := 3 // bias dial is per 3 stops
	fStop := 0 // default stop is '0' meaning no fractional stop
	bias := getBias()
	marker := []rune("                   ")

	// Draw the leading +/- icon
	w.Dot.X -= fixed.Int26_6(VFGlyphs6x13.Width * 3 * 64) // offset icon 3 glyphs to the left
	w.DrawString("+-")
	w.ResetToOrigin()

	// Calculate marker position.
	i, f := math.Modf(float64(int16(val)) / float64(1000))
	onZero := i == 0 && f == 0
	if f != 0 {
		fStop = 1
		if math.Abs(f) > 0.4 {
			fStop = 2
		}
		if math.Signbit(f) {
			fStop = -fStop
		}
	}
	pos := zero + fStop + int(i*float64(stops))

	// When we're not on a fractional number, replace the number marker with an 'empty' marker.
	if f == 0 {
		bias[pos] = '"'
	}

	// When the marker is on 0, the widget must be drawn in grey.
	if onZero {
		w.SetColour(100, 100, 100) // grey
	}

	// Now draw the basic exposure bias compensation widget.
	w.DrawString(string(bias))

	// When the marker is on 0, the the marker and '0' position must be drawn in white.
	if onZero {
		w.SetColour(255, 255, 255) // white
		for _, r := range []rune{'"', '!'} {
			w.ResetToOrigin()
			marker[pos] = r
			w.DrawString(string(marker))
		}

		return
	}

	// Draw the marker on the the calculated position in yellow!
	marker[pos] = '!'
	w.SetColour(255, 185, 10) // yellow
	w.ResetToOrigin()
	w.DrawString(string(marker))
}

func NewFujiExposureProgramModeWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Min.X) + (float64(img.Bounds().Max.X) * 0.1)
	y := img.Bounds().Max.Y - 10

	w := NewWhiteGlyphWidget(img, int(x), y)
	w.Draw = drawFujiExposureProgramMode

	return w
}

func drawFujiExposureProgramMode(w *Widget, val int64) {
	w.ResetToOrigin()

	icon := " "
	switch ptp.ExposureProgramMode(val) {
	case ptp.EPM_Manual:
		icon = "Mm"
	case ptp.EPM_Automatic:
		icon = "Pp"
	case ptp.EPM_AperturePriority:
		icon = "Nn"
	case ptp.EPM_ShutterPriority:
		icon = "Ll"
	}

	w.DrawString(icon)
}

func NewFujiISOWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.2)
	y := img.Bounds().Max.Y - 10

	w := NewWhiteGlyphWidget(img, int(x), y)
	w.Draw = drawFujiISO

	return w
}

func drawFujiISO(w *Widget, val int64) {
	w.ResetToOrigin()
	w.ResetFace()

	iso := ptpfmt.FujiExposureIndexAsString(ip.FujiExposureIndex(val))

	w.DrawString("is") // iso icon

	if strings.HasPrefix(iso, "S") {
		w.Dot.X -= fixed.Int26_6(18 * 64) // offset to the left
		w.Dot.Y -= fixed.Int26_6(8 * 64)
		w.DrawString("ISO")              // auto icon
		w.Dot.Y += fixed.Int26_6(8 * 64) // reset Y axis
		iso = string([]rune(iso)[1:])    // drop the leading S
	}

	w.Face = basicfont.Face7x13
	w.Dot.X += fixed.Int26_6(6 * 64)
	w.Dot.Y += fixed.Int26_6(2 * 64)

	w.DrawString(iso) // actual value
}

func NewFujiFilmSimulationWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Min.X) + (float64(img.Bounds().Max.X) * 0.3)
	y := 18

	w := NewWhiteGlyphWidget(img, int(x), y)
	w.Draw = drawFujiFilmSimulation

	return w
}

func drawFujiFilmSimulation(w *Widget, val int64) {
	w.ResetToOrigin()

	var flm string

	switch ip.FujiFilmSimulation(val) {
	case ip.FS_Fuji_Provia:
		flm = "()*"
	case ip.FS_Fuji_Velvia:
		flm = "#$%"
	case ip.FS_Fuji_Astia:
		flm = "(&%"
	case ip.FS_Fuji_Monochrome:
		flm = "'&%"
	case ip.FS_Fuji_Sepia:
		flm = "9DE"
	case ip.FS_Fuji_ProNegHigh:
		flm = "=>;"
	case ip.FS_Fuji_ProNegStandard:
		flm = "=><"
	case ip.FS_Fuji_MonochromeYeFilter:
		flm = "'?@"
	case ip.FS_Fuji_MonochromeRFilter:
		flm = "'?7"
	case ip.FS_Fuji_MonochromeGFilter:
		flm = "'?8"
	case ip.FS_Fuji_ClassicChrome:
		flm = ":,/"
	}

	w.DrawString(flm)
}

func NewFujiFNumberWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Min.X) + (float64(img.Bounds().Max.X) * 0.25)
	y := img.Bounds().Max.Y - 10

	w := NewWhiteFontWidget(img, int(x), y)
	w.Draw = drawFujiFNumber

	return w
}

func drawFujiFNumber(w *Widget, val int64) {
	w.ResetToOrigin()

	w.DrawString(strings.Replace(ptpfmt.FNumberAsString(uint16(val)), "f/", "F", 1))
}

func NewFujiImageSizeWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Max.X) - (float64(img.Bounds().Max.X) * 0.09)
	y := 18

	w := NewWhiteFontWidget(img, int(x), y)
	w.Draw = drawFujiImageSize

	return w
}

func drawFujiImageSize(w *Widget, val int64) {
	w.ResetToOrigin()

	w.DrawString(string([]rune(ptpfmt.FujiImageAspectRatioAsString(ip.FujiImageSize(val)))[0]))
}

func NewFujiWhiteBalanceWidget(img *image.RGBA) *Widget {
	// Calculate starting position.
	x := float64(img.Bounds().Min.X) + (float64(img.Bounds().Max.X) * 0.26)
	y := 18

	w := NewWhiteGlyphWidget(img, int(x), y)
	w.Draw = drawFujiWhiteBalance

	return w
}

func drawFujiWhiteBalance(w *Widget, val int64) {
	w.ResetToOrigin()

	var icon string

	switch ptp.WhiteBalance(val) {
	case ptp.WB_Daylight:
		icon = "XY"
	case ptp.WB_Tungsten:
		icon = "QR"
	case ip.WB_Fuji_Fluorescent1:
		icon = "FGH"
	case ip.WB_Fuji_Fluorescent2:
		icon = "FGJ"
	case ip.WB_Fuji_Fluorescent3:
		icon = "FGK"
	case ip.WB_Fuji_Shade:
		icon = "de"
	case ip.WB_Fuji_Underwater:
		icon = "VW"
	case ip.WB_Fuji_Temperature:
		icon = "]^"
	case ip.WB_Fuji_Custom:
		icon = "Z["
	}

	w.DrawString(icon)
}
