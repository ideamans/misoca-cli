package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Flags for invoice commands
var (
	invoiceType           string
	invoiceCondition      string
	invoiceContactGroupID string
	invoicePaymentStatus  string
	invoiceInvoiceStatus  string
	invoiceFrom           string
	invoiceTo             string
	invoiceDueDateFrom    string
	invoiceDueDateTo      string
	invoiceUpdatedFrom    string
	invoiceUpdatedTo      string
	invoiceOrder          string
	invoiceOrderBy        string

	invoiceNumber    string
	invoiceIssueDate string
	invoiceSubject   string
	invoiceDueOn     string
	invoiceContactID int

	invoicePaidOn string
)

var invoiceCmd = &cobra.Command{
	Use:     "invoice",
	Aliases: []string{"invoices"},
	Short:   "請求書の操作",
}

var invoiceListCmd = &cobra.Command{
	Use:   "list",
	Short: "請求書の一覧を取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		q := listQuery()
		for _, kv := range []struct{ flag, param string }{
			{invoiceType, "type"},
			{invoiceCondition, "condition"},
			{invoiceContactGroupID, "contact_group_id"},
			{invoicePaymentStatus, "payment_status"},
			{invoiceInvoiceStatus, "invoice_status"},
			{invoiceFrom, "from"},
			{invoiceTo, "to"},
			{invoiceDueDateFrom, "due_date_from"},
			{invoiceDueDateTo, "due_date_to"},
			{invoiceUpdatedFrom, "updated_at_from"},
			{invoiceUpdatedTo, "updated_at_to"},
			{invoiceOrder, "order"},
			{invoiceOrderBy, "order_by"},
		} {
			if kv.flag != "" {
				q.Set(kv.param, kv.flag)
			}
		}
		data, _, err := c.Get("/invoices", q)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoiceGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "請求書の詳細を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get(fmt.Sprintf("/invoice/%s", args[0]), nil)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoiceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "請求書を作成",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if invoiceContactID != 0 {
			flags["contact_id"] = invoiceContactID
		}
		if invoiceIssueDate != "" {
			flags["issue_date"] = invoiceIssueDate
		}
		if invoiceSubject != "" {
			flags["subject"] = invoiceSubject
		}
		if invoiceDueOn != "" {
			flags["payment_due_on"] = invoiceDueOn
		}
		if invoiceNumber != "" {
			flags["invoice_number"] = invoiceNumber
		}
		body := buildBody(flags)
		data, _, err := c.Post("/invoice", body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoicePDFCmd = &cobra.Command{
	Use:   "pdf <id>",
	Short: "請求書のPDFを取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := c.GetPDF(fmt.Sprintf("/invoice/%s/pdf", args[0]))
		if err != nil {
			return err
		}
		if outputFile == "" {
			outputFile = fmt.Sprintf("invoice_%s.pdf", args[0])
		}
		return savePDF(data, outputFile)
	},
}

