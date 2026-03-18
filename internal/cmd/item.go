package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	itemName      string
	itemUnitPrice float64
	itemUnitName  string
	itemTaxType   string
	itemExcludeWT bool
)

var itemCmd = &cobra.Command{
	Use:     "item",
	Aliases: []string{"items"},
	Short:   "品目の操作",
}

var itemListCmd = &cobra.Command{
	Use:   "list",
	Short: "品目の一覧を取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get("/dealing_items", listQuery())
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var itemGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "品目の詳細を取得",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, _, err := c.Get(fmt.Sprintf("/dealing_item/%s", args[0]), nil)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

var itemCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "品目を作成",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		flags := map[string]any{}
		if itemName != "" {
			flags["name"] = itemName
		}
		if itemUnitPrice != 0 {
			flags["unit_price"] = itemUnitPrice
		}
		if itemUnitName != "" {
			flags["unit_name"] = itemUnitName
		}
		if itemTaxType != "" {
			flags["tax_type"] = itemTaxType
		}
		if itemExcludeWT {
			flags["excluding_withholding_tax"] = true
		}
		body := buildBody(flags)
		data, _, err := c.Post("/dealing_item", body)
		if err != nil {
			return err
		}
		printJSON(data)
		return nil
	},
}

func init() {
	// list
	addListFlags(itemListCmd)

	// create
	addJSONFlag(itemCreateCmd)
	itemCreateCmd.Flags().StringVar(&itemName, "name", "", "品目名 (必須)")
	itemCreateCmd.Flags().Float64Var(&itemUnitPrice, "unit-price", 0, "単価")
	itemCreateCmd.Flags().StringVar(&itemUnitName, "unit-name", "", "単位名")
	itemCreateCmd.Flags().StringVar(&itemTaxType, "tax-type", "", "税区分 (USE_SENDER/STANDARD_TAX_10/REDUCED_TAX_8/EXEMPTED_TAX等)")
	itemCreateCmd.Flags().BoolVar(&itemExcludeWT, "exclude-withholding-tax", false, "源泉徴収対象外")

	itemCmd.AddCommand(
		itemListCmd,
		itemGetCmd,
		itemCreateCmd,
	)
}
