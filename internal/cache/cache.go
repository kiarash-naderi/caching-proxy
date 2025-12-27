package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

type CachedResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

var (
	cacheMu sync.RWMutex
	cache   = make(map[string]CachedResponse)
)

func generateKey(r *http.Request) string {
	hasher := md5.New()

	hasher.Write([]byte(r.Method + r.URL.Path))

	if len(r.URL.Query()) > 0 {
		keys := make([]string, 0, len(r.URL.Query()))
		for k := range r.URL.Query() {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			hasher.Write([]byte(k + "=" + r.URL.Query().Get(k)))
		}
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func CachingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next() 
			return
		}

		key := generateKey(c.Request)

		cacheMu.RLock()
		cached, hit := cache[key]
		cacheMu.RUnlock()

		if hit {
			c.Header("X-Cache", "HIT")
			for k, vv := range cached.Header {
				for _, v := range vv {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.WriteHeader(cached.Status)
			c.Writer.Write(cached.Body)
			c.Abort()
			return
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			status:         http.StatusOK,
		}
		c.Writer = writer

		c.Next()

		if writer.status == http.StatusOK {
			cacheMu.Lock()
			cache[key] = CachedResponse{
				Status: writer.status,
				Header: writer.Header().Clone(),
				Body:   writer.body.Bytes(),
			}
			cacheMu.Unlock()
		}

		c.Header("X-Cache", "MISS")
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func ClearCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache = make(map[string]CachedResponse)
}