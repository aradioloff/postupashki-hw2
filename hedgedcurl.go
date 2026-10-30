package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type result struct {
	url    string
	proto  string
	status string
	header http.Header
	body   []byte
	err    error
}

func fetch(ctx context.Context, url string) result {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result{url: url, err: err}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result{url: url, err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result{url: url, err: err}
	}
	return result{
		url:    url,
		proto:  resp.Proto,
		status: resp.Status,
		header: resp.Header,
		body:   body,
	}
}
func main() {
	var timeout int
	flag.IntVar(&timeout, "t", 15, "таймаут в секундах")
	flag.IntVar(&timeout, "timeout", 15, "таймаут в секундах")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Использование: hedgedcurl [-t SECONDS] URL [URL...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	urls := flag.Args()
	if len(urls) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	results := make(chan result, len(urls))
	for _, u := range urls {
		if !strings.Contains(u, "://") {
			u = "https://" + u
		}
		go func(u string) {
			results <- fetch(ctx, u)
		}(u)
	}
	var res result
	var lastErr error
	timedOut := false
	got := false
	for range urls {
		r := <-results
		if r.err != nil {
			lastErr = r.err
			if errors.Is(r.err, context.DeadlineExceeded) {
				timedOut = true
			}
			continue
		}
		res = r
		got = true
		cancel()
		break
	}
	if !got {
		if timedOut {
			fmt.Fprintln(os.Stderr, "ошибка: превышено время ожидания")
			os.Exit(228)
		}
		fmt.Fprintln(os.Stderr, "все запросы завершились ошибкой:", lastErr)
		os.Exit(1)
	}
	fmt.Println(res.proto, res.status)
	for name, values := range res.header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", name, v)
		}
	}
	fmt.Println()
	os.Stdout.Write(res.body)
}
