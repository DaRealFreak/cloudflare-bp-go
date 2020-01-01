package cloudflarebp

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/proxy"
)

func TestApplyCloudFlareByPassDefaultClient(t *testing.T) {
	client := http.DefaultClient

	res, err := client.Get("https://www.patreon.com/login")
	assert.New(t).NoError(err)
	assert.New(t).Equal(403, res.StatusCode)

	// apply our bypass for request headers and client TLS configurations
	http.DefaultClient.Transport = AddCloudFlareByPass(http.DefaultClient.Transport)

	res, err = client.Get("https://www.patreon.com/login")
	assert.New(t).NoError(err)
	assert.New(t).Equal(200, res.StatusCode)
}

func TestApplyCloudFlareByPassDefinedTransport(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{},
	}

	// if the client requests something before applying the fix some configurations are applied already
	// and our ByPass won't work anymore, so we have to apply our ByPass as the first thing
	client.Transport = AddCloudFlareByPass(client.Transport)

	res, err := client.Get("https://www.patreon.com/login")
	assert.New(t).NoError(err)
	assert.New(t).Equal(200, res.StatusCode)
}

// TestAddCloudFlareByPassSocksProxy tests the CloudFlare bypass while we're using a SOCK5 proxy transport layer
func TestAddCloudFlareByPassSocksProxy(t *testing.T) {
	auth := proxy.Auth{
		User:     os.Getenv("PROXY_USER"),
		Password: os.Getenv("PROXY_PASS"),
	}

	dialer, err := proxy.SOCKS5(
		"tcp",
		fmt.Sprintf("%s:1080", os.Getenv("PROXY_HOST")),
		&auth,
		proxy.Direct,
	)
	assert.New(t).NoError(err)

	client := &http.Client{
		Transport: &http.Transport{Dial: dialer.Dial},
	}

	// if the client requests something before applying the fix some configurations are applied already
	// and our ByPass won't work anymore, so we have to apply our ByPass as the first thing
	client.Transport = AddCloudFlareByPass(client.Transport)

	res, err := client.Get("https://www.patreon.com/login")
	assert.New(t).NoError(err)
	assert.New(t).Equal(200, res.StatusCode)
}

// TestAddCloudFlareByPassHTTPProxy tests the CloudFlare bypass while we're using a HTTP proxy transport layer
func TestAddCloudFlareByPassHTTPProxy(t *testing.T) {
	proxyURL, _ := url.Parse(
		fmt.Sprintf(
			"http://%s:%s@%s:80",
			url.QueryEscape(os.Getenv("PROXY_USER")), url.QueryEscape(os.Getenv("PROXY_PASS")),
			url.QueryEscape(os.Getenv("PROXY_HOST")),
		),
	)

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}

	// if the client requests something before applying the fix some configurations are applied already
	// and our ByPass won't work anymore, so we have to apply our ByPass as the first thing
	client.Transport = AddCloudFlareByPass(client.Transport)

	res, err := client.Get("https://www.patreon.com/login")
	assert.New(t).NoError(err)
	assert.New(t).Equal(200, res.StatusCode)
}
