package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var awsAuthCmd = &cobra.Command{
	Use:   "aws-auth",
	Short: "Handles AWS Authentication",
	Long:  `This command helps with AWS Authentication tasks.`,
}

var awsAuthCheck = &cobra.Command{
	Use:   "check",
	Short: "Checks that prerequisites are met for signing into cloud platforms",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws-auth check command invoked")
		// You would add your functionality here
	},
}

var awsAuthLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Executes the check command, and runs an external tool with the necessary parameters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws-auth login command invoked")
		// You would add your functionality here
	},
}

func init() {
	rootCmd.AddCommand(awsAuthCmd)
	awsAuthCmd.AddCommand(awsAuthCheck, awsAuthLoginCmd)
}
