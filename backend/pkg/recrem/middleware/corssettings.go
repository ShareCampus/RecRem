package middleware

import "os"

const (
	EnvironmentString string = "IOT_DRIVER_COPILOT_ENVIRONMENT"
	EnvironmentGlobal string = "global"
	EnvironmentDev    string = "dev"
)

var (
	// This should allow us to accept request from "https://test.shifu.dev",
	// "https://chuangyeshitang.staging.shifu.dev" and "https://prod.shifu.cloud"
	GlobalCORSList = []string{
		"https://copilot.test.shifu.dev",
		"https://copilot.staging.shifu.dev",
		"https://copilot.shifu.dev",
	}
	// This should allow us to accept request from "https://*.test.shifu.cloud",
	// "https://*.staging.shifu.cloud", "https://test.shifu.dev",
	// "https://chuangyeshitang.staging.shifu.dev", "https://prod.shifu.cloud",
	// "https://*.vercel.app" and local development
	DevCORSList = []string{
		"http://localhost:*",
		"*.test.shifu.dev",
		"*.vercel.app",
	}
)

// A helper package for configuring CORS settings
func GetCORSSettings() []string {
	// Set CORS settings
	switch env := os.Getenv(EnvironmentString); env {
	case EnvironmentGlobal:
		return GlobalCORSList
	case EnvironmentDev:
		return DevCORSList
	default:
		return DevCORSList
	}
}
