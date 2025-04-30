package utilities

import (
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

/**
# Configure Nginx: Ensure that Nginx is set up to pass the X-Forwarded-For header.
This is typically done with the following configuration:

server {
        ...
        location / {
            proxy_pass http://your_gin_backend;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Real-IP $remote_addr;
            ...
        }
}
*/

func GetClientIP(c *gin.Context) string {
	// Check the X-Forwarded-For header
	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		// Take the first IP in the list, which is the client's IP
		return strings.TrimSpace(ips[0])
	}

	// Fall back to X-Real-IP if X-Forwarded-For is not present
	xRealIP := c.Request.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr if no proxy headers are present
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}
