package main

import (
	"fmt"
	"log"
	"os"
	"caching-proxy/internal/cache"
	"caching-proxy/internal/server"
	"github.com/urfave/cli/v2" // optional, for nicer CLI
)

func main() {
	app := &cli.App{
		Name:  "caching-proxy",
		Usage: "A simple caching proxy server with Gin",
		Commands: []*cli.Command{
			{
				Name:  "clear-cache",
				Usage: "Clear the entire cache",
				Action: func(c *cli.Context) error {
					cache.ClearCache()
					fmt.Println("Cache cleared successfully!")
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() > 0 || c.NumFlags() < 2 {
				return fmt.Errorf("Usage: caching-proxy --port <number> --origin <url>")
			}
			port := c.Int("port")
			origin := c.String("origin")

			if port == 0 || origin == "" {
				return fmt.Errorf("--port and --origin are required")
			}

			server.Run(port, origin)
			return nil
		},
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "port", Usage: "Proxy port"},
			&cli.StringFlag{Name: "origin", Usage: "Origin server address (e.g., http://dummyjson.com)"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
