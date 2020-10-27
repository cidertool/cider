package config

import (
	"testing"

	"github.com/cidertool/asc-go/asc"
	"github.com/stretchr/testify/assert"
)

func TestPlatformAPIValue(t *testing.T) {
	var plat Platform
	plat = PlatformiOS
	assert.Equal(t, *plat.APIValue(), asc.PlatformIOS)
	plat = PlatformMacOS
	assert.Equal(t, *plat.APIValue(), asc.PlatformMACOS)
	plat = PlatformTvOS
	assert.Equal(t, *plat.APIValue(), asc.PlatformTVOS)

	bad := Platform("watchOS")
	assert.Empty(t, bad.APIValue())

	var empty *Platform

	assert.Empty(t, empty.APIValue())
}

func TestReleaseTypeAPIValue(t *testing.T) {
	var release releaseType
	release = ReleaseTypeManual
	assert.Equal(t, *release.APIValue(), "MANUAL")
	release = ReleaseTypeAfterApproval
	assert.Equal(t, *release.APIValue(), "AFTER_APPROVAL")
	release = ReleaseTypeScheduled
	assert.Equal(t, *release.APIValue(), "SCHEDULED")

	bad := releaseType("never")
	assert.Empty(t, bad.APIValue())

	var empty *releaseType

	assert.Empty(t, empty.APIValue())
}

func TestContentIntensityAPIValue(t *testing.T) {
	var intensity contentIntensity
	intensity = ContentIntensityNone
	assert.Equal(t, *intensity.APIValue(), "NONE")
	intensity = ContentIntensityInfrequentOrMild
	assert.Equal(t, *intensity.APIValue(), "INFREQUENT_OR_MILD")
	intensity = ContentIntensityFrequentOrIntense
	assert.Equal(t, *intensity.APIValue(), "FREQUENT_OR_INTENSE")

	bad := contentIntensity("nothing but violence")
	assert.Empty(t, bad.APIValue())

	var empty *contentIntensity

	assert.Empty(t, empty.APIValue())
}

func TestKidsAgeBandAPIValue(t *testing.T) {
	var band kidsAgeBand
	band = KidsAgeBandFiveAndUnder
	assert.Equal(t, *band.APIValue(), asc.KidsAgeBandFiveAndUnder)
	band = KidsAgeBandSixToEight
	assert.Equal(t, *band.APIValue(), asc.KidsAgeBandSixToEight)
	band = KidsAgeBandNineToEleven
	assert.Equal(t, *band.APIValue(), asc.KidsAgeBandNineToEleven)

	bad := kidsAgeBand("18+")
	assert.Empty(t, bad.APIValue())

	var empty *kidsAgeBand

	assert.Empty(t, empty.APIValue())
}

func TestPreviewTypeAPIValue(t *testing.T) {
	var preview previewType
	preview = PreviewTypeAppleTV
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeAppleTV)
	preview = PreviewTypeDesktop
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeDesktop)
	preview = PreviewTypeiPad105
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPad105)
	preview = PreviewTypeiPad97
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPad97)
	preview = PreviewTypeiPadPro129
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPadPro129)
	preview = PreviewTypeiPadPro3Gen11
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPadPro3Gen11)
	preview = PreviewTypeiPadPro3Gen129
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPadPro3Gen129)
	preview = PreviewTypeiPhone35
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPhone35)
	preview = PreviewTypeiPhone40
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPhone40)
	preview = PreviewTypeiPhone47
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPhone47)
	preview = PreviewTypeiPhone55
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPhone55)
	preview = PreviewTypeiPhone58
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPhone58)
	preview = PreviewTypeiPhone65
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeiPhone65)
	preview = PreviewTypeWatchSeries3
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeWatchSeries3)
	preview = PreviewTypeWatchSeries4
	assert.Equal(t, *preview.APIValue(), asc.PreviewTypeWatchSeries4)

	bad := previewType("Google Pixel")
	assert.Empty(t, bad.APIValue())

	var empty *previewType

	assert.Empty(t, empty.APIValue())
}

