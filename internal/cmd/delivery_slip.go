package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dsType           string
	dsContactGroupID string

	dsContactID    int
	dsIssueDate    string
	dsSubject      string
	dsDeliveryDate string
	dsNumber       string
)

var deliverySlipCmd = &cobra.Command{
	Use:     "delivery-slip",
	Aliases: []string{"delivery-slips", "ds"},
	Short:   "納品書の操作",
}

var deliverySlipListCmd = &cobra.Command{
	Use:   "list",
	Short: "納品書の一覧を取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		q := listQuery()
		if dsType != "" {
			q.Set("type", dsType)
		}
		if dsContactGroupID != "" {
			q.Set("contact_group_id", dsContactGroupID)
		}
		data, _, err := c.Get("/delivery_slips", q)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var deliverySlipGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "納品書の詳細を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get(fmt.Sprintf("/delivery_slip/%s", args[0]), nil)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var deliverySlipCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "納品書を作成",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if dsContactID != 0 {
			flags["contact_id"] = dsContactID
		}
		if dsIssueDate != "" {
			flags["issue_date"] = dsIssueDate
		}
		if dsSubject != "" {
			flags["subject"] = dsSubject
		}
		if dsDeliveryDate != "" {
			flags["delivery_date"] = dsDeliveryDate
		}
		if dsNumber != "" {
			flags["delivery_slip_number"] = dsNumber
		}
		body := buildBody(flags)
		data, _, err := c.Post("/delivery_slip", body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var deliverySlipPDFCmd = &cobra.Command{
	Use:   "pdf <id>",
	Short: "納品書のPDFを取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := c.GetPDF(fmt.Sprintf("/delivery_slip/%s/pdf", args[0]))
		if err != nil {
			return err
		}
		if outputFile == "" {
			outputFile = fmt.Sprintf("delivery_slip_%s.pdf", args[0])
		}
		return savePDF(data, outputFile)
	},
}

func init() {
	// list
	addListFlags(deliverySlipListCmd)
	deliverySlipListCmd.Flags().StringVar(&dsType, "type", "", "種類 (active/archived/trashed/untrashed)")
	deliverySlipListCmd.Flags().StringVar(&dsContactGroupID, "contact-group-id", "", "取引先グループIDでフィルタ")

	// create
	addJSONFlag(deliverySlipCreateCmd)
	deliverySlipCreateCmd.Flags().IntVar(&dsContactID, "contact-id", 0, "送り先ID (必須)")
	deliverySlipCreateCmd.Flags().StringVar(&dsIssueDate, "issue-date", "", "発行日 (必須, YYYY/MM/DD)")
	deliverySlipCreateCmd.Flags().StringVar(&dsSubject, "subject", "", "件名")
	deliverySlipCreateCmd.Flags().StringVar(&dsDeliveryDate, "delivery-date", "", "納品日 (YYYY/MM/DD)")
	deliverySlipCreateCmd.Flags().StringVar(&dsNumber, "delivery-slip-number", "", "納品書番号")

	// pdf
	addOutputFlag(deliverySlipPDFCmd)

	deliverySlipCmd.AddCommand(
		deliverySlipListCmd,
		deliverySlipGetCmd,
		deliverySlipCreateCmd,
		deliverySlipPDFCmd,
	)
}
