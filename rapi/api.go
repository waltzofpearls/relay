package rapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	config      *Config
	router      *mux.Router
	transformer *Transformer
}

func New(config *Config) *Api {
	a := new(Api)
	a.config = config
	a.router = mux.NewRouter()
	a.transformer = NewTransformer()
	return a
}

func (a *Api) Run() {
	if a.config.Listener.Tls.Enable {
		http.ListenAndServeTLS(
			a.config.Listener.Address,
			a.config.Listener.Tls.CertFile,
			a.config.Listener.Tls.KeyFile,
			a.router,
		)
	} else {
		http.ListenAndServe(
			a.config.Listener.Address,
			a.router,
		)
	}
}

func (a *Api) Route(method, path string, h http.Handler) {
	a.router.Handle(path, h).
		Methods(method)
}

func (a *Api) Router() *mux.Router {
	return a.router
}

func (a *Api) NewEndpoint(method, endpoint string) *Endpoint {
	return NewEndpoint(a, method, endpoint)
}
