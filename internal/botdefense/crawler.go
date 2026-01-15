package botdefense

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"
)

// verifies legitimate crawlers via reverse DNS
type CrawlerVerifier struct {
	allowedDomains []string
	cache          map[string]cacheEntry
	cacheMu        sync.RWMutex
	cacheTTL       time.Duration
}

type cacheEntry struct {
	verified  bool
	expiresAt time.Time
}

// creates a new crawler verifier
func NewCrawlerVerifier(allowedDomains []string) *CrawlerVerifier {
	return &CrawlerVerifier{
		allowedDomains: allowedDomains,
		cache:          make(map[string]cacheEntry),
		cacheTTL:       1 * time.Hour,
	}
}

// checks if an IP belongs to a verified crawler
// uses reverse DNS lookup and forward verification
func (v *CrawlerVerifier) IsVerifiedCrawler(ctx context.Context, ip string) bool {
	// check cache first
	v.cacheMu.RLock()

	if entry, ok := v.cache[ip]; ok && time.Now().Before(entry.expiresAt) {
		v.cacheMu.RUnlock()
		return entry.verified
	}

	v.cacheMu.RUnlock()

	// perform verification
	verified := v.verifyIP(ctx, ip)

	// cache result
	v.cacheMu.Lock()
	v.cache[ip] = cacheEntry{
		verified:  verified,
		expiresAt: time.Now().Add(v.cacheTTL),
	}

	v.cacheMu.Unlock()

	return verified
}

// performs the actual reverse DNS verification
func (v *CrawlerVerifier) verifyIP(ctx context.Context, ip string) bool {
	// reverse DNS lookup (IP -> hostname)
	resolver := net.Resolver{}

	names, err := resolver.LookupAddr(ctx, ip)
	if err != nil || len(names) == 0 {
		return false
	}

	hostname := strings.TrimSuffix(names[0], ".")

	// check if hostname matches allowed domains
	matchedDomain := ""
	for _, domain := range v.allowedDomains {
		if strings.HasSuffix(hostname, domain) {
			matchedDomain = domain
			break
		}
	}

	if matchedDomain == "" {
		return false
	}

	// forward DNS lookup (hostname -> IP) to verify
	ips, err := resolver.LookupIPAddr(ctx, hostname)
	if err != nil || len(ips) == 0 {
		return false
	}

	// check if original IP is in the forward lookup results
	for _, resolvedIP := range ips {
		if resolvedIP.IP.String() == ip {
			return true
		}
	}

	return false
}

// removes expired entries from the cache
func (v *CrawlerVerifier) CleanCache() {
	v.cacheMu.Lock()
	defer v.cacheMu.Unlock()

	now := time.Now()
	for ip, entry := range v.cache {
		if now.After(entry.expiresAt) {
			delete(v.cache, ip)
		}
	}
}

// returns the current cache size (for monitoring)
func (v *CrawlerVerifier) CacheSize() int {
	v.cacheMu.RLock()
	defer v.cacheMu.RUnlock()
	return len(v.cache)
}

// common crawler user-agent patterns for quick identification
// (before doing expensive DNS lookups)
var knownCrawlerPatterns = map[string]string{
	"googlebot":           "googlebot.com",
	"bingbot":             "search.msn.com",
	"slurp":               "yahoo.com", // Yahoo
	"duckduckbot":         "duckduckgo.com",
	"baiduspider":         "baidu.com",
	"yandexbot":           "yandex.ru",
	"facebookexternalhit": "facebook.com",
	"twitterbot":          "twitter.com",
	"linkedinbot":         "linkedin.com",
	"applebot":            "applebot.apple.com",
	"claudebot":           "anthropic.com",
	"gptbot":              "openai.com",
	"chatgpt-user":        "openai.com",
}

// does a quick check if the user-agent claims to be a known crawler
// used to decide whether to do the more expensive DNS verification
func MightBeKnownCrawler(userAgent string) (bool, string) {
	userAgentLower := strings.ToLower(userAgent)

	for pattern, domain := range knownCrawlerPatterns {
		if strings.Contains(userAgentLower, pattern) {
			return true, domain
		}
	}

	return false, ""
}
