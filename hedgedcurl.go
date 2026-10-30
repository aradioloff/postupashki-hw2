package main

import (
	"flag"
	"fmt"
	"os"
)

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
}
