package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type req struct {
	method  string
	url     string
	headers map[string]string
	content []byte
}

func main() {
	repeats := 5

	var r req
	r.method = "PATCH"
	r.url = "http://localhost/v1/"
	r.headers = map[string]string{
		"Content-Type": "application/json",
		"X-Origin":     "test",
	}

	if f, err := os.Open("content.json"); err != nil {
		panic(err)
	} else {
		r.content, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("repeats: %d\n", repeats)

	wg := sync.WaitGroup{}
	for i := 0; i < repeats; i++ {
		wg.Add(1)
		go func(r req) {
			doRequest(r)
			wg.Done()
		}(r)
	}
	wg.Wait()
}

func doRequest(r req) {
	client := &http.Client{}

	req, err := http.NewRequest(r.method, r.url, bytes.NewReader(r.content))
	if err != nil {
		panic(err)
	}
	for k, v := range r.headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("status: %d\n", resp.StatusCode)
}
