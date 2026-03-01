package quickbooks

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// a test client in case when we want to test the client method
// export_test.go
func NewTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	endpoint, err := url.Parse(server.URL + "/v3/company/test-realm/")
	require.NoError(t, err)

	return &Client{
		Client:       server.Client(),
		endpoint:     endpoint,
		realm:        "test-realm",
		minorVersion: "65",
	}, server
}
