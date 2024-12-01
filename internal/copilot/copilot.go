// internal/copilot/copilot.go
package copilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	deviceCodeURL   = "https://github.com/login/device/code"
	tokenURL        = "https://github.com/login/oauth/access_token"
	copilotTokenAPI = "https://api.github.com/copilot_internal/v2/token"
	clientID        = "Iv23liXYYLdDjdQbE0BB" // GitHub Copilot client ID
)

type Client struct {
	configDir string
	client    *http.Client
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type CopilotToken struct {
	Token string `json:"token"`
}

func New(configDir string) *Client {
	log.Printf("Initializing Copilot client with config directory: %s", configDir)
	return &Client{
		configDir: configDir,
		client:    &http.Client{},
	}
}

func (c *Client) RequestDeviceCode() (*DeviceCodeResponse, error) {
	log.Println("Requesting device code from GitHub...")

	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"client_id":"%s","scope":"copilot"}`, clientID)))

	req, err := http.NewRequest("POST", deviceCodeURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device code request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Device code response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("HTTP error: %s", resp.Status)
		}
		return nil, fmt.Errorf("device code error: %s - %s", errResp.Error, errResp.ErrorDescription)
	}

	var deviceCode DeviceCodeResponse
	if err := json.Unmarshal(body, &deviceCode); err != nil {
		return nil, fmt.Errorf("failed to parse device code response: %w", err)
	}

	log.Printf("Successfully received device code with verification URI: %s", deviceCode.VerificationURI)
	return &deviceCode, nil
}

func (c *Client) VerifyAuth(deviceCode string) (*AuthResponse, error) {
	log.Println("Verifying authentication status...")

	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(
		`{"client_id":"%s","device_code":"%s","grant_type":"urn:ietf:params:oauth:grant-type:device_code"}`,
		clientID,
		deviceCode,
	)))

	req, err := http.NewRequest("POST", tokenURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth verification request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Auth verification response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("HTTP error: %s", resp.Status)
		}
		return nil, fmt.Errorf("auth error: %s - %s", errResp.Error, errResp.ErrorDescription)
	}

	var auth AuthResponse
	if err := json.Unmarshal(body, &auth); err != nil {
		return nil, fmt.Errorf("failed to parse auth response: %w", err)
	}

	if auth.AccessToken == "" {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error == "authorization_pending" {
			log.Println("Authorization still pending...")
			return nil, nil
		}
		return nil, fmt.Errorf("no access token in response")
	}

	log.Printf("Successfully received access token of length: %d", len(auth.AccessToken))
	return &auth, nil
}

func (c *Client) SaveAuthToken(token string) error {
	log.Println("Saving auth token...")

	if err := os.MkdirAll(c.configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	tokenPath := filepath.Join(c.configDir, "copilot_token")
	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	log.Printf("Successfully saved token to: %s", tokenPath)
	return nil
}

func (c *Client) RemoveAuthToken() error {
	log.Println("Removing auth token...")

	tokenPath := filepath.Join(c.configDir, "copilot_token")
	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token file: %w", err)
	}

	log.Println("Successfully removed token file")
	return nil
}

func (c *Client) LoadAuthToken() (string, error) {
	log.Println("Loading auth token...")

	tokenPath := filepath.Join(c.configDir, "copilot_token")
	token, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No token file found")
			return "", nil
		}
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	log.Printf("Successfully loaded token of length: %d", len(token))
	return string(token), nil
}

func (c *Client) GetAPIToken() (string, error) {
	authToken, err := c.LoadAuthToken()
	if err != nil {
		return "", fmt.Errorf("failed to load auth token: %w", err)
	}
	if authToken == "" {
		return "", fmt.Errorf("no auth token found, please run 'copilot-login' first")
	}

	req, err := http.NewRequest("GET", copilotTokenAPI, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+authToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Editor-Version", "vscode/1.88.0")
	req.Header.Set("Editor-Plugin-Version", "copilot-chat/0.14.2024032901")
	req.Header.Set("User-Agent", "GitHubCopilotChat/0.14.2024032901")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token API error: %s", resp.Status)
	}

	var tokenResp CopilotToken
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResp.Token, nil
}

// This method is kept for backward compatibility
func (c *Client) GetCopilotToken(model string) (*CopilotToken, error) {
	log.Printf("Getting Copilot API token for model: %s", model)
	token, err := c.GetAPIToken()
	if err != nil {
		return nil, err
	}
	return &CopilotToken{Token: token}, nil
}
