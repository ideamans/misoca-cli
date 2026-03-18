package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	BaseURL  = "https://app.misoca.jp/api/v3"
	AuthURL  = "https://app.misoca.jp/oauth2/authorize"
	TokenURL = "https://app.misoca.jp/oauth2/token"
)

// TokenData is the persisted token file format.
type TokenData struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// Client is the Misoca API client.
type Client struct {
	httpClient *http.Client
	tokenData  *TokenData
	tokenPath  string
}

// TokenFilePath returns the path to the token file.
func TokenFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "misoca-cli", "token.json"), nil
}

// SaveToken saves token data to the config file.
func SaveToken(data *TokenData) error {
	path, err := TokenFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

// LoadToken loads token data from the config file.
func LoadToken() (*TokenData, string, error) {
	path, err := TokenFilePath()
	if err != nil {
		return nil, "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, path, err
	}
	var data TokenData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, path, err
	}
	return &data, path, nil
}

// NewClient creates a new API client using the persisted token file.
func NewClient() (*Client, error) {
	td, _, err := LoadToken()
	if err != nil {
		return nil, fmt.Errorf("トークンファイルが見つかりません。misoca-cli auth で認証してください: %w", err)
	}
	if td.RefreshToken == "" || td.ClientID == "" || td.ClientSecret == "" {
		return nil, fmt.Errorf("トークンファイルが不完全です。misoca-cli auth で再認証してください")
	}

	tokenPath, _ := TokenFilePath()

	c := &Client{
		httpClient: &http.Client{},
		tokenData:  td,
		tokenPath:  tokenPath,
	}

	// If no access token, refresh immediately
	if td.AccessToken == "" {
		if err := c.refreshAccessToken(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// refreshAccessToken uses the refresh token to obtain a new access token.
func (c *Client) refreshAccessToken() error {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.tokenData.RefreshToken},
		"client_id":     {c.tokenData.ClientID},
		"client_secret": {c.tokenData.ClientSecret},
	}

	resp, err := http.PostForm(TokenURL, form)
	if err != nil {
		return fmt.Errorf("トークンリフレッシュに失敗: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("トークンリフレッシュエラー (HTTP %d): %s\n  misoca-cli auth で再認証してください", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("トークンレスポンスの解析に失敗: %w", err)
	}

	c.tokenData.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		c.tokenData.RefreshToken = tokenResp.RefreshToken
	}
	c.tokenData.TokenType = tokenResp.TokenType

	// Persist the updated tokens
	return SaveToken(c.tokenData)
}

// Request performs an API request. Retries once on 401 by refreshing the token.
func (c *Client) Request(method, path string, query url.Values, body io.Reader) ([]byte, http.Header, error) {
	// Save body for potential retry
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, nil, fmt.Errorf("リクエストボディの読み取りに失敗: %w", err)
		}
	}

	data, header, err := c.doRequest(method, path, query, bodyBytes)
	if err != nil && isUnauthorized(err) {
		// Refresh and retry once
		if refreshErr := c.refreshAccessToken(); refreshErr != nil {
			return nil, nil, refreshErr
		}
		return c.doRequest(method, path, query, bodyBytes)
	}
	return data, header, err
}

func (c *Client) doRequest(method, path string, query url.Values, bodyBytes []byte) ([]byte, http.Header, error) {
	u := BaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var body io.Reader
	if bodyBytes != nil {
		body = strings.NewReader(string(bodyBytes))
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, nil, fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.tokenData.AccessToken)
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("APIリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("レスポンスの読み取りに失敗: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.Header, fmt.Errorf("APIエラー (HTTP %d): %s", resp.StatusCode, string(data))
	}

	return data, resp.Header, nil
}

func isUnauthorized(err error) bool {
	return err != nil && strings.Contains(err.Error(), "HTTP 401")
}

// Get performs a GET request.
func (c *Client) Get(path string, query url.Values) ([]byte, http.Header, error) {
	return c.Request(http.MethodGet, path, query, nil)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, jsonBody string) ([]byte, http.Header, error) {
	var body io.Reader
	if jsonBody != "" {
		body = strings.NewReader(jsonBody)
	}
	return c.Request(http.MethodPost, path, nil, body)
}

// Put performs a PUT request.
func (c *Client) Put(path string, jsonBody string) ([]byte, http.Header, error) {
	var body io.Reader
	if jsonBody != "" {
		body = strings.NewReader(jsonBody)
	}
	return c.Request(http.MethodPut, path, nil, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) ([]byte, http.Header, error) {
	return c.Request(http.MethodDelete, path, nil, nil)
}

// GetBinary performs a GET request for binary content (PDF, images). Retries once on 401.
// Returns the data and the Content-Type header.
func (c *Client) GetBinary(path, accept string) ([]byte, string, error) {
	data, ct, err := c.doGetBinary(path, accept)
	if err != nil && isUnauthorized(err) {
		if refreshErr := c.refreshAccessToken(); refreshErr != nil {
			return nil, "", refreshErr
		}
		return c.doGetBinary(path, accept)
	}
	return data, ct, err
}

func (c *Client) doGetBinary(path, accept string) ([]byte, string, error) {
	u := BaseURL + path

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, "", fmt.Errorf("リクエストの作成に失敗: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.tokenData.AccessToken)
	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("APIリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil, "", nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("レスポンスの読み取りに失敗: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("APIエラー (HTTP %d): %s", resp.StatusCode, string(data))
	}

	return data, resp.Header.Get("Content-Type"), nil
}

// GetPDF performs a GET request for PDF.
func (c *Client) GetPDF(path string) ([]byte, error) {
	data, _, err := c.GetBinary(path, "application/pdf")
	return data, err
}

// PrettyJSON formats raw JSON bytes with indentation.
func PrettyJSON(data []byte) (string, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return string(data), nil
	}
	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(data), nil
	}
	return string(pretty), nil
}
