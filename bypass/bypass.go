// Package bypass provides DNS-over-HTTPS resolution and SNI bypass
// for accessing Pixiv from restricted networks (e.g., behind the GFW).
//
// Usage:
//
//	httpClient, hosts, err := bypass.NewHTTPClient(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	client := pixivgo.NewClient(
//	    pixivgo.WithHTTPClient(httpClient),
//	    pixivgo.WithBaseURL(hosts),
//	)
package bypass

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const defaultHostname = "app-api.secure.pixiv.net"

// DoH server URLs to try in order.
var dohServers = []string{
	"https://1.0.0.1/dns-query",
	"https://1.1.1.1/dns-query",
	"https://doh.dns.sb/dns-query",
	"https://cloudflare-dns.com/dns-query",
}

// dnsAnswer represents a single answer record from a DoH JSON response.
type dnsAnswer struct {
	Data string `json:"data"`
}

// dnsResponse represents a DoH JSON response.
type dnsResponse struct {
	Answer []dnsAnswer `json:"Answer"`
}

// ResolveHost performs DNS-over-HTTPS resolution for the given hostname.
// It tries multiple DoH servers and returns the first successful A record.
func ResolveHost(ctx context.Context, hostname string) (string, error) {
	if hostname == "" {
		hostname = defaultHostname
	}

	client := &http.Client{Timeout: 3 * time.Second}

	var lastErr error
	for _, server := range dohServers {
		ip, err := resolveFromServer(ctx, client, server, hostname)
		if err != nil {
			lastErr = err
			continue
		}
		return ip, nil
	}

	return "", fmt.Errorf("bypass: failed to resolve %s via DoH: %w", hostname, lastErr)
}

func resolveFromServer(ctx context.Context, client *http.Client, server, hostname string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/dns-json")
	q := req.URL.Query()
	q.Set("name", hostname)
	q.Set("type", "A")
	q.Set("do", "false")
	q.Set("cd", "false")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var dnsResp dnsResponse
	if err := json.Unmarshal(body, &dnsResp); err != nil {
		return "", err
	}

	if len(dnsResp.Answer) == 0 {
		return "", fmt.Errorf("no DNS answer records for %s", hostname)
	}

	return dnsResp.Answer[0].Data, nil
}

// NewHTTPClient creates an http.Client configured for SNI bypass.
// It resolves the real IP of app-api.secure.pixiv.net via DoH,
// then creates a transport that connects to the IP directly while
// using the original hostname for TLS ServerName verification.
//
// Returns the configured client and the base URL to use (https://<IP>).
func NewHTTPClient(ctx context.Context) (*http.Client, string, error) {
	ip, err := ResolveHost(ctx, defaultHostname)
	if err != nil {
		return nil, "", err
	}

	hosts := "https://" + ip

	transport := &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Connect to the resolved IP instead of the hostname
			_, port, _ := net.SplitHostPort(addr)
			if port == "" {
				port = "443"
			}
			targetAddr := net.JoinHostPort(ip, port)

			// Establish TCP connection
			dialer := &net.Dialer{Timeout: 10 * time.Second}
			conn, err := dialer.DialContext(ctx, "tcp", targetAddr)
			if err != nil {
				return nil, err
			}

			// TLS handshake with the original hostname as ServerName
			tlsConn := tls.Client(conn, &tls.Config{
				ServerName: "app-api.pixiv.net",
			})
			if err := tlsConn.HandshakeContext(ctx); err != nil {
				conn.Close()
				return nil, err
			}

			return tlsConn, nil
		},
	}

	client := &http.Client{Transport: transport}
	return client, hosts, nil
}
