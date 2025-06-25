package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseURL string
	id      string
	secret  string

	httpClient *http.Client
	mu         sync.RWMutex
	token      string
	expire     time.Time
}

func NewClient(baseURL, id, secret string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		id:      id,
		secret:  secret,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Token(ctx context.Context) (string, error) {
	c.mu.RLock()
	if time.Until(c.expire) > 30*time.Second {
		tok := c.token
		c.mu.RUnlock()
		return tok, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Until(c.expire) > 30*time.Second {
		return c.token, nil
	}

	newTok, newExp, err := c.fetch(ctx)
	if err != nil {
		return "", err
	}

	c.token, c.expire = newTok, newExp
	return newTok, nil
}

func (c *Client) fetch(ctx context.Context) (string, time.Time, error) {
	body, _ := json.Marshal(map[string]string{
		"clientId":     c.id,
		"clientSecret": c.secret,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/v1/auth/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("status %s", resp.Status)
	}
	var tr struct {
		AccessToken     string `json:"accessToken"`
		AccessExpiresAt int64  `json:"accessExpiresAt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", time.Time{}, err
	}
	return tr.AccessToken, time.Unix(tr.AccessExpiresAt, 0), nil
}
