package kiwijar

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// An http cookiejar implementation that doesn't suck ass.
type KiwiJar struct {
	cm cookieMap
}

// Implements the Cookies method of the http.CookieJar interface.
// Returns empty slice if the URL's scheme is not http or https.
func (kj *KiwiJar) Cookies(u *url.URL) (cookies []*http.Cookie) {
	switch u.Scheme {
	case "http", "https":
		return kj.cm.cookies(u)
	default:
		return
	}
}

func (kj *KiwiJar) ParseString(u *url.URL, cookies string) error {
	switch u.Scheme {
	case "http", "https":
	default:
		return nil
	}

	if cookies == "" {
		return nil
	}

	cs, err := parseCookieString(u, cookies)
	if err != nil {
		return err
	}

	kj.SetCookies(u, cs)

	return nil
}

func (kj *KiwiJar) HeaderString(u *url.URL) string {
	var b strings.Builder

	cookies := kj.Cookies(u)
	for _, c := range cookies {
		fmt.Fprintf(&b, "; %s=%s", c.Name, c.Value)
	}

	return strings.TrimPrefix(b.String(), "; ")
}

func (kj *KiwiJar) GetCookie(u *url.URL, name string) *http.Cookie {
	switch u.Scheme {
	case "http", "https":
		return kj.cm.loadSiteMap(u).getCookie(name)
	default:
		return nil
	}
}

func (kj *KiwiJar) SetCookie(u *url.URL, cookie *http.Cookie) {
	switch u.Scheme {
	case "http", "https":
		kj.cm.loadSiteMap(u).setCookie(cookie)
	}
}

func (kj *KiwiJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	switch u.Scheme {
	case "http", "https":
	default:
		return
	}

	var (
		sm = kj.cm.loadSiteMap(u)
		wg sync.WaitGroup
	)

	for _, c := range cookies {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm.setCookie(c)
		}()
	}
	wg.Wait()
}

func parseCookieString(u *url.URL, cookies string) ([]*http.Cookie, error) {
	sp := strings.Split(cookies, "; ")
	cs := make([]*http.Cookie, len(sp))

	var (
		host   = u.Hostname()
		path   = u.Path
		secure = u.Scheme == "https"
	)

	if path == "" {
		path = "/"
	}

	for i, c := range sp {
		kv := strings.Split(c, "=")
		if len(kv) != 2 {
			return nil, errors.New("Invalid cookie string: " + cookies)
		}
		cs[i] = &http.Cookie{
			Domain: host,
			Name:   kv[0],
			Path:   path,
			Secure: secure,
			Value:  kv[1],
		}
	}

	return cs, nil
}
