package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	host = "localhost"
	port = 8000
)

type Service interface {
	Route(router *router.Router) *router.Router
}

type API interface {
	Run() error
	Close() error
}

var _ API = &api{}

type api struct {
	Router *router.Router
	server *fasthttp.Server
	logger *logrus.Logger
}

func NewAPI(log *logrus.Logger) *api {
	return &api{
		server: &fasthttp.Server{
			Name:            "admin",
			CloseOnShutdown: true,
		},
		Router: router.New(),
		logger: log,
	}
}

// Run assigns previously defined routes and starts server for listening for requests
func (a *api) Run() error {
	a.server.Handler = a.Router.Handler
	return a.server.ListenAndServe(fmt.Sprintf("%s:%d", host, port))
}

// Close shutdowns server
func (a *api) Close() error {
	return a.server.Shutdown()
}

func decode(body []byte, output interface{}) error {
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	fmt.Println(reflect.TypeOf(output))

	return decoder.Decode(output)
}
