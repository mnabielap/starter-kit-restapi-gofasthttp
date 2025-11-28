package middleware

import (
	"log"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
)

// Logger is a middleware that logs the incoming request details and duration
func Logger(c *routing.Context) error {
	start := time.Now()

	// Process the request
	err := c.Next()

	// Log details
	log.Printf("[HTTP] %s %s | Status: %d | Duration: %v", 
		string(c.Method()), 
		string(c.Path()), 
		c.Response.StatusCode(), 
		time.Since(start),
	)

	return err
}