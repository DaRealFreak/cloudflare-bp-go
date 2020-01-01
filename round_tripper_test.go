package cloudflarebp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
