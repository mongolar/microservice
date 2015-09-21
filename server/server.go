package server

import (
	"flag"
	"fmt"
	"github.com/mongolar/microservice/environment"
	"github.com/mongolar/microservice/herald"
	"github.com/mongolar/microservice/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var DefaultServer *Server

func init() {
	DefaultServer = New()
}

type Server struct {
	Handler http.Handler
	Service *service.Service
	Herald  herald.Herald
}

func New() *Server {
	s := &Server{
		Handler: http.DefaultServeMux,
		Service: service.DefaultService,
		Herald:  herald.DefaultHerald,
	}
	return s
}

// Set handler for the default service
func Handler(handler http.Handler) {
	DefaultServer.Handler = handler
}

// Set handler for the default service
func Service(service *service.Service) {
	DefaultServer.Service = service
}

func (s *Server) Init() {
	environment.Init()
	s.Service.Init()
	if s.Herald != nil {
		s.Herald.Init()
		s.Herald.Register(s.Service)
	}
	s.shutdown()
}

func (s *Server) Serve() {
	if !flag.Parsed() {
		flag.Parse()
	}
	s.Init()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", environment.Port()), http.HandlerFunc(s.PreServe)))
}

func (s *Server) PreServe(w http.ResponseWriter, r *http.Request) {
	if s.Service.ValidParameters(r) {
		s.Handler.ServeHTTP(w, r)
	}
}

func Serve() {
	DefaultServer.Serve()
}

func Init() {
	DefaultServer.Init()
}

func (s *Server) shutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for _ = range c {
			err := s.Herald.UnRegister(s.Service)
			if err == nil {
				os.Exit(0)
			} else {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}
	}()
}
