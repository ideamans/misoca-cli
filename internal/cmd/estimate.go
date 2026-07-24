package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	estimateType           string
	estimateContactGroupID string

	estimateContactID  int
	estimateIssueDate  string
	estimateSubject    string
	estimateExpireDate string
	estimateNumber     string

	distMailSubject   string
	distMailBody      string
	distIncludeSelfCC bool
)

var estimateCmd = &cobra.Command{
	Use:     "estimate",
	Aliases: []string{"estimates"},
	Short:   "見積書の操作",
}

var estimateListCmd = &cobra.Command{
	Use:   "list",
	Short: "見積書の一覧を取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		q := listQuery()
		if estimateType != "" {
			q.Set("type", estimateType)
		}
		if estimateContactGroupID != "" {
			q.Set("contact_group_id", estimateContactGroupID)
		}
		data, _, err := c.Get("/estimates", q)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var estimateGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "見積書の詳細を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get(fmt.Sprintf("/estimate/%s", args[0]), nil)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var estimateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "見積書を作成",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if estimateContactID != 0 {
			flags["contact_id"] = estimateContactID
		}
		if estimateIssueDate != "" {
			flags["issue_date"] = estimateIssueDate
		}
		if estimateSubject != "" {
			flags["subject"] = estimateSubject
		}
		if estimateExpireDate != "" {
			flags["expire_date"] = estimateExpireDate
		}
		if estimateNumber != "" {
			flags["estimate_number"] = estimateNumber
		}
		body := buildBody(flags)
		data, _, err := c.Post("/estimate", body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var estimatePDFCmd = &cobra.Command{
	Use:   "pdf <id>",
	Short: "見積書のPDFを取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := c.GetPDF(fmt.Sprintf("/estimate/%s/pdf", args[0]))
		if err != nil {
			return err
		}
		if outputFile == "" {
			outputFile = fmt.Sprintf("estimate_%s.pdf", args[0])
		}
		return savePDF(data, outputFile)
	},
}

var estimateLogoCmd = &cobra.Command{
	Use:   "logo <id>",
	Short: "見積書のロゴ画像を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, ct, err := c.GetBinary(fmt.Sprintf("/estimate/%s/logo", args[0]), "")
		if err != nil {
			return err
		}
		if data == nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "ロゴ画像が設定されていません")
			return nil
		}
		if outputFile == "" {
			outputFile = fmt.Sprintf("estimate_%s_logo%s", args[0], extFromContentType(ct))
		}
		return saveFile(data, outputFile)
	},
}

var estimateStampCmd = &cobra.Command{
	Use:   "stamp <id>",
	Short: "見積書の印影画像を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, ct, err := c.GetBinary(fmt.Sprintf("/estimate/%s/stamp", args[0]), "")
		if err != nil {
			return err
		}
		if data == nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "印影画像が設定されていません")
			return nil
		}
		if outputFile == "" {
			outputFile = fmt.Sprintf("estimate_%s_stamp%s", args[0], extFromContentType(ct))
		}
		return saveFile(data, outputFile)
	},
}

var estimateDistributeCmd = &cobra.Command{
	Use:   "distribute <id>",
	Short: "見積書をメールで送信",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if distMailSubject != "" {
			flags["mail_subject"] = distMailSubject
		}
		if distMailBody != "" {
			flags["mail_body"] = distMailBody
		}
		if distIncludeSelfCC {
			flags["including_self_to_cc"] = true
		}
		body := buildBody(flags)
		data, _, err := c.Post(fmt.Sprintf("/estimate/%s/distribute", args[0]), body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

func init() {
	// list
	addListFlags(estimateListCmd)
	estimateListCmd.Flags().StringVar(&estimateType, "type", "", "種類 (active/archived/trashed/untrashed)")
	estimateListCmd.Flags().StringVar(&estimateContactGroupID, "contact-group-id", "", "取引先グループIDでフィルタ")

	// create
	addJSONFlag(estimateCreateCmd)
	estimateCreateCmd.Flags().IntVar(&estimateContactID, "contact-id", 0, "送り先ID (必須)")
	estimateCreateCmd.Flags().StringVar(&estimateIssueDate, "issue-date", "", "見積日 (必須, YYYY/MM/DD)")
	estimateCreateCmd.Flags().StringVar(&estimateSubject, "subject", "", "件名")
	estimateCreateCmd.Flags().StringVar(&estimateExpireDate, "expire-date", "", "有効期限 (YYYY/MM/DD)")
	estimateCreateCmd.Flags().StringVar(&estimateNumber, "estimate-number", "", "見積書番号")

	// pdf / logo / stamp
	addOutputFlag(estimatePDFCmd)
	addOutputFlag(estimateLogoCmd)
	addOutputFlag(estimateStampCmd)

	// distribute
	addJSONFlag(estimateDistributeCmd)
	estimateDistributeCmd.Flags().StringVar(&distMailSubject, "mail-subject", "", "メール件名 (必須)")
	estimateDistributeCmd.Flags().StringVar(&distMailBody, "mail-body", "", "メール本文")
	estimateDistributeCmd.Flags().BoolVar(&distIncludeSelfCC, "cc-self", false, "自分をCCに含める")

	estimateCmd.AddCommand(
		estimateListCmd,
		estimateGetCmd,
		estimateCreateCmd,
		estimatePDFCmd,
		estimateLogoCmd,
		estimateStampCmd,
		estimateDistributeCmd,
	)
}
