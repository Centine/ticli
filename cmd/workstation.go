package cmd

import (
	"log"

	"github.com/centine/ticli/internal/checks"
	"github.com/centine/ticli/internal/config"
	"github.com/spf13/cobra"
)

var platformFlag string

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
		cmdCtx := cmd.Context()
		ticliCtx, ok := cmdCtx.Value(KeyTicliContext).(config.TicliContext)
		if !ok {
			log.Panic("TicliContext not found in command context")
			return
		}
		checks.PerformChecks(ticliCtx)
	},
}

func init() {
	rootCmd.AddCommand(workstationCmd)
	workstationCmd.AddCommand(wsCheckCmd)
	wsCheckCmd.Flags().StringVar(&platformFlag, "platform", "", "Override platform detection")
}
