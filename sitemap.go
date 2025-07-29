package kiwijar

import (
	"net/http"
	"net/url"
	"slices"
	"sync"
)

type siteMap struct {
	m  map[string]*http.Cookie
	mx sync.Mutex
}

func (sm *siteMap) getCookie(name string) *http.Cookie {
	out := make(chan *http.Cookie, 1)

	go func() {
		defer close(out)

		sm.mx.Lock()
		defer sm.mx.Unlock()

		out <- sm.m[name]
	}()

	return <-out
}

func (sm *siteMap) setCookie(c *http.Cookie) {
	done := make(chan struct{}, 1)

	go func() {
		defer close(done)

		sm.mx.Lock()
		defer sm.mx.Unlock()

		sm.m[c.Name] = c
	}()

	<-done
}

func (sm *siteMap) cookies(u *url.URL) <-chan []*http.Cookie {
	out := make(chan []*http.Cookie, 1)

	path := u.Path
	if path == "" {
		path = "/"
	}

	go func() {
		defer close(out)

		sm.mx.Lock()
		defer sm.mx.Unlock()

		cs := make([]*http.Cookie, 0, len(sm.m))
		for _, c := range sm.m {
			if u.Hostname() == c.Domain && path == c.Path && (u.Scheme == "https" || !c.Secure) {
				cs = append(cs, c)
			}
		}

		out <- slices.Clip(cs)
	}()

	return out
}
