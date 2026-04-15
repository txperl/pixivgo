package pixivgo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// testdataDir returns the absolute path to the testdata directory.
func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

// loadFixture reads a file from the testdata directory.
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(testdataDir(), name))
	if err != nil {
		t.Fatalf("load fixture %s: %v", name, err)
	}
	return data
}

// newTestServer creates an httptest server and a Client configured to use it.
// The client has auth tokens pre-set so API methods that require auth will work.
func newTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	client := NewClient(
		WithHTTPClient(ts.Client()),
		WithBaseURL(ts.URL),
	)
	client.SetAuth("test_access_token", "test_refresh_token")
	return client, ts
}

// newUnauthTestServer creates an httptest server and an unauthenticated Client.
func newUnauthTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	client := NewClient(
		WithHTTPClient(ts.Client()),
		WithBaseURL(ts.URL),
	)
	return client, ts
}
