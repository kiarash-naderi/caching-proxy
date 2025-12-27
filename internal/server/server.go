package server

import (
	"fmt"
	"net/url"

	"caching-proxy/internal/cache"
	"caching-proxy/internal/proxy"

	"github.com/gin-gonic/gin"
)

func Run(port int, originStr string) {
	originURL, err := url.Parse(originStr)
	if err != nil {
		panic(fmt.Sprintf("invalid origin URL: %v", err))
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery()) 

	r.Use(cache.CachingMiddleware())
	r.Any("/*path", proxy.ProxyHandler(originURL))

	fmt.Printf("Caching Proxy running on :%d → %s\n", port, originStr)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}