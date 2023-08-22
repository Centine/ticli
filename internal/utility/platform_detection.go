package utility

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Check if the platform demands is compatible with the actual platform
func is_platform_compatible(wanted string, actual string) bool {
	return (wanted == "*" || actual == "*") || (strings.ToLower(wanted) == strings.ToLower(actual))
}

// Checks if it's running on WSL
func isWSL() bool {
	if _, err := os.Stat("/proc/version"); err == nil {
		f, err := os.Open("/proc/version")
		if err != nil {
			return false
		}
		defer f.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(f)

		return strings.Contains(buf.String(), "microsoft")
	}
	return false
}

func detectPlatform(platformOverride string) string {
	platform := runtime.GOOS
	fmt.Println("Detected platform: " + platform)
	if platformOverride == "macos" {
		platform = "darwin"
	}
	if (platform == "windows") && isWSL() {
		platform = "wsl"
	}
	if platformOverride != "" {
		platform = platformOverride
	}
	switch platform {
	case "windows":
		fmt.Println("Native Windows platform set")
	case "wsl":
		fmt.Println("Windows WSL environment set")
		// Call WSL-specific function here
	case "linux":
		fmt.Println("Linux platform set")
		// Call linux-specific function here
	case "darwin":
		fmt.Println("macOS platform set")
		// Call macOS-specific function here
	default:
		fmt.Println("platform not supported")
	}

	return platform
}
