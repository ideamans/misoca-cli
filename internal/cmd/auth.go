package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ideamans/misoca-cli/internal/api"
	"github.com/ideamans/misoca-cli/internal/oauth"
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
トークンは ~/.config/misoca-cli/token.json に保存され、自動的にリフレッシュされます。

開くURLは常に表示されます。ブラウザの無いサーバーでは、表示されたURLを
手元のブラウザで開いて許可し、リダイレクト先（http://localhost:18080/callback?code=...、
接続エラーのページになって構いません）のURLをアドレスバーからコピーして
ターミナルに貼り付けてください。認証コードだけを貼り付けても構いません。
DISPLAY の無い Linux ではブラウザを自動では開きません（--no-browser と同じ）。`,
	RunE: runAuth,
}

var authNoBrowser bool

func init() {
	authCmd.Flags().BoolVar(&authNoBrowser, "no-browser", false, "ブラウザを開かず、URLの表示だけを行う（SSH先のサーバーなど）")
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

	return doOAuth2(reader, clientID, clientSecret)
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

	fmt.Printf("  開発者ページ: %s\n", developURL)
	fmt.Println()
	if canOpenBrowser() {
		waitEnter(reader, "ブラウザで開発者ページを開きます。Enterを押してください...")
		openBrowser(developURL)
		fmt.Println()
	}

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

func doOAuth2(reader *bufio.Reader, clientID, clientSecret string) error {
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

	fmt.Println("  以下のURLをブラウザで開き、「許可」をクリックしてください:")
	fmt.Println()
	fmt.Printf("  %s\n", authURL)
	fmt.Println()

	if canOpenBrowser() {
		if err := openBrowser(authURL); err == nil {
			fmt.Println("  （ブラウザを自動で開きました）")
			fmt.Println()
		}
	}

	fmt.Println("  ブラウザの無いサーバーでは、許可後に開かれるURL")
	fmt.Printf("  （%s?code=... 接続エラーで構いません）を\n", callbackURL)
	fmt.Println("  アドレスバーからコピーして、ここに貼り付けてEnterを押してください。")
	fmt.Println()
	fmt.Println("  認証完了を待っています...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	code, err := waitAuthCode(ctx, reader, state)
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

// waitAuthCode はローカルのコールバックと、標準入力への貼り付けの早い方から
// 認証コードを受け取ります。後者はブラウザの無いサーバーで使います。
func waitAuthCode(ctx context.Context, reader *bufio.Reader, state string) (string, error) {
	type result struct {
		code string
		err  error
	}
	ch := make(chan result, 2)

	go func() {
		code, err := oauth.StartCallbackServer(ctx, callbackPort)
		ch <- result{code, err}
	}()

	go func() {
		for {
			line, err := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line != "" {
				code, perr := parsePastedCode(line, state)
				if perr != nil {
					fmt.Printf("  %v。もう一度貼り付けてください: ", perr)
					continue
				}
				ch <- result{code, nil}
				return
			}
			if err != nil {
				// 標準入力が閉じている（パイプ等）ときはコールバックだけを待つ
				ch <- result{"", fmt.Errorf("標準入力が閉じられました")}
				return
			}
		}
	}()

	// 片方が失敗しても（ポートが使用中、標準入力が無いなど）もう片方を待つ
	var firstErr error
	for range 2 {
		select {
		case r := <-ch:
			if r.err == nil {
				return r.code, nil
			}
			if firstErr == nil {
				firstErr = r.err
				fmt.Printf("  （%v。引き続き待っています）\n", r.err)
			}
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	return "", firstErr
}

// parsePastedCode は貼り付けられたリダイレクト先URL、または認証コードそのものから
// 認証コードを取り出します。
func parsePastedCode(input, state string) (string, error) {
	if !strings.Contains(input, "://") && !strings.Contains(input, "?") {
		return input, nil
	}
	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("URLを解釈できません")
	}
	q := u.Query()
	if e := q.Get("error"); e != "" {
		return "", fmt.Errorf("認証が拒否されました (%s)", e)
	}
	if s := q.Get("state"); s != "" && s != state {
		return "", fmt.Errorf("state が一致しません（別の認証のURLです）")
	}
	code := q.Get("code")
	if code == "" {
		return "", fmt.Errorf("URLに code が含まれていません")
	}
	return code, nil
}

// canOpenBrowser はブラウザを自動で開いてよいかを返します。
// DISPLAY の無い Linux で xdg-open を呼ぶと、端末内のテキストブラウザが
// 起動して入力を奪うことがあるため開きません。
func canOpenBrowser() bool {
	if authNoBrowser {
		return false
	}
	if runtime.GOOS == "linux" && os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return false
	}
	return true
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