var invoiceTrashCmd = &cobra.Command{
	Use:   "trash <id>",
	Short: "請求書をゴミ箱に移動",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Put(fmt.Sprintf("/invoice/%s/trashed", args[0]), "")
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoiceRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "請求書をゴミ箱から復元",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Delete(fmt.Sprintf("/invoice/%s/trashed", args[0]))
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoiceSubmitCmd = &cobra.Command{
	Use:   "submit <id>",
	Short: "請求書を送付済みにする",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Put(fmt.Sprintf("/invoice/%s/submitted", args[0]), "")
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoiceUnsubmitCmd = &cobra.Command{
	Use:   "unsubmit <id>",
	Short: "請求書の送付済みを取り消す",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Delete(fmt.Sprintf("/invoice/%s/submitted", args[0]))
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoicePayCmd = &cobra.Command{
	Use:   "pay <id>",
	Short: "請求書を入金済みにする",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if invoicePaidOn != "" {
			flags["paid_on"] = invoicePaidOn
		}
		body := buildBody(flags)
		data, _, err := c.Put(fmt.Sprintf("/invoice/%s/paid", args[0]), body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoiceUnpayCmd = &cobra.Command{
	Use:   "unpay <id>",
	Short: "請求書の入金済みを取り消す",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Delete(fmt.Sprintf("/invoice/%s/paid", args[0]))
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var invoiceMailCmd = &cobra.Command{
	Use:   "mail <id>",
	Short: "請求書を郵送依頼する",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Post(fmt.Sprintf("/invoice/%s/send_by_postal_mail", args[0]), "")
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

func init() {
	// list
	addListFlags(invoiceListCmd)
	invoiceListCmd.Flags().StringVar(&invoiceType, "type", "", "種類 (active/archived/trashed/untrashed)")
	invoiceListCmd.Flags().StringVar(&invoiceCondition, "condition", "", "キーワード検索 (番号・取引先名・件名・メモ)")
	invoiceListCmd.Flags().StringVar(&invoiceContactGroupID, "contact-group-id", "", "取引先グループIDでフィルタ")
	invoiceListCmd.Flags().StringVar(&invoicePaymentStatus, "payment-status", "", "入金状況 (paid/unpaid)")
	invoiceListCmd.Flags().StringVar(&invoiceInvoiceStatus, "invoice-status", "", "請求状況 (submitted/unsubmitted)")
	invoiceListCmd.Flags().StringVar(&invoiceFrom, "from", "", "請求日の開始日 (YYYY/MM/DD)")
	invoiceListCmd.Flags().StringVar(&invoiceTo, "to", "", "請求日の終了日 (YYYY/MM/DD)")
	invoiceListCmd.Flags().StringVar(&invoiceDueDateFrom, "due-date-from", "", "支払期限の開始日 (YYYY/MM/DD)")
	invoiceListCmd.Flags().StringVar(&invoiceDueDateTo, "due-date-to", "", "支払期限の終了日 (YYYY/MM/DD)")
	invoiceListCmd.Flags().StringVar(&invoiceUpdatedFrom, "updated-from", "", "更新日の開始日 (YYYY/MM/DD)")
	invoiceListCmd.Flags().StringVar(&invoiceUpdatedTo, "updated-to", "", "更新日の終了日 (YYYY/MM/DD)")
	invoiceListCmd.Flags().StringVar(&invoiceOrder, "order", "", "並び順 (asc/desc)")
	invoiceListCmd.Flags().StringVar(&invoiceOrderBy, "order-by", "", "並び順基準 (created_at/updated_at/issue_date/payment_due_on)")

	// create
	addJSONFlag(invoiceCreateCmd)
	invoiceCreateCmd.Flags().IntVar(&invoiceContactID, "contact-id", 0, "送り先ID (必須)")
	invoiceCreateCmd.Flags().StringVar(&invoiceIssueDate, "issue-date", "", "請求日 (必須, YYYY/MM/DD)")
	invoiceCreateCmd.Flags().StringVar(&invoiceSubject, "subject", "", "件名")
	invoiceCreateCmd.Flags().StringVar(&invoiceDueOn, "payment-due-on", "", "支払期限 (YYYY/MM/DD)")
	invoiceCreateCmd.Flags().StringVar(&invoiceNumber, "invoice-number", "", "請求書番号")

	// pdf
	addOutputFlag(invoicePDFCmd)

	// pay
	addJSONFlag(invoicePayCmd)
	invoicePayCmd.Flags().StringVar(&invoicePaidOn, "paid-on", "", "入金日 (YYYY/MM/DD)")

	invoiceCmd.AddCommand(
		invoiceListCmd,
		invoiceGetCmd,
		invoiceCreateCmd,
		invoicePDFCmd,
		invoiceTrashCmd,
		invoiceRestoreCmd,
		invoiceSubmitCmd,
		invoiceUnsubmitCmd,
		invoicePayCmd,
		invoiceUnpayCmd,
		invoiceMailCmd,
	)
}
