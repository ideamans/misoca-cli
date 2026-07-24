package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	contactGroupIDFilter string
	contactTrashed       bool

	contactRecipientName  string
	contactRecipientTitle string
	contactRecipientRuby  string
	contactTaxOption      string
	contactCode           string
	contactNotes          string
	contactMail           string
	contactMailCC         string
	contactZipCode        string
	contactAddress1       string
	contactAddress2       string
	contactName1          string
	contactName2          string
	contactName3          string
	contactName4          string
	contactTitle2         string
	contactTelNo          string
	contactFaxNo          string
	contactMemo           string
	contactCreateGroupID  int
)

var contactCmd = &cobra.Command{
	Use:     "contact",
	Aliases: []string{"contacts"},
	Short:   "送り先の操作",
}

var contactListCmd = &cobra.Command{
	Use:   "list",
	Short: "送り先の一覧を取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		q := url.Values{}
		if contactGroupIDFilter != "" {
			q.Set("contact_group_id", contactGroupIDFilter)
		}
		if contactTrashed {
			q.Set("trashed", "true")
		}
		data, _, err := c.Get("/contacts", q)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "送り先の詳細を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get(fmt.Sprintf("/contact/%s", args[0]), nil)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "送り先を作成",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if contactCreateGroupID != 0 {
			flags["contact_group_id"] = contactCreateGroupID
		}
		if contactRecipientName != "" {
			flags["recipient_name"] = contactRecipientName
		}
		if contactRecipientTitle != "" {
			flags["recipient_title"] = contactRecipientTitle
		}
		if contactRecipientRuby != "" {
			flags["recipient_ruby"] = contactRecipientRuby
		}
		if contactTaxOption != "" {
			flags["tax_option"] = contactTaxOption
		}
		if contactCode != "" {
			flags["recipient_code"] = contactCode
		}
		if contactNotes != "" {
			flags["recipient_notes"] = contactNotes
		}
		if contactMail != "" {
			flags["recipient_mail_address"] = contactMail
		}
		if contactMailCC != "" {
			flags["recipient_mail_address_cc"] = contactMailCC
		}
		if contactZipCode != "" {
			flags["recipient_zip_code"] = contactZipCode
		}
		if contactAddress1 != "" {
			flags["recipient_address1"] = contactAddress1
		}
		if contactAddress2 != "" {
			flags["recipient_address2"] = contactAddress2
		}
		if contactName1 != "" {
			flags["recipient_name1"] = contactName1
		}
		if contactName2 != "" {
			flags["recipient_name2"] = contactName2
		}
		if contactName3 != "" {
			flags["recipient_name3"] = contactName3
		}
		if contactName4 != "" {
			flags["recipient_name4"] = contactName4
		}
		if contactTitle2 != "" {
			flags["recipient_title2"] = contactTitle2
		}
		if contactTelNo != "" {
			flags["recipient_tel_no"] = contactTelNo
		}
		if contactFaxNo != "" {
			flags["recipient_fax_no"] = contactFaxNo
		}
		if contactMemo != "" {
			flags["memo"] = contactMemo
		}
		body := buildBody(flags)
		data, _, err := c.Post("/contact", body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactHideCmd = &cobra.Command{
	Use:   "hide <id>",
	Short: "送り先を非表示にする",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Put(fmt.Sprintf("/contact/%s/trashed", args[0]), "")
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var contactRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "送り先を再表示する",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Delete(fmt.Sprintf("/contact/%s/trashed", args[0]))
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

func init() {
	// list
	contactListCmd.Flags().StringVar(&contactGroupIDFilter, "contact-group-id", "", "取引先グループIDでフィルタ")
	contactListCmd.Flags().BoolVar(&contactTrashed, "trashed", false, "非表示の送り先を取得")

	// create
	addJSONFlag(contactCreateCmd)
	contactCreateCmd.Flags().IntVar(&contactCreateGroupID, "contact-group-id", 0, "取引先グループID")
	contactCreateCmd.Flags().StringVar(&contactRecipientName, "name", "", "送り先名")
	contactCreateCmd.Flags().StringVar(&contactRecipientTitle, "title", "", "敬称 (御中/様)")
	contactCreateCmd.Flags().StringVar(&contactRecipientRuby, "ruby", "", "フリガナ")
	contactCreateCmd.Flags().StringVar(&contactTaxOption, "tax-option", "", "税表示 (USE_SENDER/INCLUDE/EXCLUDE/EXEMPT/NONENTRY)")
	contactCreateCmd.Flags().StringVar(&contactCode, "code", "", "取引先コード")
	contactCreateCmd.Flags().StringVar(&contactNotes, "notes", "", "備考")
	contactCreateCmd.Flags().StringVar(&contactMail, "mail", "", "メールアドレス")
	contactCreateCmd.Flags().StringVar(&contactMailCC, "mail-cc", "", "CCメールアドレス")
	contactCreateCmd.Flags().StringVar(&contactZipCode, "zip-code", "", "郵便番号")
	contactCreateCmd.Flags().StringVar(&contactAddress1, "address1", "", "住所1")
	contactCreateCmd.Flags().StringVar(&contactAddress2, "address2", "", "住所2")
	contactCreateCmd.Flags().StringVar(&contactName1, "name1", "", "宛名1")
	contactCreateCmd.Flags().StringVar(&contactName2, "name2", "", "宛名2")
	contactCreateCmd.Flags().StringVar(&contactName3, "name3", "", "宛名3")
	contactCreateCmd.Flags().StringVar(&contactName4, "name4", "", "宛名4")
	contactCreateCmd.Flags().StringVar(&contactTitle2, "title2", "", "敬称2")
	contactCreateCmd.Flags().StringVar(&contactTelNo, "tel", "", "電話番号")
	contactCreateCmd.Flags().StringVar(&contactFaxNo, "fax", "", "FAX番号")
	contactCreateCmd.Flags().StringVar(&contactMemo, "memo", "", "社内メモ")

	contactCmd.AddCommand(
		contactListCmd,
		contactGetCmd,
		contactCreateCmd,
		contactHideCmd,
		contactRestoreCmd,
	)
}
