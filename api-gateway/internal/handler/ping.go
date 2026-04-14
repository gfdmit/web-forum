package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func pingHandler(services map[string]string) gin.HandlerFunc {
	client := &http.Client{Timeout: 3 * time.Second}
	return func(c *gin.Context) {
		results := make(map[string]string, len(services))
		var mu sync.Mutex
		var wg sync.WaitGroup

		for name, url := range services {
			wg.Add(1)
			go func(name, url string) {
				defer wg.Done()
				status := "ok"
				resp, err := client.Get(url)
				if err != nil {
					status = "unavailable"
				} else {
					resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						status = "unavailable"
					}
				}
				mu.Lock()
				results[name] = status
				mu.Unlock()
			}(name, url)
		}

		wg.Wait()
		c.JSON(http.StatusOK, results)
	}
}
