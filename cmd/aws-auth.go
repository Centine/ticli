package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var awsAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Handles authentication",
	Long:  `This command helps with authentication tasks.`,
}

var awsAuthCheck = &cobra.Command{
	Use:   "aws-check",
	Short: "Checks that prerequisites are met for signing into AWS cloud platforms",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws-auth check command invoked")

		// Not sure what to do here.... validate AD credentials?
	},
}

var awsAuthFix = &cobra.Command{
	Use:   "aws-cli-fix",
	Short: "Attempts to fix common issues with AWS Authentication (go-aws-sso)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws-auth fix command invoked")
		// You would add your functionality here
	},
}

var awsAuthLoginCmd = &cobra.Command{
	Use:   "show-aws-login",
	Short: "Shows the necessary command to log into AWS",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To login to AWS, run the following command:")
		fmt.Println(" go-aws-sso login --profile saml --persist --start-url https://dfds.awsapps.com/start --region eu-west-1")
		fmt.Println()
		fmt.Println("A browser window will open, and you will be prompted to login with your AD credentials.")
		fmt.Println("After successful login, you will be able to use the AWS CLI with the 'saml' profile.")
		fmt.Println()

	},
}

func init() {
	rootCmd.AddCommand(awsAuthCmd)
	awsAuthCmd.AddCommand(awsAuthCheck, awsAuthLoginCmd, awsAuthFix)
}
