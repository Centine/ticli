package cmd

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/centine/ticli/internal/checks"
	"github.com/spf13/cobra"
)

var platformFlag string

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

var workstationCmd = &cobra.Command{
	Use:   "workstation",
	Short: "Runs platform specific submodules to validate that the local development environment is properly set up",
	Long: `This command validates the local development environment setup, 
        checks for the necessary user permissions and software installations.`,
}

var wsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks that prerequisites are met for signing into cloud platforms",
	Run: func(cmd *cobra.Command, args []string) {
		platform := runtime.GOOS
		fmt.Println("Detected platform: " + platform)
		if platformFlag == "macos" {
			platform = "darwin"
		}
		if (platform == "windows") && isWSL() {
			platform = "wsl"
		}
		if platformFlag != "" {
			platform = platformFlag
		}
		switch platform {
		case "windows":
			fmt.Println("Native Windows platform set")
			checks.PerformChecks()
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
	},
}

func init() {
	rootCmd.AddCommand(workstationCmd)
	workstationCmd.AddCommand(wsCheckCmd)
	wsCheckCmd.Flags().StringVar(&platformFlag, "platform", "", "Override platform detection")
}
