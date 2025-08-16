package utils

import (
	"github.com/mssola/user_agent"
)

// ParsedUserAgent menampung hasil parsing dari User-Agent string.
type ParsedUserAgent struct {
	BrowserName    string
	BrowserVersion string
	OSName         string
	OSVersion      string
	DeviceType     string
}

// ParseUserAgent mem-parsing user agent string dan mengembalikan data terstruktur.
func ParseUserAgent(uaString string) *ParsedUserAgent {
	ua := user_agent.New(uaString)

	browserName, browserVersion := ua.Browser()
	osInfo := ua.OSInfo()

	deviceType := "unknown"
	if ua.Mobile() {
		deviceType = "mobile"
	} else if !ua.Bot() {
		deviceType = "desktop"
	}

	return &ParsedUserAgent{
		BrowserName:    browserName,
		BrowserVersion: browserVersion,
		OSName:         osInfo.Name,
		OSVersion:      osInfo.Version,
		DeviceType:     deviceType,
	}
}
