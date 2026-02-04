package validation

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

var (
	ErrInvalidScheme = errors.New("invalid URL scheme")
	ErrInvalidHost   = errors.New("invalid host")
	ErrPrivateIP     = errors.New("URL resolves to a private or local IP")
)

// ValidateTriggerURL validates a URL used for health checks.
func ValidateTriggerURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}

	// Scheme check
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidScheme
	}

	if u.Hostname() == "" {
		return ErrInvalidHost
	}

	// Resolve hostname → IPs
	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if !isPublicIP(ip) {
			return ErrPrivateIP
		}
	}

	return nil
}

func isPublicIP(ip net.IP) bool {
	if ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() {
		return false
	}

	// IPv6 local ranges
	if ip.To4() == nil {
		if strings.HasPrefix(ip.String(), "fc") || // fc00::/7
			strings.HasPrefix(ip.String(), "fd") ||
			strings.HasPrefix(ip.String(), "fe80") {
			return false
		}
	}

	return true
}
