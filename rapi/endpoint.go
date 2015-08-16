package rapi

import (
	"io"
	"log"
	"net/http"
)

type Endpoint struct {
	api          *Api
	config       *ConfigItem
	extPath      string
	intPath      string
	method       string
	reqExtStruct interface{}
	reqIntStruct interface{}
	resExtStruct interface{}
	resIntStruct interface{}
}

func NewEndpoint(a *Api, method, path string) *Endpoint {
	ep := &Endpoint{
		api:          a,
		config:       a.config,
		extPath:      path,
		intPath:      path,
		method:       method,
		reqExtStruct: nil,
		reqIntStruct: nil,
		resExtStruct: nil,
		resIntStruct: nil,
	}

	a.router.
		Handle(ep.config.ExtPathPrefix+path, ep).
		Methods(method)

	return ep
}

func (ep *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tr := http.DefaultTransport

	r.URL.Host = ep.config.Downstream
	r.URL.Scheme = "http"
	r.URL.Path = "/api" + ep.intPath

	res, resErr := tr.RoundTrip(r)
	if resErr == nil {
		defer res.Body.Close()
	}

	w.WriteHeader(res.StatusCode)
	_, ioErr := io.Copy(w, res.Body)

	if ioErr != nil {
		log.Printf("Error writting response: %s", ioErr)
	}

	// fmt.Fprint(w, "Protected!!!!\n"+r.URL.Query().Get("querystring1"))
	// w.Write([]byte("Gorilla!\n"))
}

func (ep *Endpoint) InternalPath(path string) *Endpoint {
	ep.intPath = path
	return ep
}

func (ep *Endpoint) TransformRequest(external, internal interface{}) *Endpoint {
	ep.reqExtStruct = external
	ep.reqIntStruct = internal
	return ep
}

func (ep *Endpoint) TransformResponse(external, internal interface{}) *Endpoint {
	ep.resExtStruct = external
	ep.resIntStruct = internal
	return ep
}

// To be done
func (ep *Endpoint) TransformRequestCb(cb TransformCb) *Endpoint {
	err := cb()
	if err != nil {
		log.Print("Something went wrong")
	}
	return ep
}

// To be done
func (ep *Endpoint) TransformResponseCb(cb TransformCb) *Endpoint {
	err := cb()
	if err != nil {
		log.Print("Something went wrong")
	}
	return ep
}
