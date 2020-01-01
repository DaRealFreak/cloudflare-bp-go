package cloudflarebp

import (
	"crypto/tls"
	browser "github.com/EDDYCJY/fake-useragent"
	"net/http"
)

// cloudFlareRoundTripper is a custom round tripper add the validated request headers
type cloudFlareRoundTripper struct {
	inner     http.RoundTripper
	userAgent string
}

// addCloudFlareByPass returns a round tripper adding the required headers for the CloudFlare checks
// and updates the TLS configuration of the passed inner transport
func AddCloudFlareByPass(inner http.RoundTripper) http.RoundTripper {
	if trans, ok := inner.(*http.Transport); ok {
		trans.TLSClientConfig = getCloudFlareTLSConfiguration()
	}

	return &cloudFlareRoundTripper{
		inner:     inner,
		userAgent: browser.Firefox(),
	}
}

// RoundTrip adds the required request headers to pass CloudFlare checks
func (ug *cloudFlareRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// Accept-Encoding header not needed here since the http library will add gzip automatically if not set manually
	// would be required for porting this for other libraries
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	r.Header.Set("User-Agent", ug.userAgent)

	// in case we don't have an inner transport layer from the round tripper
	if ug.inner == nil {
		return (&http.Transport{
			TLSClientConfig: getCloudFlareTLSConfiguration(),
		}).RoundTrip(r)
	}

	return ug.inner.RoundTrip(r)
}

func getCloudFlareTLSConfiguration() *tls.Config {
	return &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519},
	}
}
