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

var DefaultService *Service
var Frequency uint64

func init() {
	noconfig := flag.Bool("noconf", false, "Do not load Services config file.")
	flag.Uint64Var(&Frequency, "frequency", 10, "The frequency at which the service updates statuses.")
	if *noconfig {
		return
	}
	service := new(Service)
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
	service.Handler = http.DefaultServeMux
	DefaultService = service
}

func Serve() {
	DefaultService.Serve()
}

func SetHandler(handler http.Handler) {
	DefaultService.Handler = handler
}

func (s *Service) Serve() {
	if !flag.Parsed() {
		flag.Parse()
	}
	env.bootstrap()
	s.bootstrap()
	if s.Private {
		http.ListenAndServe(fmt.Sprintf(":%v", env.Port), http.HandlerFunc(s.servePrivate))
	} else {
		http.ListenAndServe(fmt.Sprintf(":%v", env.Port), s.Handler)
	}
}

func (s *Service) bootstrap() {
	s.register()
	fmt.Println(s.Private)
	if s.Private {
		s.registerPrivateServer()
	}
	s.heartbeat()
	s.watchPrivateClientKeys()
	s.shutdown()
}

func (s *Service) register() {
	client := etcd.NewClient(env.Machines())
	servicetype, _ := json.Marshal(s)
	//TODO: ERROR handling needs to be added
	client.Set(s.backendPath(), string(servicetype), 0)
	serviceurl, _ := json.Marshal(env)
	client.Set(s.serverPath(), string(serviceurl), Frequency)
}

func (s *Service) unregister() error {
	client := etcd.NewClient(env.Machines())
	_, err := client.Delete(s.serverPath(), false)
	//TODO unregister private key
	if s.Private {
		_, err = client.Delete(s.privateServiceKeyPath(), false)
	}
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
				//TODO: ERROR handling needs to be added
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}()
}
