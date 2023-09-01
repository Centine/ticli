package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	ssu_apiclient "github.com/centine/ssu_openapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var selfServiceCmd = &cobra.Command{
	Use:     "selfservice",
	Aliases: []string{"ssu"},
	Short:   "Self service operations",
	Long:    `All things related to self service`,
}

var capabilityCmd = &cobra.Command{
	Use:     "capability",
	Aliases: []string{"cap"},
	Short:   "Capability operations",
	Long:    `Operations related to capabilities`,
}

var getCapabilitiesCmd = &cobra.Command{
	Use:   "get",
	Short: "Get capability",
	Long:  `Fetch information about a specific capability`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("ssu get capability invoked")

		// fmt.Println("Flag auth-token:", cmd.Flag("auth-token").Value.String())
		// fmt.Println("Viper auth-token:", viper.GetString("auth-token"))
		// fmt.Println("Flag cookie:", cmd.Flag("cookie").Value.String())
		// fmt.Println("Viper cookie:", viper.GetString("cookie"))

		configuration := ssu_apiclient.NewConfiguration()
		at := os.Getenv("TICLI_AUTH_TOKEN")
		ct := os.Getenv("TICLI_COOKIE")
		configuration.AddDefaultHeader("Authorization", at)
		configuration.AddDefaultHeader("cookie", ct)
		// configuration.AddDefaultHeader("Authorization", cmd.Flag("auth-token").Value.String())
		// configuration.AddDefaultHeader("cookie", cmd.Flag("cookie").Value.String())
		apiClient := ssu_apiclient.NewAPIClient(configuration)
		resp, httpRes, err := apiClient.CapabilityApi.CapabilitiesGet(context.Background()).Execute()
		if err != nil {
			fmt.Println("Error fetching capabilities: ", err)
			return
		}
		if httpRes.StatusCode != 200 {
			fmt.Println("Error fetching capabilities: ", httpRes.Status)
			return
		}
		fmt.Println("Capabilities:")
		for _, cap := range resp.GetItems() {

			fmt.Printf("Cap %v\n", cap.GetName())
		}

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
var ssAuthCookie string

func init() {
	rootCmd.AddCommand(selfServiceCmd)
	selfServiceCmd.AddCommand(capabilityCmd)
	capabilityCmd.AddCommand(getCapabilitiesCmd)
	capabilityCmd.AddCommand(describeCapabilitiesCmd)

	selfServiceCmd.PersistentFlags().StringVarP(&ssAuthToken, "auth-token", "t", "", "Self service auth token")
	selfServiceCmd.PersistentFlags().StringVarP(&ssAuthCookie, "cookie", "c", "", "Self service cookie token")
	// selfServiceCmd.MarkPersistentFlagRequired("auth-token")
	// selfServiceCmd.MarkPersistentFlagRequired("cookie")

	// Bind flags to Viper configuration
	viper.BindPFlag("auth-token", selfServiceCmd.PersistentFlags().Lookup("auth-token"))
	viper.BindPFlag("cookie", selfServiceCmd.PersistentFlags().Lookup("cookie"))

	getCapabilitiesCmd.Flags().BoolVarP(&getAll, "all", "a", false, "Get all capabilities, not just the ones you have access to")
	getTopicsCmd.Flags().BoolVarP(&getAll, "all", "a", false, "Get all topics, not just the ones you have access to")

	selfServiceCmd.AddCommand(topicsCmd)
	topicsCmd.AddCommand(getTopicsCmd)
	topicsCmd.AddCommand(describeTopicsCmd)
}
