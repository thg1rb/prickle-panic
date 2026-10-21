package clientip

import (
	"net/netip"
	"testing"
)

func TestResolverTrustBoundary(t *testing.T) {
	r, err := NewResolver([]netip.Prefix{netip.MustParsePrefix("10.0.0.0/8")}, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, remote, xff, want string
		bad                     bool
	}{
		{"direct ignores spoofed headers", "192.0.2.9:1234", "198.51.100.1", "192.0.2.9", false},
		{"untrusted peer cannot spoof source", "203.0.113.8:1234", "198.51.100.1, 10.0.0.2", "203.0.113.8", false},
		{"trusted proxy selects first untrusted hop", "10.0.0.2:443", "198.51.100.8, 10.1.0.2", "198.51.100.8", false},
		{"trusted peer needs chain", "10.0.0.2:443", "", "", true},
		{"malformed chain", "10.0.0.2:443", "198.51.100.1, garbage", "", true},
		{"zone rejected", "[fe80::1%lo]:443", "", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := r.Resolve(tc.remote, tc.xff)
			if tc.bad {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil || got.String() != tc.want {
				t.Fatalf("got=%v err=%v want=%s", got, err, tc.want)
			}
		})
	}
}

func TestResolverRejectsInvalidProxyConfiguration(t *testing.T) {
	if _, err := NewResolver([]netip.Prefix{netip.Prefix{}}, 0); err == nil {
		t.Fatal("invalid prefix accepted")
	}
	if _, err := NewResolver([]netip.Prefix{netip.MustParsePrefix("10.0.0.1/8")}, 0); err == nil {
		t.Fatal("non-canonical prefix accepted")
	}
	if _, err := NewResolver(nil, 257); err == nil {
		t.Fatal("unbounded hop count accepted")
	}
}
