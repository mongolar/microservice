package service

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

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
		Port:         port,
		env.Machines: machines,
		env.Host:     host,
	}
}

type Environment struct {
	Port     int
	Host     string
	Machines []string
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
	Public   bool
	Category string
	Title    string
	Version  string
	Handler  http.Handler
}

func (s *Service) Serve() {
	s.register()
	http.ListenAndServe(":"+env.Port, s.Handler)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for sig := range c {
			s.unregister()
		}
	}()

}

func (s *Service) register() {
	client := etcd.NewClient(env.Machines)
}

func (s *Service) unregister() {
	client := etcd.NewClient(env.Machines)
}

func (s *Service) heartbeat() {

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
