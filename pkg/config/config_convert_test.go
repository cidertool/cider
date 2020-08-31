package config

import (
	"testing"

	"github.com/aaronsky/asc-go/asc"
	"github.com/stretchr/testify/assert"
)

func TestPlatformAPIValue(t *testing.T) {
	ios, err := PlatformiOS.APIValue()
	assert.NoError(t, err)
	assert.Equal(t, ios, asc.PlatformIOS)
	macos, err := PlatformMacOS.APIValue()
	assert.NoError(t, err)
	assert.Equal(t, macos, asc.PlatformMACOS)
	tvos, err := PlatformTvOS.APIValue()
	assert.NoError(t, err)
	assert.Equal(t, tvos, asc.PlatformTVOS)
	_, err = platform("watchOS").APIValue()
	assert.Error(t, err)
}

func TestReleaseTypeAPIValue(t *testing.T) {
	manual, err := ReleaseTypeManual.APIValue()
	assert.NoError(t, err)
	assert.Equal(t, manual, "MANUAL")
	afterApproval, err := ReleaseTypeAfterApproval.APIValue()
	assert.NoError(t, err)
	assert.Equal(t, afterApproval, "AFTER_APPROVAL")
	scheduled, err := ReleaseTypeScheduled.APIValue()
	assert.NoError(t, err)
	assert.Equal(t, scheduled, "SCHEDULED")
	_, err = releaseType("never").APIValue()
	assert.Error(t, err)
}

func TestPreviewTypeAPIValue(t *testing.T) {
	appleTV := PreviewTypeAppleTV.APIValue()
	assert.Equal(t, appleTV, asc.PreviewTypeAppleTV)
	desktop := PreviewTypeDesktop.APIValue()
	assert.Equal(t, desktop, asc.PreviewTypeDesktop)
	iPad105 := PreviewTypeiPad105.APIValue()
	assert.Equal(t, iPad105, asc.PreviewTypeiPad105)
	iPad97 := PreviewTypeiPad97.APIValue()
	assert.Equal(t, iPad97, asc.PreviewTypeiPad97)
	iPadPro129 := PreviewTypeiPadPro129.APIValue()
	assert.Equal(t, iPadPro129, asc.PreviewTypeiPadPro129)
	iPadPro3Gen11 := PreviewTypeiPadPro3Gen11.APIValue()
	assert.Equal(t, iPadPro3Gen11, asc.PreviewTypeiPadPro3Gen11)
	iPadPro3Gen129 := PreviewTypeiPadPro3Gen129.APIValue()
	assert.Equal(t, iPadPro3Gen129, asc.PreviewTypeiPadPro3Gen129)
	iPhone35 := PreviewTypeiPhone35.APIValue()
	assert.Equal(t, iPhone35, asc.PreviewTypeiPhone35)
	iPhone40 := PreviewTypeiPhone40.APIValue()
	assert.Equal(t, iPhone40, asc.PreviewTypeiPhone40)
	iPhone47 := PreviewTypeiPhone47.APIValue()
	assert.Equal(t, iPhone47, asc.PreviewTypeiPhone47)
	iPhone55 := PreviewTypeiPhone55.APIValue()
	assert.Equal(t, iPhone55, asc.PreviewTypeiPhone55)
	iPhone58 := PreviewTypeiPhone58.APIValue()
	assert.Equal(t, iPhone58, asc.PreviewTypeiPhone58)
	iPhone65 := PreviewTypeiPhone65.APIValue()
	assert.Equal(t, iPhone65, asc.PreviewTypeiPhone65)
	watchSeries3 := PreviewTypeWatchSeries3.APIValue()
	assert.Equal(t, watchSeries3, asc.PreviewTypeWatchSeries3)
	watchSeries4 := PreviewTypeWatchSeries4.APIValue()
	assert.Equal(t, watchSeries4, asc.PreviewTypeWatchSeries4)
	bad := previewType("Google Pixel").APIValue()
	assert.Empty(t, bad)
}

