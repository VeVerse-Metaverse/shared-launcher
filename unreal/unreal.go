package unreal

import (
	"fmt"
	"runtime"
)

func GetPlatformName() (string, error) {
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		return "Win64", nil
	}

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		return "Mac", nil
	}

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "linux" {
		return "Linux", nil
	}

	return "", fmt.Errorf("unsupported platform: %v", runtime.GOOS)
}

func GetEnvironmentConfiguration(env string) string {
	if env == "dev" {
		return "Development"
	} else if env == "test" {
		return "Test"
	} else if env == "prod" {
		return "Shipping"
	}
	return "Shipping"
}
