package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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
		fmt.Fprintln(os.Stderr, "использование: hedgedcurl [-t SECONDS] URL [URL...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	urls := flag.Args()
	if len(urls) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Println("timeout:", timeout, "urls:", urls)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	res := fetch(ctx, urls[0])
	if res.err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", res.err)
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
