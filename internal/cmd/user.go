package cmd

import "github.com/spf13/cobra"

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "ユーザー情報の操作",
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "認証ユーザーの情報を取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get("/user/me", nil)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

func init() {
	userCmd.AddCommand(userMeCmd)
}
