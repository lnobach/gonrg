package util

import (
	"net/http"
	"net/url"
	"strings"
)

func CorsWsIsOriginAllowed(r *http.Request, allowed []string) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		// no Origin set, so the client doesn't care
		return true
	}

	for _, a := range allowed {
		if DomainMatches(origin[0], a) {
			return true
		}
	}

	// Lastly, check if the origin corresponds to the hostname
	origurl, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	return EqualFoldNonUnicode(origurl.Host, r.Host)
}

func DomainMatches(domain, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == "" {
		return false
	}
	if len(domain) > 0 && strings.HasPrefix(pattern, "*") &&
		len(pattern)-1 <= len(domain) &&
		EqualFoldNonUnicode(pattern[1:], domain[len(domain)-len(pattern)+1:]) {
		return true
	}
	return EqualFoldNonUnicode(domain, pattern)
}
