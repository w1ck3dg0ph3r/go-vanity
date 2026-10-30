// go-vanity serves Go vanity URLs.
// It reads list of modules from vanities.list file, and watches the vanities
// file for changes.
//
// Usage:
//
//	go-vanity [flags]
//
// The flags are:
//
//	-listen {addr}
//		Host:port to listen on. Default is ":5000".
//
//	-vanities {path}
//		Path to the vanities.list file. Default is "./vanities.list".
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	"go.w1ck3dg0ph3r.dev/vanity/forges"
)

func main() {
	f := ParseFlags()

	vs, err := ReadVanities(f.ListFile)
	if err != nil {
		Fatal(fmt.Errorf("read vanities: %w", err))
	}
	fmt.Printf("read vanities:\n%s\n", vs)

	h := &DynamicHandler{}
	h.Set(vs.Handler())

	w := WatchVanities(f.ListFile, h)
	defer w.Close()

	if err := http.ListenAndServe(f.Listen, Log(h)); err != nil {
		Fatal(fmt.Errorf("listen: %w", err))
	}
}

type Flags struct {
	ListFile string
	Listen   string
}

func ParseFlags() Flags {
	var f Flags
	flag.StringVar(&f.Listen, "listen", ":5000", "Address to listen on")
	flag.StringVar(&f.ListFile, "vanities", "./vanities.list", "Vanities list file path")
	flag.Parse()
	return f
}

func Fatal(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}

func ReadVanities(fn string) (Vanities, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var vanities Vanities

	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		v, err := forges.Parse(l)
		if err != nil {
			return nil, fmt.Errorf("parse vanity: %w", err)
		}
		vanities = append(vanities, v)
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("scan file: %w", err)
	}

	return vanities, nil
}

func WatchVanities(fn string, h *DynamicHandler) *fsnotify.Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		Fatal(fmt.Errorf("watch vanities: %w", err))
	}
	go func() {
		for {
			select {
			case e, ok := <-w.Events:
				if !ok {
					return
				}
				if e.Op&fsnotify.Write > 0 {
					vs, err := ReadVanities(fn)
					if err != nil {
						fmt.Printf("error: %v\n", fmt.Errorf("read vanities: %w", err))
						continue
					}
					fmt.Printf("read vanities:\n%s\n", vs)
					h.Set(vs.Handler())
				}
			case err := <-w.Errors:
				fmt.Printf("watch error: %v\n", err)
			}
		}
	}()
	w.Add(fn)
	return w
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	})
}

type DynamicHandler struct {
	mu sync.RWMutex
	h  http.Handler
}

func (d *DynamicHandler) Set(h http.Handler) {
	d.mu.Lock()
	d.h = h
	d.mu.Unlock()
}

func (d *DynamicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	h := d.h
	d.mu.RUnlock()
	h.ServeHTTP(w, r)
}

type Vanities []forges.Vanity

func (vs Vanities) Handler() http.Handler {
	m := http.NewServeMux()
	for _, v := range vs {
		h := VanityHandler(v)
		m.Handle("GET /"+v.Module(), h)
		m.Handle("GET /"+v.Module()+"/", h)
	}
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	m.Handle("/_list", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		for _, v := range vs {
			w.Write(fmt.Appendf(nil, "%s/%s -> %s\n", r.Host, v.Module(), v.VCSHome()))
		}
	}))
	return m
}

func VanityHandler(v forges.Vanity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkg := strings.TrimSuffix(r.Host+r.URL.Path, "/")
		root := r.Host + "/" + v.Module()

		if r.URL.Query().Get("go-get") != "1" {
			w.Header().Add("Location", "https://pkg.go.dev/"+pkg)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		res := &bytes.Buffer{}

		res.WriteString(`<meta name="go-import" content="`)
		res.WriteString(root)
		res.WriteString(" git ")
		res.WriteString(v.VCSHome())
		res.WriteString("\">\n")

		if v.VCSSourceDir() != "" && v.VCSSourceFile() != "" {
			res.WriteString(`<meta name="go-source" content="`)
			res.WriteString(root)
			res.WriteString(" ")
			res.WriteString(v.VCSHome())
			res.WriteString(" ")
			res.WriteString(v.VCSSourceDir())
			res.WriteString(" ")
			res.WriteString(v.VCSSourceFile())
			res.WriteString("\">\n")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.WriteHeader(200)
		w.Write(res.Bytes())
	})
}

func (vs Vanities) String() string {
	var s strings.Builder
	for i, v := range vs {
		s.WriteString(v.String())
		if i < len(vs)-1 {
			s.WriteRune('\n')
		}
	}
	return s.String()
}
