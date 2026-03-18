package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	contactGroupTrashed bool

	cgRecipientName  string
	cgRecipientTitle string
	cgRecipientRuby  string
	cgTaxOption      string
	cgPeppolID       string
)

var contactGroupCmd = &cobra.Command{
	Use:     "contact-group",
	Aliases: []string{"contact-groups", "cg"},
	Short:   "取引先グループの操作",
}

var contactGroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "取引先グループの一覧を取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		q := url.Values{}
		if contactGroupTrashed {
			q.Set("trashed", "true")
		}
		data, _, err := c.Get("/contact_groups", q)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactGroupGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "取引先グループの詳細を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get(fmt.Sprintf("/contact_group/%s", args[0]), nil)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactGroupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "取引先グループを作成",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if cgRecipientName != "" {
			flags["recipient_name"] = cgRecipientName
		}
		if cgRecipientTitle != "" {
			flags["recipient_title"] = cgRecipientTitle
		}
		if cgRecipientRuby != "" {
			flags["recipient_ruby"] = cgRecipientRuby
		}
		if cgTaxOption != "" {
			flags["tax_option"] = cgTaxOption
		}
		if cgPeppolID != "" {
			flags["peppol_id"] = cgPeppolID
		}
		body := buildBody(flags)
		data, _, err := c.Post("/contact_group", body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactGroupHideCmd = &cobra.Command{
	Use:   "hide <id>",
	Short: "取引先グループを非表示にする",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Put(fmt.Sprintf("/contact_group/%s/trashed", args[0]), "")
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactGroupRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "取引先グループを再表示する",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Delete(fmt.Sprintf("/contact_group/%s/trashed", args[0]))
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

func init() {
	// list
	contactGroupListCmd.Flags().BoolVar(&contactGroupTrashed, "trashed", false, "非表示の取引先グループを取得")

	// create
	addJSONFlag(contactGroupCreateCmd)
	contactGroupCreateCmd.Flags().StringVar(&cgRecipientName, "name", "", "取引先名 (必須)")
	contactGroupCreateCmd.Flags().StringVar(&cgRecipientTitle, "title", "", "敬称 (御中/様)")
	contactGroupCreateCmd.Flags().StringVar(&cgRecipientRuby, "ruby", "", "フリガナ")
	contactGroupCreateCmd.Flags().StringVar(&cgTaxOption, "tax-option", "", "税表示 (USE_SENDER/INCLUDE/EXCLUDE/EXEMPT/NONENTRY)")
	contactGroupCreateCmd.Flags().StringVar(&cgPeppolID, "peppol-id", "", "Peppol ID")

	contactGroupCmd.AddCommand(
		contactGroupListCmd,
		contactGroupGetCmd,
		contactGroupCreateCmd,
		contactGroupHideCmd,
		contactGroupRestoreCmd,
	)
}
