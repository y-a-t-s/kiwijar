package kiwijar

import (
	"net/http"
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

		oc := sm.m[c.Name]
		if oc == nil {
			sm.m[c.Name] = c
			return
		}

		// Overwrite the value instead of the pointer if cookie exists.
		sm.m[c.Name].Value = c.Value
	}()

	<-done
}

func (sm *siteMap) cookies() <-chan []*http.Cookie {
	out := make(chan []*http.Cookie, 1)

	go func() {
		defer close(out)

		sm.mx.Lock()
		defer sm.mx.Unlock()

		cs := make([]*http.Cookie, 0, len(sm.m))
		for _, c := range sm.m {
			cs = append(cs, c)
		}

		out <- cs
	}()

	return out
}
