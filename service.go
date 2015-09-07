package service

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const vulcanpath = "/vulcand/backends"

var env Environment

func init() {
	port := flag.String("port", "", "The microservice port.")
	host := flag.String("host", "", "The microservice host.")
	etcdmachines := flag.String("etcd", "", "The etcd machines.")
	flag.Parse()
	if *port == "" {
		log.Fatal(errors.New("Port parameter is required."))
	}
	if *host == "" {
		var err error
		host, err = getEnvHost()
		if err != nil {
			log.Fatal(err)
		}
	}
	if *etcdmachines == "" {
		var err error
		etcdmachines, err = getEnvEtcdMachines()
		if err != nil {
			log.Fatal(err)
		}
	}
	env = Environment{
		Port:     *port,
		Machines: strings.Split(*etcdmachines, "|"),
		Host:     *host,
		URL:      fmt.Sprintf("http://%v:%v", *host, *port),
	}
}

type Environment struct {
	Port     string   `json:"-"`
	Host     string   `json:"-"`
	Machines []string `json:"-"`
	URL      string   `json:"URL"`
}

func (e *Environment) refresh() {
	etcdmachines, err := getEnvEtcdMachines()
	if err != nil {
		// TODO pass error to service here
		fmt.Println(err)
	} else {
		env.Machines = strings.Split(*etcdmachines, "|")
	}
}

type Service struct {
	Title   string           `json:"-"`
	Version string           `json:"-"`
	Type    string           `json:"Type"`
	backend string           `json:"-"`
	server  string           `json:"-"`
	Handler http.HandlerFunc `json:"-"`
}

func (s *Service) Serve() {
	base := fmt.Sprintf("%v/%v.%v", vulcanpath, s.Title, s.Version)
	s.backend = fmt.Sprintf("%v/backend", base)
	s.server = fmt.Sprintf("%v/servers/%v.%v", base, env.Host, env.Port)
	s.shutdown()
	s.register()
	go s.heartbeat()
	http.ListenAndServe(fmt.Sprintf(":%v", env.Port), s.Handler)
}

func (s *Service) register() {
	client := etcd.NewClient(env.Machines)
	servicetype, _ := json.Marshal(s)
	client.Set(s.backend, string(servicetype), 0)
	serviceurl, _ := json.Marshal(env)
	client.Set(s.server, string(serviceurl), 10)
}

func (s *Service) unregister() {
	client := etcd.NewClient(env.Machines)
	client.Delete(s.server, false)
}

func (s *Service) heartbeat() {
	for _ = range time.Tick(9 * time.Second) {
		s.register()
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
			s.unregister()
			os.Exit(0)
		}
	}()
}

func getEnvHost() (*string, error) {
	host := os.Getenv("MONGOLAR_SERVICES_HOST")
	if host == "" {
		return &host, errors.New("MONGOLAR_SERVICES_HOST is not set, service host environement value is required.")
	}
	return &host, nil
}

func getEnvEtcdMachines() (*string, error) {
	etcdmachines := os.Getenv("MONGOLAR_ETCD_MACHINES")
	if etcdmachines == "" {
		return &etcdmachines, errors.New("MONGOLAR_ETCD_MACHINES is not set, etcd machines environmental value is required.")
	}
	return &etcdmachines, nil
}
