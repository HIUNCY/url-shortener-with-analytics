package utils

import (
	"github.com/mssola/user_agent"
)

type ParsedUserAgent struct {
	BrowserName    string
	BrowserVersion string
	OSName         string
	OSVersion      string
	DeviceType     string
}

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
