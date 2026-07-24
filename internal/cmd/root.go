package cmd

import (
	"github.com/ideamans/go-llm-cli-kit/llmcmd"
	"github.com/spf13/cobra"
)

// PluginVersion はこのCLIのリリースバージョンです。
// plugins/misoca-cli/.claude-plugin/plugin.json の version と一致していることを
// テストが、git タグと一致していることをリリースワークフローが検査します。
const PluginVersion = "0.2.0"

// Version はビルド時に上書きされるバージョン文字列です。
var Version = PluginVersion

var rootCmd = &cobra.Command{
	Use:     "misoca",
	Short:   "Misoca API CLI",
	Long:    "Misoca APIを操作するためのコマンドラインツール",
	Version: Version,
}

func Execute() error {
	rootCmd.Version = Version
	return rootCmd.Execute()
}

// Root は組み立て済みのコマンドツリーを実行せずに返します。
// カタログ生成器が Execute と同じ定義から生成するために使います。
func Root() *cobra.Command { return rootCmd }

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
	llmcmd.AddTo(rootCmd, LLMConfig())
}
