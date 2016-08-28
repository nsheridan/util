package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	revokedURL = flag.String("url", "https://keyserver.evil.ie/revoked", "URL of ssh KRL")
	dest       = flag.String("dest", "/etc/ssh/revoked", "Output file destination")
)

func main() {
	flag.Parse()
	resp, err := http.Get(*revokedURL)
	if err != nil {
		fmt.Printf("Error retrieving revoked list: %s\n", err)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bad response from server")
		os.Exit(1)
	}
	krl, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	t, err := ioutil.TempFile("", "revoked")
	if err != nil {
		fmt.Printf("Error writing file: %s\n", err)
		os.Exit(1)
	}
	if _, err := t.Write(krl); err != nil {
		fmt.Printf("Error writing file: %s\n", err)
		os.Exit(1)
	}
	if err := t.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := os.Rename(t.Name(), *dest); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := os.Chmod(*dest, 0444); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
