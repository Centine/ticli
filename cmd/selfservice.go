package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var selfServiceCmd = &cobra.Command{
	Use:     "selfservice",
	Aliases: []string{"ssu"},
	Short:   "Self service operations",
	Long:    `All things related to self service`,
}

var capabilityCmd = &cobra.Command{
	Use:   "capability",
	Short: "Capability operations",
	Long:  `Operations related to capabilities`,
}

var getCapabilitiesCmd = &cobra.Command{
	Use:   "get",
	Short: "Get capability",
	Long:  `Fetch information about a specific capability`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get invoked")
	},
}

var describeCapabilitiesCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe capability",
	Long:  `Describe a specific capability`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("describe invoked")
	},
}

var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "Topics operations",
	Long:  `Operations related to topics`,
}

var getAll bool

var getTopicsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get topics",
	Long:  `Fetch information about a specific topic or all topics`,
	Run: func(cmd *cobra.Command, args []string) {
		if getAll {
			fmt.Println("get all topics invoked")
		} else {
			fmt.Println("get topic invoked")
		}
	},
}

var describeTopicsCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe topics",
	Long:  `Describe a specific topic`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("describe topic invoked")
	},
}

var ssAuthToken string

func init() {
	rootCmd.AddCommand(selfServiceCmd)
	selfServiceCmd.AddCommand(capabilityCmd)
	capabilityCmd.AddCommand(getCapabilitiesCmd)
	capabilityCmd.AddCommand(describeCapabilitiesCmd)

	selfServiceCmd.PersistentFlags().StringVarP(&ssAuthToken, "auth-token", "t", "", "Self service auth token")
	selfServiceCmd.MarkPersistentFlagRequired("auth-token")

	getCapabilitiesCmd.Flags().BoolVarP(&getAll, "all", "a", false, "Get all capabilities, not just the ones you have access to")
	getTopicsCmd.Flags().BoolVarP(&getAll, "all", "a", false, "Get all topics, not just the ones you have access to")

	selfServiceCmd.AddCommand(topicsCmd)
	topicsCmd.AddCommand(getTopicsCmd)
	topicsCmd.AddCommand(describeTopicsCmd)
}
