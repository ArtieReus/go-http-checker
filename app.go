package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	appName  = "go-http-checker"
	version  = "1.0.0"
	checkURL = ""
)

func main() {
	app := cli.NewApp()

	app.Name = appName
	app.Version = version
	app.Authors = []cli.Author{
		{
			Name:  "Arturo Reuschenbach Puncernau",
			Email: "a.reuschenbach.puncernau@sap.com",
		},
	}
	app.Usage = "check http connections"
	app.Action = runChecker
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url,u",
			Usage: "url to check",
			Value: "http://www.google.com",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// private

func runChecker(c *cli.Context) {
	if c.GlobalString("url") != "" {
		checkURL = c.GlobalString("url")
	} else {
		log.Fatalf("Url not provided")
	}

	log.Infof("Checking for URL: %s", checkURL)

	// client
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			req.Header["Accept"] = via[0].Header["Accept"]
			req.Header["User-Agent"] = via[0].Header["User-Agent"]
			return nil
		},
	}

	// request
	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		log.Fatalf("failed to create http request: %s", err)
	}
	req.Header["User-Agent"] = []string{appName + " " + version}

	// pretty print req
	log.Println()
	formatRequest(req)

	// response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	log.Println()
	log.Println()
	formatResponse(resp)
}

func formatResponse(resp *http.Response) {
	log.Println("RESPONSE:")
	log.Printf("Status: %v Status code: %v", resp.Status, resp.StatusCode)

	// Loop through headers
	log.Println()
	log.Println("Response Headers:")
	for name, headers := range resp.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			log.Printf("%v: %v", name, h)
		}
	}
	// print dump
	log.Println()
	log.Println("Dump:")
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Warnf("response can not be dumped: %s", err)
	}

	// add raw request
	log.Printf("%q", dump)
}

func formatRequest(r *http.Request) {
	log.Println("REQUEST:")
	// Add the request string
	log.Printf("%v %v %v", r.Method, r.URL, r.Proto)
	// Add the host
	log.Printf("Host: %v", r.Host)

	// Loop through headers
	log.Println()
	log.Println("Request Headers:")
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			log.Printf("%v: %v", name, h)
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Warnf("Post form could not be parsed: %s", err)
		} else {
			log.Println(r.Form.Encode())
		}
	}

	// print dump
	log.Println()
	log.Println("Dump:")
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Warnf("request can not be dumped: %s", err)
	}

	// add raw request
	log.Printf("%q", dump)
}
