package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ideamans/misoca-cli/internal/api"
	"github.com/spf13/cobra"
)

// Global flags
var (
	outputFile string
	jsonBody   string
	page       int
	perPage    int
)

func newClient() (*api.Client, error) {
	return api.NewClient()
}

func printJSON(data []byte) {
	pretty, _ := api.PrettyJSON(data)
	fmt.Println(pretty)
}

func savePDF(data []byte, path string) error {
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("ファイルの保存に失敗: %w", err)
	}
	fmt.Fprintf(os.Stderr, "保存しました: %s (%d bytes)\n", path, len(data))
	return nil
}

func addListFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&page, "page", 0, "ページ番号")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "1ページあたりの件数 (最大100)")
}

func listQuery() url.Values {
	q := url.Values{}
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		q.Set("per_page", fmt.Sprintf("%d", perPage))
	}
	return q
}

func addJSONFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&jsonBody, "json", "", "リクエストボディ (JSON文字列、フラグと併用可)")
}

func addOutputFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "出力ファイルパス")
}

// buildBody merges individual flag values into --json.
// Flag values take precedence over --json fields.
func buildBody(flags map[string]any) string {
	base := map[string]any{}

	// Parse --json first as base
	if jsonBody != "" {
		json.Unmarshal([]byte(jsonBody), &base)
	}

	// Overlay individual flags
	for k, v := range flags {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case string:
			if val != "" {
				base[k] = val
			}
		case int:
			if val != 0 {
				base[k] = val
			}
		case float64:
			if val != 0 {
				base[k] = val
			}
		case bool:
			if val {
				base[k] = val
			}
		default:
			base[k] = val
		}
	}

	if len(base) == 0 {
		return ""
	}

	b, _ := json.Marshal(base)
	return string(b)
}

func saveFile(data []byte, path string) error {
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("ファイルの保存に失敗: %w", err)
	}
	fmt.Fprintf(os.Stderr, "保存しました: %s (%d bytes)\n", path, len(data))
	return nil
}

func extFromContentType(ct string) string {
	ct = strings.ToLower(ct)
	switch {
	case strings.Contains(ct, "png"):
		return ".png"
	case strings.Contains(ct, "jpeg"), strings.Contains(ct, "jpg"):
		return ".jpg"
	case strings.Contains(ct, "gif"):
		return ".gif"
	case strings.Contains(ct, "svg"):
		return ".svg"
	case strings.Contains(ct, "pdf"):
		return ".pdf"
	default:
		return ".bin"
	}
}
