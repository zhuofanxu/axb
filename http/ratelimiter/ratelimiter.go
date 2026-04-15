package ratelimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	cleanupInterval = 5 * time.Minute  // 清理周期
	entryTTL        = 10 * time.Minute // IP 不活跃超过此时间后被清除
)

type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips map[string]*ipEntry
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter 创建一个基于IP的限流器
// r: 每秒允许的请求数
// b: 允许的突发请求数
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	l := &IPRateLimiter{
		ips: make(map[string]*ipEntry),
		r:   r,
		b:   b,
	}
	go l.cleanup()
	return l
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	e, exists := i.ips[ip]
	if !exists {
		e = &ipEntry{limiter: rate.NewLimiter(i.r, i.b)}
		i.ips[ip] = e
	}
	e.lastSeen = time.Now()
	return e.limiter
}

// cleanup 定期清理长时间不活跃的 IP，防止 map 无限增长
func (i *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		i.mu.Lock()
		for ip, e := range i.ips {
			if time.Since(e.lastSeen) > entryTTL {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}