func TestScreenshotTypeAPIValue(t *testing.T) {
	var screenshot screenshotType
	screenshot = ScreenshotTypeAppleTV
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppAppleTV)
	screenshot = ScreenshotTypeDesktop
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppDesktop)
	screenshot = ScreenshotTypeiPad105
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPad105)
	screenshot = ScreenshotTypeiPad97
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPad97)
	screenshot = ScreenshotTypeiPadPro129
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPadPro129)
	screenshot = ScreenshotTypeiPadPro3Gen11
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPadPro3Gen11)
	screenshot = ScreenshotTypeiPadPro3Gen129
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPadPro3Gen129)
	screenshot = ScreenshotTypeiPhone35
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPhone35)
	screenshot = ScreenshotTypeiPhone40
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPhone40)
	screenshot = ScreenshotTypeiPhone47
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPhone47)
	screenshot = ScreenshotTypeiPhone55
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPhone55)
	screenshot = ScreenshotTypeiPhone58
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPhone58)
	screenshot = ScreenshotTypeiPhone65
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppiPhone65)
	screenshot = ScreenshotTypeWatchSeries3
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppWatchSeries3)
	screenshot = ScreenshotTypeWatchSeries4
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeAppWatchSeries4)
	screenshot = ScreenshotTypeiMessageiPad105
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPad105)
	screenshot = ScreenshotTypeiMessageiPad97
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPad97)
	screenshot = ScreenshotTypeiMessageiPadPro129
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPadPro129)
	screenshot = ScreenshotTypeiMessageiPadPro3Gen11
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen11)
	screenshot = ScreenshotTypeiMessageiPadPro3Gen129
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen129)
	screenshot = ScreenshotTypeiMessageiPhone40
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPhone40)
	screenshot = ScreenshotTypeiMessageiPhone47
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPhone47)
	screenshot = ScreenshotTypeiMessageiPhone55
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPhone55)
	screenshot = ScreenshotTypeiMessageiPhone58
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPhone58)
	screenshot = ScreenshotTypeiMessageiPhone65
	assert.Equal(t, *screenshot.APIValue(), asc.ScreenshotDisplayTypeiMessageAppIPhone65)

	bad := screenshotType("Google Pixel")
	assert.Empty(t, bad.APIValue())

	var empty *screenshotType

	assert.Empty(t, empty.APIValue())
}

func TestPreviewSetsGetPreviews(t *testing.T) {
	sets := make(PreviewSets)

	sets[PreviewTypeAppleTV] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeAppleTV))
	sets[PreviewTypeDesktop] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeDesktop))
	sets[PreviewTypeiPad105] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPad105))
	sets[PreviewTypeiPad97] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPad97))
	sets[PreviewTypeiPadPro129] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPadPro129))
	sets[PreviewTypeiPadPro3Gen11] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPadPro3Gen11))
	sets[PreviewTypeiPadPro3Gen129] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPadPro3Gen129))
	sets[PreviewTypeiPhone35] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPhone35))
	sets[PreviewTypeiPhone40] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPhone40))
	sets[PreviewTypeiPhone47] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPhone47))
	sets[PreviewTypeiPhone55] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPhone55))
	sets[PreviewTypeiPhone58] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPhone58))
	sets[PreviewTypeiPhone65] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeiPhone65))
	sets[PreviewTypeWatchSeries3] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeWatchSeries3))
	sets[PreviewTypeWatchSeries4] = []Preview{}
	assert.Empty(t, sets.GetPreviews(asc.PreviewTypeWatchSeries4))
	assert.Empty(t, sets.GetPreviews(""))
}

func TestGetScreenshotSetsGetScreenshots(t *testing.T) {
	sets := make(ScreenshotSets)

	sets[ScreenshotTypeAppleTV] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppAppleTV))
	sets[ScreenshotTypeDesktop] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppDesktop))
	sets[ScreenshotTypeiPad105] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPad105))
	sets[ScreenshotTypeiPad97] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPad97))
	sets[ScreenshotTypeiPadPro129] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPadPro129))
	sets[ScreenshotTypeiPadPro3Gen11] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPadPro3Gen11))
	sets[ScreenshotTypeiPadPro3Gen129] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPadPro3Gen129))
	sets[ScreenshotTypeiPhone35] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPhone35))
	sets[ScreenshotTypeiPhone40] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPhone40))
	sets[ScreenshotTypeiPhone47] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPhone47))
	sets[ScreenshotTypeiPhone55] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPhone55))
	sets[ScreenshotTypeiPhone58] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPhone58))
	sets[ScreenshotTypeiPhone65] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppiPhone65))
	sets[ScreenshotTypeWatchSeries3] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppWatchSeries3))
	sets[ScreenshotTypeWatchSeries4] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeAppWatchSeries4))
	sets[ScreenshotTypeiMessageiPad105] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPad105))
	sets[ScreenshotTypeiMessageiPad97] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPad97))
	sets[ScreenshotTypeiMessageiPadPro129] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPadPro129))
	sets[ScreenshotTypeiMessageiPadPro3Gen11] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen11))
	sets[ScreenshotTypeiMessageiPadPro3Gen129] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen129))
	sets[ScreenshotTypeiMessageiPhone40] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPhone40))
	sets[ScreenshotTypeiMessageiPhone47] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPhone47))
	sets[ScreenshotTypeiMessageiPhone55] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPhone55))
	sets[ScreenshotTypeiMessageiPhone58] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPhone58))
	sets[ScreenshotTypeiMessageiPhone65] = []File{}
	assert.Empty(t, sets.GetScreenshots(asc.ScreenshotDisplayTypeiMessageAppIPhone65))
	assert.Empty(t, sets.GetScreenshots(""))
}
