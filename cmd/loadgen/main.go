package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type requestBody struct {
	Numbers []int64 `json:"numbers"`
	Repeat  int     `json:"repeat"`
}

func main() {
	var (
		url     = flag.String("url", "http://127.0.0.1:8080/sum", "target URL")
		workers = flag.Int("workers", 16, "concurrent workers")
		n       = flag.Int("n", 20000, "total requests")

		numbers = flag.Int("numbers", 1000, "numbers in payload")
		repeat  = flag.Int("repeat", 200, "repeat summation inside handler")
		seed    = flag.Int64("seed", 42, "rng seed")

		timeout = flag.Duration("timeout", 5*time.Second, "request timeout")
	)
	flag.Parse()

	// Генерим payload один раз и переиспользуем байты во всех запросах
	payload := requestBody{
		Numbers: make([]int64, *numbers),
		Repeat:  *repeat,
	}
	rng := rand.New(rand.NewSource(*seed))
	for i := range payload.Numbers {
		payload.Numbers[i] = int64(rng.Intn(1000))
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 200,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Timeout:   *timeout,
		Transport: transport,
	}

	tasks := make(chan int, *n)
	for i := 0; i < *n; i++ {
		tasks <- i
	}
	close(tasks)

	var okCount int64
	var errCount int64

	latencies := make([]time.Duration, 0, *n)
	var latMu sync.Mutex

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for range tasks {
				req, _ := http.NewRequest(http.MethodPost, *url, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				t0 := time.Now()
				resp, err := client.Do(req)
				lat := time.Since(t0)

				if err != nil {
					atomic.AddInt64(&errCount, 1)
					continue
				}

				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()

				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					atomic.AddInt64(&okCount, 1)
				} else {
					atomic.AddInt64(&errCount, 1)
				}

				latMu.Lock()
				latencies = append(latencies, lat)
				latMu.Unlock()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	pct := func(p float64) time.Duration {
		if len(latencies) == 0 {
			return 0
		}
		idx := int(float64(len(latencies)-1) * p)
		return latencies[idx]
	}

	fmt.Printf("url=%s workers=%d requests=%d ok=%d err=%d elapsed=%s rps=%.1f\n",
		*url, *workers, *n, okCount, errCount, elapsed.Truncate(time.Millisecond), float64(*n)/elapsed.Seconds())
	fmt.Printf("p50=%s p95=%s p99=%s\n", pct(0.50), pct(0.95), pct(0.99))
}
