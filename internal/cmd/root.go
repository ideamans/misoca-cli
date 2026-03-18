package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "misoca",
	Short: "Misoca API CLI",
	Long:  "Misoca APIを操作するためのコマンドラインツール",
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
