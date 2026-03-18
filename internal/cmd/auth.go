package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/miyanaga/misoca-cli/internal/api"
	"github.com/miyanaga/misoca-cli/internal/oauth"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	callbackPort = 18080
	callbackURL  = "http://localhost:18080/callback"
	developURL   = "https://app.misoca.jp/oauth2/applications"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Misoca APIのOAuth2認証を行います",
	Long: `Misoca APIのOAuth2認証を行います。

環境変数 MISOCA_CLIENT_ID, MISOCA_CLIENT_SECRET が設定されていれば
即座にブラウザ認証を開始します。

未設定の場合はアプリケーション作成から案内します。
トークンは ~/.config/misoca-cli/token.json に保存され、自動的にリフレッシュされます。`,
	RunE: runAuth,
}

// resolveCredentials returns clientID and clientSecret from:
// 1. Existing token file
// 2. Environment variables
// 3. Interactive input (returns empty strings if neither found)
func resolveCredentials() (clientID, clientSecret string) {
	// Try token file first
	td, _, err := api.LoadToken()
	if err == nil && td.ClientID != "" && td.ClientSecret != "" {
		return td.ClientID, td.ClientSecret
	}

	// Try environment variables
	clientID = os.Getenv("MISOCA_CLIENT_ID")
	clientSecret = os.Getenv("MISOCA_CLIENT_SECRET")
	return
}

func runAuth(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	clientID, clientSecret := resolveCredentials()

	if clientID != "" && clientSecret != "" {
		// Fast path: credentials already available
		fmt.Println("クレデンシャルを検出しました。ブラウザ認証を開始します...")
		fmt.Println()
	} else {
		// Interactive setup
		var err error
		clientID, clientSecret, err = interactiveSetup(reader)
		if err != nil {
			return err
		}
	}

	return doOAuth2(clientID, clientSecret)
}

func interactiveSetup(reader *bufio.Reader) (clientID, clientSecret string, err error) {
	fmt.Println("=== Misoca API 初期セットアップ ===")
	fmt.Println()
	fmt.Println("  Misocaの開発者ページでアプリケーションを作成してください:")
	fmt.Println()
	fmt.Println("    アプリケーション名: misoca-cli (任意の名前)")
	fmt.Printf("    コールバックURL:    %s\n", callbackURL)
	fmt.Println()

	if clipErr := copyToClipboard(callbackURL); clipErr == nil {
		fmt.Println("  ※ コールバックURLをクリップボードにコピーしました")
		fmt.Println()
	}

	waitEnter(reader, "ブラウザで開発者ページを開きます。Enterを押してください...")
	openBrowser(developURL)
	fmt.Println()

	fmt.Println("  作成したアプリケーションの情報を貼り付けてください:")
	fmt.Println()

	fmt.Print("  アプリケーションID: ")
	line, readErr := reader.ReadString('\n')
	if readErr != nil {
		return "", "", fmt.Errorf("入力エラー: %w", readErr)
	}
	clientID = strings.TrimSpace(line)
	if clientID == "" {
		return "", "", fmt.Errorf("アプリケーションIDが入力されていません")
	}

	fmt.Print("  シークレット: ")
	line, readErr = reader.ReadString('\n')
	if readErr != nil {
		return "", "", fmt.Errorf("入力エラー: %w", readErr)
	}
	clientSecret = strings.TrimSpace(line)
	if clientSecret == "" {
		return "", "", fmt.Errorf("シークレットが入力されていません")
	}

	fmt.Println()
	return clientID, clientSecret, nil
}

func doOAuth2(clientID, clientSecret string) error {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  api.AuthURL,
			TokenURL: api.TokenURL,
		},
		RedirectURL: callbackURL,
	}

	state := fmt.Sprintf("%d", time.Now().UnixNano())
	authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	fmt.Println("  ブラウザで認証ページを開きます。「許可」をクリックしてください。")

	if err := openBrowser(authURL); err != nil {
		fmt.Println()
		fmt.Println("  ブラウザを自動で開けませんでした。以下のURLを開いてください:")
		fmt.Println()
		fmt.Printf("  %s\n", authURL)
	}

	fmt.Println()
	fmt.Println("  認証完了を待っています...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	code, err := oauth.StartCallbackServer(ctx, callbackPort)
	if err != nil {
		return fmt.Errorf("認証コードの受信に失敗しました: %w", err)
	}

	fmt.Println("  トークンを取得中...")

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("トークンの取得に失敗しました: %w", err)
	}

	td := &api.TokenData{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}
	if err := api.SaveToken(td); err != nil {
		return fmt.Errorf("トークンの保存に失敗しました: %w", err)
	}

	tokenPath, _ := api.TokenFilePath()
	fmt.Println()
	fmt.Println("認証成功！")
	fmt.Printf("  トークン保存先: %s\n", tokenPath)
	fmt.Println("  以降のコマンドは自動的に認証されます。")

	return nil
}

func waitEnter(reader *bufio.Reader, prompt string) {
	fmt.Printf("  → %s", prompt)
	reader.ReadString('\n')
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return fmt.Errorf("unsupported platform")
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