func TestScreenshotTypeAPIValue(t *testing.T) {
	appleTV := ScreenshotTypeAppleTV.APIValue()
	assert.Equal(t, appleTV, asc.ScreenshotDisplayTypeAppAppleTV)
	desktop := ScreenshotTypeDesktop.APIValue()
	assert.Equal(t, desktop, asc.ScreenshotDisplayTypeAppDesktop)
	iPad105 := ScreenshotTypeiPad105.APIValue()
	assert.Equal(t, iPad105, asc.ScreenshotDisplayTypeAppiPad105)
	iPad97 := ScreenshotTypeiPad97.APIValue()
	assert.Equal(t, iPad97, asc.ScreenshotDisplayTypeAppiPad97)
	iPadPro129 := ScreenshotTypeiPadPro129.APIValue()
	assert.Equal(t, iPadPro129, asc.ScreenshotDisplayTypeAppiPadPro129)
	iPadPro3Gen11 := ScreenshotTypeiPadPro3Gen11.APIValue()
	assert.Equal(t, iPadPro3Gen11, asc.ScreenshotDisplayTypeAppiPadPro3Gen11)
	iPadPro3Gen129 := ScreenshotTypeiPadPro3Gen129.APIValue()
	assert.Equal(t, iPadPro3Gen129, asc.ScreenshotDisplayTypeAppiPadPro3Gen129)
	iPhone35 := ScreenshotTypeiPhone35.APIValue()
	assert.Equal(t, iPhone35, asc.ScreenshotDisplayTypeAppiPhone35)
	iPhone40 := ScreenshotTypeiPhone40.APIValue()
	assert.Equal(t, iPhone40, asc.ScreenshotDisplayTypeAppiPhone40)
	iPhone47 := ScreenshotTypeiPhone47.APIValue()
	assert.Equal(t, iPhone47, asc.ScreenshotDisplayTypeAppiPhone47)
	iPhone55 := ScreenshotTypeiPhone55.APIValue()
	assert.Equal(t, iPhone55, asc.ScreenshotDisplayTypeAppiPhone55)
	iPhone58 := ScreenshotTypeiPhone58.APIValue()
	assert.Equal(t, iPhone58, asc.ScreenshotDisplayTypeAppiPhone58)
	iPhone65 := ScreenshotTypeiPhone65.APIValue()
	assert.Equal(t, iPhone65, asc.ScreenshotDisplayTypeAppiPhone65)
	watchSeries3 := ScreenshotTypeWatchSeries3.APIValue()
	assert.Equal(t, watchSeries3, asc.ScreenshotDisplayTypeAppWatchSeries3)
	watchSeries4 := ScreenshotTypeWatchSeries4.APIValue()
	assert.Equal(t, watchSeries4, asc.ScreenshotDisplayTypeAppWatchSeries4)
	iMessageiPad105 := ScreenshotTypeiMessageiPad105.APIValue()
	assert.Equal(t, iMessageiPad105, asc.ScreenshotDisplayTypeiMessageAppIPad105)
	iMessageiPad97 := ScreenshotTypeiMessageiPad97.APIValue()
	assert.Equal(t, iMessageiPad97, asc.ScreenshotDisplayTypeiMessageAppIPad97)
	iMessageiPadPro129 := ScreenshotTypeiMessageiPadPro129.APIValue()
	assert.Equal(t, iMessageiPadPro129, asc.ScreenshotDisplayTypeiMessageAppIPadPro129)
	iMessageiPadPro3Gen11 := ScreenshotTypeiMessageiPadPro3Gen11.APIValue()
	assert.Equal(t, iMessageiPadPro3Gen11, asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen11)
	iMessageiPadPro3Gen129 := ScreenshotTypeiMessageiPadPro3Gen129.APIValue()
	assert.Equal(t, iMessageiPadPro3Gen129, asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen129)
	iMessageiPhone40 := ScreenshotTypeiMessageiPhone40.APIValue()
	assert.Equal(t, iMessageiPhone40, asc.ScreenshotDisplayTypeiMessageAppIPhone40)
	iMessageiPhone47 := ScreenshotTypeiMessageiPhone47.APIValue()
	assert.Equal(t, iMessageiPhone47, asc.ScreenshotDisplayTypeiMessageAppIPhone47)
	iMessageiPhone55 := ScreenshotTypeiMessageiPhone55.APIValue()
	assert.Equal(t, iMessageiPhone55, asc.ScreenshotDisplayTypeiMessageAppIPhone55)
	iMessageiPhone58 := ScreenshotTypeiMessageiPhone58.APIValue()
	assert.Equal(t, iMessageiPhone58, asc.ScreenshotDisplayTypeiMessageAppIPhone58)
	iMessageiPhone65 := ScreenshotTypeiMessageiPhone65.APIValue()
	assert.Equal(t, iMessageiPhone65, asc.ScreenshotDisplayTypeiMessageAppIPhone65)
	bad := screenshotType("Google Pixel").APIValue()
	assert.Empty(t, bad)
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
