package cmd

import "github.com/spf13/cobra"

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "misoca",
	Short:   "Misoca API CLI",
	Long:    "Misoca APIを操作するためのコマンドラインツール",
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(
		authCmd,
		invoiceCmd,
		estimateCmd,
		deliverySlipCmd,
		contactCmd,
		contactGroupCmd,
		itemCmd,
		userCmd,
	)
}
