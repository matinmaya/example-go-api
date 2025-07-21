package loghelper

import (
	"fmt"

	"github.com/mssola/user_agent"
)

func ParseUserAgent(uaString string) (device, platform, browser, os string) {
	ua := user_agent.New(uaString)

	name, version := ua.Browser()
	browser = fmt.Sprintf("%s %s", name, version)

	if ua.Mobile() {
		device = "Mobile"
	} else if !ua.Mobile() && !ua.Bot() {
		device = "Desktop"
	} else {
		device = "Unknown"
	}

	platform = ua.Platform()
	os = ua.OS()

	return
}
