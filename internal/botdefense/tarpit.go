package botdefense

import (
	"time"

	"github.com/gin-gonic/gin"
)

// slowly drip-feeds a response to waste bot resources
func Tarpit(c *gin.Context, duration time.Duration, chunkDelay time.Duration) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Connection", "keep-alive")
	c.Status(200)

	chunks := max(int(duration/chunkDelay), 1)

	opener := `<!DOCTYPE html><html><head><title>Loading...</title></head><body><p>Please wait`

	if _, err := c.Writer.Write([]byte(opener)); err != nil {
		return
	}
	c.Writer.Flush()

	for range chunks {
		select {
		case <-c.Request.Context().Done():
			return
		case <-time.After(chunkDelay):
			if _, err := c.Writer.Write([]byte(".")); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}

	c.Writer.Write([]byte("</p></body></html>")) //nolint:errcheck,gosec
	c.Writer.Flush()
}

// slowly drip-feeds a JSON response
func TarpitJSON(c *gin.Context, duration time.Duration, chunkDelay time.Duration) {
	c.Header("Content-Type", "application/json")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(200)

	chunks := max(int(duration/chunkDelay), 1)

	if _, err := c.Writer.Write([]byte(`{"status":"loading","data":[`)); err != nil {
		return
	}
	c.Writer.Flush()

	for i := range chunks {
		select {
		case <-c.Request.Context().Done():
			return
		case <-time.After(chunkDelay):
			comma := ""
			if i > 0 {
				comma = ","
			}
			if _, err := c.Writer.Write([]byte(comma + `{}`)); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}

	c.Writer.Write([]byte(`]}`)) //nolint:errcheck,gosec
	c.Writer.Flush()
}
