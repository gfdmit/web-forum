package proxy

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	rp *httputil.ReverseProxy
}

func New(targetURL string) (*Proxy, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("proxy.New: %w", err)
	}
	return &Proxy{rp: httputil.NewSingleHostReverseProxy(u)}, nil
}

func (p *Proxy) Forward() gin.HandlerFunc {
	return func(c *gin.Context) {
		p.rp.ServeHTTP(c.Writer, c.Request)
	}
}
