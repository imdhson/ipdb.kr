package main

import (
	"net/http"
	"testing"
)

func TestGetIP(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		expected string
	}{
		{
			name: "CF-Connecting-IP",
			req: &http.Request{
				Header: http.Header{
					"Cf-Connecting-Ip": []string{"1.2.3.4"},
				},
			},
			expected: "1.2.3.4",
		},
		{
			name: "X-Forwarded-For Multiple",
			req: &http.Request{
				Header: http.Header{
					"X-Forwarded-For": []string{"1.2.3.4, 5.6.7.8"},
				},
			},
			expected: "1.2.3.4",
		},
		{
			name: "X-Forwarded-For Single",
			req: &http.Request{
				Header: http.Header{
					"X-Forwarded-For": []string{"1.2.3.4"},
				},
			},
			expected: "1.2.3.4",
		},
		{
			name: "RemoteAddr IPv4 with Port",
			req: &http.Request{
				RemoteAddr: "192.168.1.1:8080",
			},
			expected: "192.168.1.1",
		},
		{
			name: "RemoteAddr IPv6 with Port",
			req: &http.Request{
				RemoteAddr: "[2001:db8::1]:8080",
			},
			expected: "2001:db8::1",
		},
		{
			name: "RemoteAddr IPv4 no Port",
			req: &http.Request{
				RemoteAddr: "192.168.1.1",
			},
			expected: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Header == nil {
				tt.req.Header = http.Header{}
			}
			res := getIP(tt.req)

			if res != tt.expected {
				t.Errorf("Mismatch for %s: expected=%q, got=%q", tt.name, tt.expected, res)
			}
		})
	}
}
