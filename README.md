```markdown
# Caching Proxy Server

A lightweight **caching reverse proxy** built with Go and Gin.

It forwards HTTP requests to an upstream origin server, caches successful `GET` responses in memory, and serves subsequent identical requests directly from cache — adding an `X-Cache: HIT` or `X-Cache: MISS` header.

Perfect for development, testing, or reducing load on external APIs.

## Features

- In-memory thread-safe cache
- Caches only `GET` requests with status `200 OK`
- Adds `X-Cache: HIT` / `X-Cache: MISS` header
- CLI command to clear the entire cache
- Clean modular architecture
- Built with standard library `httputil.ReverseProxy` and Gin

## Quick Start

Build the binary:

```bash
go build -o caching-proxy ./cmd
```

Run the proxy:

```bash
./caching-proxy --port 3000 --origin https://dummyjson.com
```

The proxy will listen on `http://localhost:3000` and forward all requests to the origin.

### Example

```bash
# First request → MISS (fetched from origin)
curl -i http://localhost:3000/products/1

# Second request → HIT (served from cache)
curl -i http://localhost:3000/products/1
```

### Clear Cache

```bash
./caching-proxy clear-cache
# Output: Cache cleared successfully!
```

## CLI Usage

```bash
./caching-proxy --help
```

```
NAME:
   caching-proxy - A simple caching reverse proxy

USAGE:
   caching-proxy [command]

COMMANDS:
   clear-cache    Clear the entire cache

FLAGS:
   --port value     Proxy listening port
   --origin value   Origin server URL (e.g. https://dummyjson.com)
```

## Project Structure

```
caching-proxy/
├── cmd/
│   └── main.go              # CLI entry point
├── internal/
│   ├── cache/   cache.go    # Caching middleware & logic
│   ├── proxy/   proxy.go    # Reverse proxy handler
│   └── server/  server.go   # Gin server setup
├── go.mod
├── go.sum
└── README.md
```

## Possible Future Improvements

- Add TTL / expiration for cache entries
- Respect `Cache-Control` headers from origin
- Optional Redis backend
- Prometheus metrics
- TLS support

## License

MIT License – free to use, modify, and distribute.



