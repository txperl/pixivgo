package pixivgo

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// --- HTTP method routing tests (migrated from test_request_call_*) ---

func TestDoRequest_Get(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	resp, err := client.doGet(context.Background(), "/v1/test", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestDoRequest_Post(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q, want application/x-www-form-urlencoded", ct)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm error: %v", err)
		}
		if got := r.PostForm.Get("key"); got != "value" {
			t.Errorf("form key = %q, want \"value\"", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	data := map[string][]string{"key": {"value"}}
	resp, err := client.doPost(context.Background(), "/v1/test", data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestDoRequest_Delete(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	resp, err := client.doRequest(context.Background(), http.MethodDelete, "/v1/test", nil, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

// --- Header tests (migrated from test_set_accept_language / test_set_additional_headers) ---

func TestSetAcceptLanguage(t *testing.T) {
	client := NewClient()
	client.SetAcceptLanguage("en-us")

	if got := client.additionalHeaders.Get("Accept-Language"); got != "en-us" {
		t.Errorf("Accept-Language = %q, want \"en-us\"", got)
	}
}

func TestSetAdditionalHeaders(t *testing.T) {
	client := NewClient()
	headers := http.Header{"Keep-Alive": {"timeout=5, max=1000"}}
	client.SetAdditionalHeaders(headers)

	if got := client.additionalHeaders.Get("Keep-Alive"); got != "timeout=5, max=1000" {
		t.Errorf("Keep-Alive = %q, want \"timeout=5, max=1000\"", got)
	}
}

func TestDoRequest_SetsHeaders(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != defaultUserAgent {
			t.Errorf("User-Agent = %q, want %q", got, defaultUserAgent)
		}
		if got := r.Header.Get("App-OS"); got != defaultAppOS {
			t.Errorf("App-OS = %q, want %q", got, defaultAppOS)
		}
		if got := r.Header.Get("App-OS-Version"); got != defaultAppOSVer {
			t.Errorf("App-OS-Version = %q, want %q", got, defaultAppOSVer)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	resp, err := client.doGet(context.Background(), "/v1/test", nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestDoRequest_AuthHeader(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test_access_token" {
			t.Errorf("Authorization = %q, want \"Bearer test_access_token\"", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	resp, err := client.doGet(context.Background(), "/v1/test", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestDoRequest_NoAuth_SkipsBearer(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("Authorization should be empty when noAuth=true, got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	resp, err := client.doGet(context.Background(), "/v1/test", nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestDoRequest_AdditionalHeaders(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept-Language"); got != "zh-cn" {
			t.Errorf("Accept-Language = %q, want \"zh-cn\"", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	client.SetAcceptLanguage("zh-cn")

	resp, err := client.doGet(context.Background(), "/v1/test", nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

// --- Auth requirement test (migrated from test_require_auth) ---

func TestRequireAuth(t *testing.T) {
	client := NewClient()

	if client.getAccessToken() != "" {
		t.Error("access token should be empty before auth")
	}

	err := client.requireAuth()
	if !errors.Is(err, ErrAuthRequired) {
		t.Errorf("requireAuth() = %v, want ErrAuthRequired", err)
	}
}

func TestRequireAuth_WithToken(t *testing.T) {
	client := NewClient()
	client.SetAuth("some_token", "")

	if err := client.requireAuth(); err != nil {
		t.Errorf("requireAuth() after SetAuth should return nil, got %v", err)
	}
}

// --- parseResponse tests (migrated from test_parse_json_valid / test_parse_json_invalid) ---

func TestParseResponse_ValidJSON(t *testing.T) {
	type testResp struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(`{"name":"test","age":25}`)),
	}

	result, err := parseResponse[testResp](resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "test" || result.Age != 25 {
		t.Errorf("got {Name:%q, Age:%d}, want {Name:\"test\", Age:25}", result.Name, result.Age)
	}
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	invalidJSON := loadFixture(t, "general_invalid_json.txt")

	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(string(invalidJSON))),
	}

	type dummy struct{}
	_, err := parseResponse[dummy](resp)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}

	var pe *PixivError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PixivError, got %T", err)
	}
	if !strings.Contains(pe.Message, "json unmarshal error") {
		t.Errorf("error message = %q, want to contain 'json unmarshal error'", pe.Message)
	}
}

func TestParseResponse_ValidFixtureJSON(t *testing.T) {
	// Use the general_valid_json.json fixture to verify valid JSON parsing
	validJSON := loadFixture(t, "general_valid_json.json")

	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(string(validJSON))),
	}

	var raw json.RawMessage
	result, err := parseResponse[json.RawMessage](resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw = *result
	if len(raw) == 0 {
		t.Error("parsed result should not be empty")
	}
}

func TestParseResponse_HTTPError(t *testing.T) {
	resp := &http.Response{
		StatusCode: 403,
		Header:     http.Header{"X-Test": {"value"}},
		Body:       io.NopCloser(strings.NewReader(`{"error":"forbidden"}`)),
	}

	type dummy struct{}
	_, err := parseResponse[dummy](resp)
	if err == nil {
		t.Fatal("expected error for HTTP 403, got nil")
	}

	var pe *PixivError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PixivError, got %T", err)
	}
	if pe.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403", pe.StatusCode)
	}
	if pe.Body != `{"error":"forbidden"}` {
		t.Errorf("Body = %q, want %q", pe.Body, `{"error":"forbidden"}`)
	}
	if pe.Header.Get("X-Test") != "value" {
		t.Errorf("Header X-Test = %q, want \"value\"", pe.Header.Get("X-Test"))
	}
}

func TestDoRequest_AuthRequired_ReturnsError(t *testing.T) {
	client, _ := newUnauthTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not be made when auth is required but not set")
	})

	_, err := client.doGet(context.Background(), "/v1/test", nil, false)
	if !errors.Is(err, ErrAuthRequired) {
		t.Errorf("got %v, want ErrAuthRequired", err)
	}
}
