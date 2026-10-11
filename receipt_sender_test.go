package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestCreateTemplateRequest(t *testing.T) {
	client := &Client{key: "test", http: &http.Client{Transport: transportFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/v1/sms/template/create" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		var fields map[string]json.RawMessage
		if err := json.NewDecoder(req.Body).Decode(&fields); err != nil {
			t.Fatal(err)
		}
		if len(fields) != 2 || string(fields["name"]) != `"approved-payment"` || string(fields["body"]) != `"Payment received"` {
			t.Fatalf("unexpected template fields: %v", fields)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"ok":true}`))}, nil
	})}}
	id, err := client.CreateTemplate(context.Background(), "approved-payment", "Payment received")
	if err != nil || id != "approved-payment" {
		t.Fatalf("id = %q, err = %v", id, err)
	}
}
