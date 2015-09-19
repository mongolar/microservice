// Copyright © 2015 Jason Smith <jasonrichardsmith@gmail.com>.
//
// Use of this source code is governed by the GPL-3
// license that can be found in the LICENSE file.

package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DefaultService is default server similair to the DefaultMuxServer provided to easily start new server
var DefaultService *Service

// Frequency with which to refresh service values
var Frequency uint64

func init() {
	flag.Uint64Var(&Frequency, "frequency", 10, "The frequency at which the service updates statuses.")
	service := New()
	_, err := os.Stat("Service.yaml")
	if err == nil {
		v := viper.New()
		v.SetConfigName("Service")
		err := v.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = v.Marshal(service)
		if err != nil {
			log.Fatal(err)
		}
	}
	DefaultService = service
}

// Service definition
type Service struct {
	Title      string       `json:"Title"`
	Version    string       `json:"Version"`
	Type       string       `json:"Type"`
	Private    bool         `json:"Private"`
	Requires   []Service    `json:"Requires,omitempty"`
	Parameters Parameters   `json:"Parameters"`
	Method     string       `json:"Method"`
	Handler    http.Handler `json:"-"`
}

// Get a new Service and set the default Handler to the DefaultServerMux
func New() *Service {
	service := new(Service)
	service.Handler = http.DefaultServeMux
	return service
}

// Set handler for the deault service
func Handler(handler http.Handler) {
	DefaultService.Handler = handler
}

// Return service description based on title and version
func GetService(title string, version string) (*Service, error) {
	service := &Service{Title: title, Version: version, foreign: true}
	err := service.GetService()
	return service, err
}

// GetService based on instantiated service, requires Title and Version to be set
func (s *Service) GetService() error {
	if s.Title == "" || Version == "" {
		return errors.New("Title and Version is required to retrieve a service")
	}
	herald.GetService(s)

}

func Serve() {
	DefaultService.Serve()
}

func (s *Service) Serve() {
	if !flag.Parsed() {
		flag.Parse()
	}
	Env.bootstrap()
	s.bootstrap()
	if s.Private {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", Env.Port), http.HandlerFunc(s.servePrivate)))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", Env.Port), s.Handler))
	}
}

func (s *Service) bootstrap() {
	s.register()
	if s.Private {
		s.registerPrivateService()
	}
	s.checkRequired()
	s.shutdown()
}

func (s *Service) register() {
	herald.Register(s)
}

func (s *Service) unregister() error {
	herald.UnRegister(s)
}

func (s *Service) checkRequired() {
	client := etcd.NewClient(Env.Machines())
	defer client.Close()
	for _, rs := range s.Requires {
		err := rs.GetService()
		if err != nil {
			fmt.Fprintf(os.Stderr, err)
		}
	}
}

func (s *Service) shutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for _ = range c {
			err := s.unregister()
			if err == nil {
				os.Exit(0)
			} else {
				fmt.Fprintf(os.Stderr, err)
				os.Exit(1)
			}
		}
	}()
}
