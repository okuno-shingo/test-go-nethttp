package main

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        0,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: 3 * time.Second,
	}

	servers := []string{"http://server1:8080", "http://server2:8080"}

	for {
		for _, server := range servers {
			time.Sleep(50 * time.Millisecond)
			go func(server string) {
				reqID := uuid.New().String()

				operation := func() error {
					fmt.Printf("[%s] Request to %s\n", reqID, server)

					response, err := client.Get(server)
					if err != nil {
						fmt.Printf("[%s] Error fetching %s: %v\n", reqID, server, err)
						return err
					}
					defer func(Body io.ReadCloser) {
						err := Body.Close()
						if err != nil {
							fmt.Printf("[%s] Error closing body: %v\n", reqID, err)
						}
					}(response.Body)

					_, _ = io.ReadAll(response.Body)

					if response.StatusCode == 500 {
						return fmt.Errorf("bad response from server")
					}

					fmt.Printf("[%s] Response from %s: %s\n", reqID, server, response.Status)
					return nil
				}

				notify := func(err error, t time.Duration) {
					fmt.Printf("[%s] Failed with %s, retrying in %s\n", reqID, server, err)
				}

				bo := backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 3)

				err := backoff.RetryNotify(operation, bo, notify)
				if err != nil {
					fmt.Printf("[%s] Operation failed after retries: %v\n", reqID, err)
				}
			}(server)
		}
	}
}
