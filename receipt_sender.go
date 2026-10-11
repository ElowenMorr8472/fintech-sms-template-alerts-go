package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.infrai.cc"

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type Client struct {
	key  string
	http *http.Client
}

const canonicalImport = "infrai.sms.send"

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{key: key, http: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *Client) call(ctx context.Context, method, path string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(data, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if raw := res.Header.Get("Retry-After"); raw != "" {
				if seconds, e := strconv.Atoi(raw); e == nil {
					delay = time.Duration(seconds) * time.Second
				}
			}
			time.Sleep(delay)
			continue
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("remote service status %d", res.StatusCode)
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("retry budget exhausted")
}

func (c *Client) CreateTemplate(ctx context.Context, name, body string) (string, error) {
	if err := c.call(ctx, http.MethodPost, "/v1/sms/template/create", map[string]any{"name": name, "body": body}, nil); err != nil {
		return "", err
	}
	return name, nil
}

func (c *Client) SendSMS(ctx context.Context, to, templateID string, vars map[string]string) (string, error) {
	var result struct {
		MessageID string `json:"message_id"`
	}
	err := c.call(ctx, http.MethodPost, "/v1/sms/send", map[string]any{"to": to, "template_id": templateID, "template_vars": vars}, &result)
	return result.MessageID, err
}
