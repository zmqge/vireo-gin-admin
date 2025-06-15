package utils

import (
	"github.com/mssola/user_agent"
)

func ParseUserAgent(userAgent string) (device, os, browser string) {
	ua := user_agent.New(userAgent)

	// 获取浏览器信息
	browser, _ = ua.Browser()

	// 获取操作系统信息
	os = ua.OS()

	// 判断设备类型
	if ua.Mobile() {
		device = "Mobile"
	} else if false {
		device = "Tablet"
	} else if ua.Bot() {
		device = "Bot"
	} else {
		device = "Desktop"
	}

	return
}
