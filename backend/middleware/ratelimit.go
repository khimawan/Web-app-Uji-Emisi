package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// bucket menyimpan token tersisa untuk satu klien (berdasarkan IP).
type bucket struct {
	tokens   float64
	lastSeen time.Time
}

// limiter adalah token bucket in-memory dengan penghitung single-goroutine
// (mutex). Cocok untuk satu proses; kalau nanti dijalankan lebih dari satu
// instance, state ini harus dipindah ke Redis.
type limiter struct {
	mu       sync.Mutex
	entries  map[string]*bucket
	perSec   float64
	capacity float64

	lastGC  time.Time
	gcEvery time.Duration
	idleTTL time.Duration
}

func newLimiter(requestsPerMinute, burst int) *limiter {
	if requestsPerMinute < 1 {
		requestsPerMinute = 1
	}
	if burst < 1 {
		burst = 1
	}
	return &limiter{
		entries:  make(map[string]*bucket),
		perSec:   float64(requestsPerMinute) / 60.0,
		capacity: float64(burst),
		lastGC:   time.Now(),
		gcEvery:  time.Minute,
		idleTTL:  10 * time.Minute,
	}
}

// allow mengembalikan true bila request diizinkan, beserta lama tunggu sampai
// satu token berikutnya tersedia bila ditolak.
func (l *limiter) allow(key string) (bool, time.Duration) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.collectLocked(now)

	entry, ok := l.entries[key]
	if !ok {
		// Klien baru: boleh sekaligus sampai burst, lalu ambil satu token.
		l.entries[key] = &bucket{tokens: l.capacity - 1, lastSeen: now}
		return true, 0
	}

	entry.tokens += now.Sub(entry.lastSeen).Seconds() * l.perSec
	if entry.tokens > l.capacity {
		entry.tokens = l.capacity
	}
	entry.lastSeen = now

	if entry.tokens < 1 {
		wait := time.Duration((1 - entry.tokens) / l.perSec * float64(time.Second))
		if wait < time.Second {
			wait = time.Second
		}
		if wait > 60*time.Second {
			wait = 60 * time.Second
		}
		return false, wait
	}

	entry.tokens--
	return true, 0
}

// collectLocked membuang bucket yang sudah lama tidak dipakai.
//
// WAJIB ADA: tanpa ini setiap IP baru akan menambah entri selamanya dan
// server yang uptime-nya panjang akan kehabisan RAM.
func (l *limiter) collectLocked(now time.Time) {
	if now.Sub(l.lastGC) < l.gcEvery {
		return
	}
	l.lastGC = now
	for key, entry := range l.entries {
		if now.Sub(entry.lastSeen) > l.idleTTL {
			delete(l.entries, key)
		}
	}
}

// RateLimit membatasi jumlah request per menit per klien dengan token bucket.
//
// Key diambil dari c.ClientIP(), yang hanya benar bila router().SetTrustedProxies
// sudah dikonfigurasi. Kalau tidak, klien bisa memalsukan X-Forwarded-For untuk
// menghindari limit.
func RateLimit(requestsPerMinute, burst int) gin.HandlerFunc {
	l := newLimiter(requestsPerMinute, burst)

	return func(c *gin.Context) {
		allowed, retryAfter := l.allow(c.ClientIP())
		if allowed {
			c.Next()
			return
		}

		seconds := int(retryAfter.Seconds())
		if seconds < 1 {
			seconds = 1
		}
		c.Header("Retry-After", strconv.Itoa(seconds))
		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"success":            false,
			"message":            "Terlalu banyak permintaan. Coba lagi sebentar lagi.",
			"retry_after_seconds": seconds,
		})
	}
}

// MaxConcurrency membatasi berapa request yang boleh diproses bersamaan.
//
// Rate limit menahan frekuensi, tapi tidak menahan request yang lambat dan
// berat. Aturan ini yang menjaga agar backend tidak tumbang ketika banyak
// pengguna membuka halaman laporan bersamaan.
func MaxConcurrency(maxInFlight int) gin.HandlerFunc {
	if maxInFlight < 1 {
		maxInFlight = 1
	}
	semaphore := make(chan struct{}, maxInFlight)

	return func(c *gin.Context) {
		select {
		case semaphore <- struct{}{}:
			// Pakai defer supaya slot kembali meski handler panic;
			// tanpa itu satu panic akan mengurangi kapasitas secara permanen.
			defer func() { <-semaphore }()
			c.Next()
		default:
			c.Header("Retry-After", "2")
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "Server sedang sibuk. Coba lagi sebentar lagi.",
			})
		}
	}
}