// Package clientip resolves a canonical source address at the trusted proxy boundary.
package clientip

import (
	"errors"
	"net/netip"
	"strings"
)

var ErrInvalidSource = errors.New("invalid client source metadata")

// Resolver trusts X-Forwarded-For only when the direct peer matches a
// configured trusted proxy prefix. Forwarded and X-Real-IP are never used.
type Resolver struct {
	trusted []netip.Prefix
	maxHops int
}

// NewResolver validates and copies trusted proxy CIDRs. A zero maxHops selects
// the default limit; larger limits are capped to bound header parsing work.
func NewResolver(trusted []netip.Prefix, maxHops int) (Resolver, error) {
	if maxHops < 0 || maxHops > 256 {
		return Resolver{}, ErrInvalidSource
	}
	copyOfTrusted := make([]netip.Prefix, len(trusted))
	for i, prefix := range trusted {
		if !prefix.IsValid() || prefix != prefix.Masked() {
			return Resolver{}, ErrInvalidSource
		}
		copyOfTrusted[i] = prefix
	}
	if maxHops == 0 {
		maxHops = 32
	}
	return Resolver{trusted: copyOfTrusted, maxHops: maxHops}, nil
}

func (r Resolver) Resolve(remoteAddr, forwardedFor string) (netip.Addr, error) {
	peer, err := parseAddr(remoteAddr)
	if err != nil {
		return netip.Addr{}, ErrInvalidSource
	}
	peer = peer.Unmap()
	if !r.trustedAddr(peer) {
		return peer, nil
	}
	if strings.TrimSpace(forwardedFor) == "" {
		return netip.Addr{}, ErrInvalidSource
	}
	parts := strings.Split(forwardedFor, ",")
	max := r.maxHops
	if len(parts) > max {
		return netip.Addr{}, ErrInvalidSource
	}
	chain := make([]netip.Addr, 0, len(parts)+1)
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return netip.Addr{}, ErrInvalidSource
		}
		a, e := netip.ParseAddr(p)
		if e != nil || a.Zone() != "" {
			return netip.Addr{}, ErrInvalidSource
		}
		chain = append(chain, a.Unmap())
	}
	chain = append(chain, peer)
	idx := len(chain) - 2
	for idx > 0 && r.trustedAddr(chain[idx]) {
		idx--
	}
	return chain[idx], nil
}

func parseAddr(s string) (netip.Addr, error) {
	if a, err := netip.ParseAddrPort(s); err == nil {
		if a.Addr().Zone() != "" {
			return netip.Addr{}, ErrInvalidSource
		}
		return a.Addr(), nil
	}
	a, err := netip.ParseAddr(s)
	if err != nil || a.Zone() != "" {
		return netip.Addr{}, ErrInvalidSource
	}
	return a, nil
}

func (r Resolver) trustedAddr(addr netip.Addr) bool {
	for _, p := range r.trusted {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}
