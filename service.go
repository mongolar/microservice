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

var DefaultService *Service
var Frequency uint64

// Service defines the service that will be declared to
// Vulcand.
type Service struct {
	Title                string       `json:"Title"`
	Version              string       `json:"Version"`
	Type                 string       `json:"Type"`
	Private              bool         `json:"Private"`
	Requires             []Service    `json:"Requires,omitempty"`
	Parameters           []string     `json:"Parameters"`
	Method               string       `json:"Method"`
	Handler              http.Handler `json:"-"`
	privateClientKeys    map[string]string
	privateServiceKeyOld string
	privateServiceKey    string
}

func New() *Service {
	service := new(Service)
	service.Handler = http.DefaultServeMux
	return service
}

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

func Serve() {
	DefaultService.Serve()
}

func Handler(handler http.Handler) {
	DefaultService.Handler = handler
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
		s.registerPrivateServer()
	}
	s.heartbeat()
	s.watchPrivateClientKeys()
	s.shutdown()
}

func (s *Service) register() {
	client := etcd.NewClient(Env.Machines())
	servicetype, err := json.Marshal(s)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	//TODO: ERROR handling needs to be added
	_, err = client.Set(s.backendPath(), string(servicetype), 0)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	serviceurl, err := json.Marshal(Env)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	_, err = client.Set(s.serverPath(), string(serviceurl), Frequency)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	client.Close()
}

func (s *Service) unregister() error {
	client := etcd.NewClient(Env.Machines())
	_, err := client.Delete(s.serverPath(), false)
	//TODO unregister private key
	if s.Private {
		_, err = client.Delete(s.privateServiceKeyPath(), false)
	}
	client.Close()
	return err
}

func (s *Service) heartbeat() {
	go func() {
		for _ = range time.Tick(time.Duration(Frequency-1) * time.Second) {
			s.register()
		}
	}()
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
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}
	}()
}
