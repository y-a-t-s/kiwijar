package kiwijar

import (
	"net/http"
	"net/url"
	"sync"
)

type cookieMap struct {
	m sync.Map // Stores *siteMap
}

func (cm *cookieMap) loadSiteMap(u *url.URL) *siteMap {
	sm, _ := cm.m.LoadOrStore(u.Hostname(), &siteMap{
		m: make(map[string]*http.Cookie),
	})

	return sm.(*siteMap)
}

func (cm *cookieMap) cookies(u *url.URL) []*http.Cookie {
	return <-cm.loadSiteMap(u).cookies()
}
