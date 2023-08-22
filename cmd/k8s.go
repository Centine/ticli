package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes related operations",
	Long:  `All things related to Kubernetes`,
}

var validateLocalCmd = &cobra.Command{
	Use:   "validatelocal",
	Short: "Validate local Kubernetes",
	Long:  `Validate Kubernetes configuration on local system`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validatelocal invoked")
	},
}

var validateRemoteCmd = &cobra.Command{
	Use:   "validateremote",
	Short: "Validate remote Kubernetes",
	Long:  `Validate Kubernetes configuration on remote system`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validateremote invoked")
	},
}

var getKubeCfgCmd = &cobra.Command{
	Use:   "get-kubecfg",
	Short: "Get Kubernetes config",
	Long:  `Fetch the Kubernetes configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get-kubecfg invoked")
	},
}

func init() {
	rootCmd.AddCommand(k8sCmd)
	k8sCmd.AddCommand(validateLocalCmd)
	k8sCmd.AddCommand(validateRemoteCmd)
	k8sCmd.AddCommand(getKubeCfgCmd)
}
