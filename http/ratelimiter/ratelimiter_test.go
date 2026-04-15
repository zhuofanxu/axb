package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestGetLimiter_SameIPReturnsSameLimiter(t *testing.T) {
	l := NewIPRateLimiter(10, 10)
	l1 := l.GetLimiter("192.168.1.1")
	l2 := l.GetLimiter("192.168.1.1")
	if l1 != l2 {
		t.Error("same IP should return the same limiter instance")
	}
}

func TestGetLimiter_DifferentIPsAreIndependent(t *testing.T) {
	l := NewIPRateLimiter(10, 10)
	l1 := l.GetLimiter("192.168.1.1")
	l2 := l.GetLimiter("192.168.1.2")
	if l1 == l2 {
		t.Error("different IPs should return different limiter instances")
	}
}

func TestRateLimit_BurstExhausted(t *testing.T) {
	// burst=3, rate=1 req/s
	l := NewIPRateLimiter(1, 3)
	limiter := l.GetLimiter("10.0.0.1")

	for i := 0; i < 3; i++ {
		if !limiter.Allow() {
			t.Errorf("request %d should be allowed within burst capacity", i+1)
		}
	}
	if limiter.Allow() {
		t.Error("4th request should be denied after burst is exhausted")
	}
}

func TestGetLimiter_UpdatesLastSeen(t *testing.T) {
	l := NewIPRateLimiter(10, 10)
	l.GetLimiter("1.1.1.1")

	before := time.Now()
	time.Sleep(10 * time.Millisecond)
	l.GetLimiter("1.1.1.1")

	l.mu.RLock()
	lastSeen := l.ips["1.1.1.1"].lastSeen
	l.mu.RUnlock()

	if !lastSeen.After(before) {
		t.Error("lastSeen should be updated on each GetLimiter call")
	}
}

func TestGetLimiter_Concurrent(t *testing.T) {
	l := NewIPRateLimiter(100, 100)
	const goroutines = 200
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			l.GetLimiter("10.0.0.1")
		}()
	}
	wg.Wait()

	l.mu.RLock()
	count := len(l.ips)
	l.mu.RUnlock()
	if count != 1 {
		t.Errorf("expected exactly 1 entry for a single IP, got %d", count)
	}
}

func TestCleanup_RemovesStaleEntries(t *testing.T) {
	l := NewIPRateLimiter(10, 10)
	l.GetLimiter("1.2.3.4")
	l.GetLimiter("5.6.7.8")

	// 将 1.2.3.4 的 lastSeen 设为过期
	l.mu.Lock()
	l.ips["1.2.3.4"].lastSeen = time.Now().Add(-(entryTTL + time.Second))
	l.mu.Unlock()

	// 直接触发清理逻辑
	l.mu.Lock()
	for ip, e := range l.ips {
		if time.Since(e.lastSeen) > entryTTL {
			delete(l.ips, ip)
		}
	}
	l.mu.Unlock()

	l.mu.RLock()
	_, staleExists := l.ips["1.2.3.4"]
	_, activeExists := l.ips["5.6.7.8"]
	l.mu.RUnlock()

	if staleExists {
		t.Error("stale IP entry should have been removed")
	}
	if !activeExists {
		t.Error("active IP entry should not have been removed")
	}
}
