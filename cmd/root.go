package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/centine/ticli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ticli",
	Short: "ticli is a utility tool for developers",
	Long:  `ticli is a utility tool for developers that helps execute various tasks.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// debug that context set correctly
		ctx := cmd.Context()
		value := ctx.Value(KeyTicliContext)
		switch v := value.(type) {
		case config.TicliContext:
			// DBG fmt.Println("All OK, value is of type TicliContext")
		case nil:
			return errors.New("no value in the context for key 'config'")
		default:
			return fmt.Errorf("unexpected context type %T", v)
		}
		ticliCtx := ctx.Value(KeyTicliContext).(config.TicliContext)
		if ticliCtx.ConfigOrigin == "" {
			return errors.New("config origin not set")
		}
		if ticliCtx.TicliDir == "" {
			return errors.New("config TicliDir not set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

type keyType string

const (
	KeyTicliContext keyType = "TicliContext"
)

func Execute() {
	ticliDir, err := config.SetupTicliEnv()
	if err != nil {
		log.Fatalf("Error setting up ticli environment: %v", err)
	}
	configOrigin, cfg := config.LoadConfig()
	ticliContext := config.TicliContext{
		TicliDir:     ticliDir,
		ConfigOrigin: configOrigin,
		Config:       cfg,
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyTicliContext, ticliContext)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.Println("Initializing root command")
	viper.SetEnvPrefix("TICLI")
	viper.AutomaticEnv()
}
