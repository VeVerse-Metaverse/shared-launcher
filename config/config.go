package config

import "os"

func GetApi2Url() string {
	if os.Getenv("VE_API_URL") != "" {
		return os.Getenv("VE_API_URL")
	}
	return ApiUrl
}

func GetLauncherId() string {
	return LauncherId
}
