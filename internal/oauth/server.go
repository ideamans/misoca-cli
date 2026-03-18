package oauth

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// StartCallbackServer starts a local HTTP server to receive the OAuth2 callback.
// It returns the authorization code received via the callback.
func StartCallbackServer(ctx context.Context, port int) (string, error) {
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error")
			if errMsg == "" {
				errMsg = "authorization code not found"
			}
			fmt.Fprintf(w, "<html><body><h2>認証エラー</h2><p>%s</p><p>このウィンドウを閉じてください。</p></body></html>", errMsg)
			errCh <- fmt.Errorf("OAuth2 error: %s", errMsg)
			return
		}

		fmt.Fprint(w, `<html><body><h2>認証成功！</h2><p>このウィンドウを閉じてターミナルに戻ってください。</p></body></html>`)
		codeCh <- code
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return "", fmt.Errorf("failed to start callback server on port %d: %w", port, err)
	}

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	defer server.Shutdown(context.Background())

	select {
	case code := <-codeCh:
		return code, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
