package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ticli",
	Short: "ticli is a utility tool for developers",
	Long:  `ticli is a utility tool for developers that helps execute various tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ticli command invoked")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
