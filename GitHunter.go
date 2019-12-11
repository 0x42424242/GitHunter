package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Settings struct {
	filename string
	quiet    bool
}

type Query struct {
	hostname string
	status   int
	error    bool
	git      bool
	https    bool
}

func gitWorker(q Query, rChan chan Query, c *http.Client, wg *sync.WaitGroup) {
	var res *http.Response
	var err error

	if q.https {
		res, err = c.Get(fmt.Sprintf("https://%s/.git/config", q.hostname))
	} else {
		res, err = c.Get(fmt.Sprintf("http://%s/.git/config", q.hostname))
	}

	if err != nil {
		q.error = true

		rChan <- q

		wg.Done()
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		q.error = true

		rChan <- q

		wg.Done()
		return
	}

	if strings.Contains(string(body), "[core]") {
		q.git = true
		q.error = false
		q.status = res.StatusCode

		rChan <- q

		wg.Done()
		return
	} else {
		q.git = false
		q.error = false
		q.status = res.StatusCode

		rChan <- q

		wg.Done()
		return
	}
}

func resultPrinter(r chan Query, s Settings) {
	for i := range r {
		if !s.quiet {
			if i.error {
				fmt.Printf("\033[1;33m[!]\033[0m %s\n", i.hostname)
			} else if i.git {
				fmt.Printf("\033[1;32m[+]\033[0m [%d] [HTTPS: %t] %s\n", i.status, i.https, i.hostname)
			} else {
				fmt.Printf("\033[1;31m[-]\033[0m [%d] %s\n", i.status, i.hostname)
			}
		} else if i.git {
			fmt.Printf("\n\033[1;32m[+]\033[0m [%d] [HTTPS: %t] %s\n", i.status, i.https, i.hostname)
		} else {
			fmt.Printf(".")
		}
	}
}

func main() {
	var settings Settings

	flag.StringVar(&settings.filename, "filename", "", "Filename of hosts. 1 per line.")
	flag.BoolVar(&settings.quiet, "quiet", false, "Quiet print or not")
	flag.Parse()

	f, err := os.Open(settings.filename)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	s := bufio.NewScanner(f)
	rChan := make(chan Query)

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
	}

	for s.Scan() {
		wg.Add(2)
		q := Query{
			hostname: s.Text(),
			https:    true,
		}
		go gitWorker(q, rChan, client, &wg)

		q.https = false
		go gitWorker(q, rChan, client, &wg)
	}

	go resultPrinter(rChan, settings)

	defer close(rChan)
	wg.Wait()
}
