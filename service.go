package service

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const vulcanpath = "/vulcand/backends"

var env Environment

func init() {
	var port int = flag.Int("port", 0, "The microservice port.")
	flag.Parse()
	env = Environment{
		Port:     port,
		Machines: make([]string),
	}
	if port == 0 {
		log.Fatal(errors.New("Port parameter is required."))
	}
	host = getHost()
	if host == "" {
		log.Fatal(errors.New("MONGOLAR_SERVICE_HOST is not set, service host environement value is required."))
	}
	machines, err := getEtcdMachines()
	if err != nil {
		log.Fatal(err)
	}
	env = &Environment{
		Port:     port,
		Machines: machines,
		Host:     host,
		URL:      fmt.Sprintf("http://%v:%v", host, port),
	}
}

type Environment struct {
	Port     int
	Host     string
	Machines []string
	URL      string `json:"URL"`
}

func (e *Environment) refresh() {
	machines, err := getEtcdMachines()
	if err != nil {
		// TODO pass error to service here
		fmt.Println(err)
	} else {
		etc.Machines = machines
	}
}

type Service struct {
	Title   string
	Version string
	Type    string `json:"Type"`
	Handler http.Handler
	backend string
	server  string
}

func (s *Service) Serve() {
	base := fmt.Sprintf("%v/%v.%v", vulcanpath, s.Title, s.Version)
	s.backend = fmt.Sprintf("%v/backend", base)
	s.server = fmt.Sprintf("%v/servers/%v.%v", base, env.Host, env.Port)

	s.register()
	go s.heartbeat()
	http.ListenAndServe(":"+env.Port, s.Handler)
}

func (s *Service) register() {
	client := etcd.NewClient(env.Machines)
	client.set(s.backend, json.Marshall(s.Type), 0)
	client.set(s.server, json.Marshall(env.URL), 10)
}

func (s *Service) unregister() {
	client := etcd.NewClient(env.Machines)
	client.delete(s.server, FALSE)
}

func (s *Service) heartbeat() {
	for _ := range time.Tick(9 * time.Second) {
		s.register()
	}
}

func (s *Serice) shutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for sig := range c {
			s.unregister()
		}
	}()
}

func getHost() (string, error) {
	host := os.Getenv("MONGOLAR_SERVICES_HOST")
	if host == "" {
		return host, errors.New("MONGOLAR_SERVICE_HOST is not set, service host environement value is required.")
	}
	return host, nil
}

func getEtcdMachines() ([]string, error) {
	etcd := os.Getenv("MONGOLAR_ETCD_MACHINES")
	machines := make([]string, 0)
	if etcd == "" {
		return machines, errors.New("MONGOLAR_ETCD_MACHINES is not set, etcd machines environmental value is required.")
	}
	err := json.Unmarshal([]byte(etcd), &machines)
	if err != nil {
		err = fmt.Errorf("Unable to unmarshall MONGOLAR_ETCD_MACHINES: %s", err)
	}
	return machines, err
}
