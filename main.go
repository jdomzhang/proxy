package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
)

var PrintableTypes map[string]bool = map[string]bool{
	"application/atom+xml":              true,
	"application/ecmascript":            true,
	"application/json":                  true,
	"application/javascript":            true,
	"application/rdf+xml":               true,
	"application/rss+xml":               true,
	"application/soap+xml":              true,
	"application/xhtml+xml":             true,
	"application/xml":                   true,
	"application/xml-dtd":               true,
	"application/x-www-form-urlencoded": true,
	"text/css":                          true,
	"text/csv":                          true,
	"text/html":                         true,
	"text/javascript":                   true,
	"text/plain":                        true,
	"text/vcard":                        true,
	"text/xml":                          true,
}

func main() {
	port := flag.Int("port", 7777, "The port to listen on.")
	help := flag.Bool("help", false, "Show this help message, then exit.")

	flag.Parse()
	flag.Usage = usage

	if *help {
		usage()
		os.Exit(0)
	}

	fmt.Printf("Proxy listening on :%d ...\n", *port)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

			fmt.Println("requesting...")

			body := isPrintable(req.Header)
			reqDump, _ := httputil.DumpRequest(req, body)
			if !body {
				reqDump = append(reqDump, []byte("BINARY\n\n")...)
			}

			fmt.Printf("\n---\n\n")
			fmt.Printf("> %s:\n%s", req.RequestURI, reqDump)
			// fmt.Printf("< %s:\n%s", flag.Arg(0), req.RemoteAddr, resDump)
			fmt.Printf("\n")

			return req, nil
		})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), proxy))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: proxy [PORT]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.VisitAll(func(flag *flag.Flag) {
		fmt.Fprintf(os.Stderr, "  --%s (%v)  \t%s\n", flag.Name, flag.DefValue, flag.Usage)
	})

	fmt.Fprintf(os.Stderr, "\nFor more information, see https://github.com/jdomzhang/proxy.\n")
}

func isPrintable(header http.Header) bool {
	mimeType := header.Get(http.CanonicalHeaderKey("content-type"))
	mimeType = strings.SplitN(mimeType, ";", 2)[0]
	mimeType = strings.TrimSpace(mimeType)

	if mimeType == "" {
		return true
	}

	return PrintableTypes[mimeType]
}
