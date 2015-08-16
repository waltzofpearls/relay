package rapi_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/waltzofpearls/relay-api/rapi"
)

func TestEndpointUnchanged(t *testing.T) {

	var requestContent string

	expectedResult := `test`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _ := ioutil.ReadAll(r.Body)
		requestContent = string(req)
		w.Write([]byte(expectedResult))
	}))
	defer ts.Close()

	conf := rapi.NewConfig()
	conf.Backend.Address = strings.TrimPrefix(ts.URL, "http://")

	api := rapi.New(conf)
	require.NotNil(t, api)

	ep := rapi.NewEndpoint(api, "POST", "/foo")
	assert.NotNil(t, ep)

	fixture := `{"One":"this is the one", "Two":"this is the second"}`
	req, err := http.NewRequest("POST", "/foo", strings.NewReader(fixture))
	require.Nil(t, err)
	require.NotNil(t, req)

	resp := httptest.NewRecorder()
	require.NotNil(t, resp)

	ep.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, fixture, requestContent, "request body is unchanged")
}
