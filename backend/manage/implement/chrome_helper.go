package implement

import (
	"tiktok-wails/backend/global"

	"github.com/chromedp/chromedp"
)

// DefaultChromeOpts trả về chromedp options chung cho tất cả browser automation.
// Tránh copy-paste options ở nhiều nơi.
func DefaultChromeOpts(userDataDir string, headless bool) []chromedp.ExecAllocatorOption {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(global.PathAppChrome),
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("start-maximized", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.UserDataDir(userDataDir),
	)
	return opts
}

// HeadlessChromeOpts trả về options cho headless mode (dùng khi scrape Douyin).
func HeadlessChromeOpts(userDataDir string) []chromedp.ExecAllocatorOption {
	opts := DefaultChromeOpts(userDataDir, true)
	opts = append(opts,
		chromedp.Flag("enable-logging", true),
		chromedp.Flag("log-level", "0"),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.Flag("disable-ipc-flooding-protection", true),
	)
	return opts
}

// VisibleChromeOpts trả về options cho visible mode (dùng khi upload TikTok).
func VisibleChromeOpts(userDataDir string) []chromedp.ExecAllocatorOption {
	opts := DefaultChromeOpts(userDataDir, false)
	opts = append(opts,
		chromedp.Flag("lang", "vi-VN"),
		chromedp.Flag("accept-language", "vi-VN,vi;q=0.9,en;q=0.8"),
	)
	return opts
}
