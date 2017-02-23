package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"golang.org/x/net/http2"
)

var (
	url        = flag.String("url", "https://http2.golang.org/", "URL to fetch")
	skipVerify = flag.Bool("skip_verify", false, "Verify TLS certificates")
)

func fatal(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}

func main() {
	flag.Parse()
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: *skipVerify,
		},
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		fatal(err)
	}
	cl := &http.Client{
		Transport: transport,
	}
	resp, err := cl.Get(*url)
	if err != nil {
		fatal(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "Protocol:\t%s\n", resp.Proto)
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		fmt.Fprintf(w, "Certificate Common Name:\t%s\n", resp.TLS.PeerCertificates[0].Subject.CommonName)
		fmt.Fprintf(w, "Certificate Alt Names:\t%s\n", strings.Join(resp.TLS.PeerCertificates[0].DNSNames, ", "))
	}
	fmt.Fprintf(w, "\nHeaders:\n")
	for h, v := range resp.Header {
		fmt.Fprintf(w, "%s:\t%s\n", h, strings.Join(v, ""))
	}
	w.Flush()
}
