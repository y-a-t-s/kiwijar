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

func (kj *KiwiJar) Cookies(u *url.URL) []*http.Cookie {
	return kj.cm.cookies(u)
}

func (kj *KiwiJar) ParseString(u *url.URL, cookies string) error {
	if cookies == "" {
		return nil
	}

	cs, err := parseCookieString(cookies)
	if err != nil {
		return err
	}

	kj.SetCookies(u, cs)

	return nil
}

func (kj *KiwiJar) CookieString(u *url.URL) string {
	var b strings.Builder

	cookies := kj.Cookies(u)
	for _, c := range cookies {
		fmt.Fprintf(&b, "; %s=%s", c.Name, c.Value)
	}

	return strings.TrimPrefix(b.String(), "; ")
}

func (kj *KiwiJar) GetCookie(u *url.URL, name string) *http.Cookie {
	return kj.cm.loadSiteMap(u).getCookie(name)
}

func (kj *KiwiJar) SetCookie(u *url.URL, cookie *http.Cookie) {
	kj.cm.loadSiteMap(u).setCookie(cookie)
}

func (kj *KiwiJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	var wg sync.WaitGroup
	for _, c := range cookies {
		wg.Add(1)
		go func() {
			defer wg.Done()
			kj.SetCookie(u, c)
		}()
	}
	wg.Wait()
}

func parseCookieString(cookies string) ([]*http.Cookie, error) {
	sp := strings.Split(cookies, "; ")
	cs := make([]*http.Cookie, len(sp))

	for i, c := range sp {
		kv := strings.Split(c, "=")
		if len(kv) != 2 {
			return nil, errors.New("Invalid cookie string: " + cookies)
		}
		cs[i] = &http.Cookie{
			Name:  kv[0],
			Value: kv[1],
		}
	}

	return cs, nil
}
