package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(originURL *url.URL) gin.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(originURL)

	return func(c *gin.Context) {
		req := c.Request
		req.URL.Scheme = originURL.Scheme
		req.URL.Host = originURL.Host
		req.Host = originURL.Host 

		proxy.ServeHTTP(c.Writer, req)
	}
}