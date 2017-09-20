// uploader is a minimal service for accepting uploads over HTTP.

package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"html"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/crypto/acme/autocert"
)

var filesDir = ""

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for %q from %q\n", r.RequestURI, r.RemoteAddr)
	if r.Method != "POST" {
		log.Printf("Unexpected http method %q, returning StatusBadRequest\n", r.Method)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	p := path.Join(
		filesDir,
		path.Base(html.EscapeString(r.URL.Path)),
	)
	buf := bufio.NewReaderSize(r.Body, 1000)
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to write %q: %v\n", p, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	// TODO: Do read-write loop until EOF, to avoid allocating entire buf at once during Copy?
	n, err := io.Copy(f, buf)
	log.Printf("Wrote %d bytes to %q\n", n, p)
	if err := f.Close(); err != nil {
		log.Printf("Failed to close %q file: %v\n", p, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Thanks!")
}



func main() {
	filesDir = os.Getenv("UPLOADER_DIR")
	if filesDir == "" {
		filesDir = "/var/www"
	}
	http.HandleFunc("/", handler)
	addr := os.Getenv("UPLOADER_ADDR")
        if addr == "" {
                addr = ":8080"
        }
        s := &http.Server{ Addr: addr }
        if addr == ":443" {
                log.Printf("uploader accepting uploads to %q over TLS..", filesDir)
                m := autocert.Manager{
                        Prompt:     autocert.AcceptTOS,
                        Cache:      autocert.DirCache("/etc/secrets/acme/"),
                        HostPolicy: autocert.HostWhitelist("admin1.hkjn.me"),
                }
                s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
                log.Fatal(s.ListenAndServeTLS("", ""))
        } else {
                log.Printf("uploader accepting uploads to %q on plaintext HTTP on %q..\n", filesDir, addr)
                log.Fatal(s.ListenAndServe())
        }
	log.Printf("uploader accepting uploads to %q on %q..\n", filesDir, addr)
	panic(http.ListenAndServe(addr, nil))
}
