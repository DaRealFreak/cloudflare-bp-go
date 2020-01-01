package cloudflarebp

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestApplyCloudFlareByPass(t *testing.T) {
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
