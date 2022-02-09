package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type req struct {
	method  string
	url     string
	headers map[string]string
	body    []byte
}

var (
	client = &http.Client{}
)

func main() {
	u := flag.String("u", "", "request URL")
	b := flag.String("b", "", "request body")
	m := flag.String("m", "POST", "HTTP method (PUT, POST, PATCH, GET, ...)")
	p := flag.Int("p", 10, "count of parallel requests")
	h := flag.String("h", "Content-Type: application/json", "request headers, multiple separated by \\n")
	debug := flag.Bool("debug", false, "debug mode (print response body etc.)")
	flag.Parse()

	// validate flags
	if *u == "" || *b == "" {
		flag.Usage()
		return
	}

	// build request
	var r req
	r.method = *m
	r.url = *u
	r.body = []byte(*b)

	r.headers = make(map[string]string)
	if *h != "" {
		for _, h := range strings.Split(*h, "\\n") {
			kv := strings.Split(h, ":")
			if len(kv) != 2 {
				flag.Usage()
				return
			}
			r.headers[kv[0]] = kv[1]
		}
	}

	// fire parallel requests
	wg := sync.WaitGroup{}
	for i := 0; i < *p; i++ {
		wg.Add(1)
		go func(r req) {
			if err := doRequest(r, *debug); err != nil {
				fmt.Printf("error: %s\n", err)
			}
			wg.Done()
		}(r)
	}
	wg.Wait()
}

func doRequest(r req, debug bool) error {
	req, err := http.NewRequest(r.method, r.url, bytes.NewReader(r.body))
	if err != nil {
		return err
	}
	for k, v := range r.headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("status: %d\n", resp.StatusCode)

	if debug {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf(">> body: %s\n", string(b))
	}
	return nil
}
