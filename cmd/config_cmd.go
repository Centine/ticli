// in ticli.go
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/centine/ticli/internal/config"
	"github.com/spf13/cobra"
)

var configRootCmd = &cobra.Command{
	Use:   "config",
	Short: "Ticli config",
	Long:  `All things related to configuration of the ticli utility`,
}

// showConfigCmd represents the config command
var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Prints the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		value := ctx.Value(KeyTicliContext)
		// debug if there are problems understanding how command context works
		switch v := value.(type) {
		case config.TicliContext:
			// DBG fmt.Println("All OK, value is of type TicliContext")
		case nil:
			panic("No value in the context for key 'config'")
		default:
			msg := fmt.Sprintf("Unexpected context type %T\n", v)
			panic(msg)

		}
		// end debug
		cfg, ok := ctx.Value(KeyTicliContext).(config.TicliContext)
		if !ok {
			log.Panic("Config not found in context")
			return
		}
		cfgJson, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			fmt.Println("Failed to marshal config:", err)
			return
		}
		fmt.Println(string(cfgJson))
	},
}

var refreshConfigCmd = &cobra.Command{
	Use:   "Refresh",
	Short: "Refresh the configuration from source NOT IMPLEMENTED",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Download and verify new config
		panic("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(configRootCmd)
	configRootCmd.AddCommand(showConfigCmd)
	configRootCmd.AddCommand(refreshConfigCmd)
}
