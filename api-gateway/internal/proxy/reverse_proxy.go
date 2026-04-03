package proxy

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	target *url.URL
}

func New(targetURL string) (*Proxy, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("error when setup proxy url: %v", err)
	}
	return &Proxy{target: u}, nil
}

func (p *Proxy) Forward() gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := httputil.NewSingleHostReverseProxy(p.target)

		c.Request.Host = p.target.Host

		rp.ServeHTTP(c.Writer, c.Request)
	}
}
