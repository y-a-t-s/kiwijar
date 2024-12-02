package kiwijar

import (
	"log"
	"net/url"
	"os"
	"testing"
)

const TEST_HOST = "kiwifarms.net"

func TestKiwiJar(t *testing.T) {
	tc := os.Getenv("TEST_COOKIES")

	u, err := url.Parse("https://" + TEST_HOST)
	if err != nil {
		t.Error(err)
	}

	jar := KiwiJar{}
	err = jar.ParseString(u, tc)
	if err != nil {
		t.Error(err)
	}

	log.Printf("Cookies from jar: %s\n", jar.CookieString(u))
}
