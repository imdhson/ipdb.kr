package main

import (
	"net/http"
	"testing"
)

func TestGetIP(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		cfIP           string
		xffIP          string
		expected       string
	}{
		// Basic direct connections (No trusted proxy)
		{"Direct IPv4", "192.0.2.1:12345", "", "", "192.0.2.1"},
		{"Direct IPv6", "[2001:db8::1]:12345", "", "", "2001:db8::1"},
		{"Direct IPv4 without port", "192.0.2.1", "", "", "192.0.2.1"},
		{"Direct IPv6 without port", "2001:db8::1", "", "", "2001:db8::1"},
		{"Invalid IP format", "not-an-ip", "", "", "Unknown"},

		// Cloudflare Header Tests (Requires trusted proxy)
		{"CF Header with Localhost proxy", "127.0.0.1:12345", "192.0.2.2", "", "192.0.2.2"},
		{"CF Header with IPv6 Localhost proxy", "[::1]:12345", "192.0.2.2", "", "192.0.2.2"},
		{"CF Header with Private IP proxy", "10.0.0.1:12345", "192.0.2.2", "", "192.0.2.2"},
		{"CF Header with External Proxy (Spoofing attempt)", "198.51.100.1:12345", "192.0.2.2", "", "198.51.100.1"},

		// X-Forwarded-For Header Tests (Requires trusted proxy)
		{"XFF Header single IP", "127.0.0.1:12345", "", "192.0.2.3", "192.0.2.3"},
		{"XFF Header multiple IPs", "127.0.0.1:12345", "", "198.51.100.1, 192.0.2.3", "192.0.2.3"}, // Should pick right-most
		{"XFF Header multiple IPs (lots of spaces)", "127.0.0.1:12345", "", "198.51.100.1,   192.0.2.3  ", "192.0.2.3"},
		{"XFF Header with External Proxy (Spoofing attempt)", "198.51.100.1:12345", "", "192.0.2.3", "198.51.100.1"},

		// Priority Tests (Requires trusted proxy)
		{"CF overrides XFF", "127.0.0.1:12345", "192.0.2.2", "192.0.2.3", "192.0.2.2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.RemoteAddr = tt.remoteAddr
			if tt.cfIP != "" {
				req.Header.Set("CF-Connecting-IP", tt.cfIP)
			}
			if tt.xffIP != "" {
				req.Header.Set("X-Forwarded-For", tt.xffIP)
			}

			ip := getIP(req)
			if ip != tt.expected {
				t.Errorf("getIP() = %v, want %v", ip, tt.expected)
			}
		})
	}
}

func TestValidateCSRF(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		origin  string
		referer string
		want    bool
	}{
		{"Valid Origin HTTP", "example.com:8080", "http://example.com:8080", "", true},
		{"Valid Origin HTTPS", "example.com:8080", "https://example.com:8080", "", true},
		{"Valid Referer HTTP", "example.com:8080", "", "http://example.com:8080/path", true},
		{"Valid Referer HTTPS", "example.com:8080", "", "https://example.com:8080/path?q=1", true},

		{"Invalid Origin (wrong host)", "example.com:8080", "http://attacker.com:8080", "", false},
		{"Invalid Referer (wrong host)", "example.com:8080", "", "http://attacker.com/path", false},

		{"Empty Headers", "example.com:8080", "", "", false},

		// Edge cases where the URL parser might fail but we should handle it
		{"Referer with query but no path", "example.com", "", "http://example.com?query=1", true},
		{"Origin with fragment but no path", "example.com", "http://example.com#frag", "", true},

		{"Invalid Scheme", "example.com", "ftp://example.com", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Host = tt.host
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.referer != "" {
				req.Header.Set("Referer", tt.referer)
			}

			if got := validateCSRF(req); got != tt.want {
				t.Errorf("validateCSRF() = %v, want %v", got, tt.want)
			}
		})
	}
}
