package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ideamans/go-llm-cli-kit/llmcmd"

	"github.com/ideamans/misoca-cli/internal/llmdocs"
)

func TestEmbeddedReference(t *testing.T) {
	g, err := llmdocs.Docs().Markdown()
	if err != nil {
		t.Fatalf("Markdown: %v", err)
	}
	for _, want := range []string{
		"AIエージェント向けリファレンス",
		"MISOCA_CLIENT_ID", // 認証
		"--type untrashed", // 一覧の落とし穴
		"郵送",               // 費用の発生する操作
		"帳票と状態",            // 2章目
		"コマンドカタログ",         // 生成物
		"misoca invoice",
	} {
		if !strings.Contains(g, want) {
			t.Errorf("埋め込みリファレンスに %q がありません", want)
		}
	}
}

func TestChapterOrder(t *testing.T) {
	sections, err := llmdocs.Docs().Sections()
	if err != nil {
		t.Fatalf("Sections: %v", err)
	}
	var files []string
	for _, s := range sections {
		files = append(files, s.File)
	}
	want := "00-guide.md,10-documents.md,90-commands.md"
	if got := strings.Join(files, ","); got != want {
		t.Errorf("chapters = %s, want %s", got, want)
	}
}

// TestLegacyLLMFlag は互換の約束を守る: --llm は従来どおりコマンドラインの
// どの位置でも効く。
func TestLegacyLLMFlag(t *testing.T) {
	for _, args := range [][]string{{"--llm"}, {"invoice", "list", "--llm"}} {
		var out bytes.Buffer
		handled, err := llmcmd.HandleLegacy(args, LLMConfig(), &out)
		if err != nil {
			t.Fatalf("HandleLegacy(%v): %v", args, err)
		}
		if !handled {
			t.Errorf("HandleLegacy(%v) が --llm を処理しませんでした", args)
		}
		if !strings.Contains(out.String(), "MISOCA_CLIENT_ID") {
			t.Errorf("HandleLegacy(%v) の出力が想定と違います", args)
		}
	}
}
