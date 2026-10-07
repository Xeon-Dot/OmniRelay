package main

import "testing"

func TestIsSameOrigin(t *testing.T) {
	cases := []struct {
		origin string
		host   string
		want   bool
	}{
		// Same-origin: Origin scheme+host matches the request Host.
		{"https://relay.example.com", "relay.example.com", true},
		{"http://relay.example.com", "relay.example.com", true},
		{"https://relay.example.com:8443", "relay.example.com:8443", true},
		// Cross-origin must be rejected when no allow list is configured.
		{"https://evil.example.com", "relay.example.com", false},
		{"https://relay.example.com", "evil.example.com", false},
		{"null", "relay.example.com", false},
		{"", "relay.example.com", false},
		// Subdomain is a different origin.
		{"https://sub.relay.example.com", "relay.example.com", false},
	}
	for _, tc := range cases {
		if got := isSameOrigin(tc.origin, tc.host); got != tc.want {
			t.Errorf("isSameOrigin(%q, %q) = %v, want %v", tc.origin, tc.host, got, tc.want)
		}
	}
}
