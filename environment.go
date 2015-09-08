package service

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const VULCANPATH = "/vulcand/backends"

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
		host, err = getEnvValue("MONGOLAR_SERVICES_HOST")
		if err != nil {
			log.Fatal(err)
		}
	}
	if *etcdmachines == "" {
		var err error
		etcdmachines, err = getEnvValue("MONGOLAR_ETCD_MACHINES")
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
	go func() {
		for _ = range time.Tick(10 * time.Second) {
			etcdmachines, err := getEnvValue("MONGOLAR_ETCD_MACHINES")
			if err != nil || *etcdmachines != "" {
				//TODO: ERROR handling needs to be added
				fmt.Println(err)
			} else {
				env.Machines = strings.Split(*etcdmachines, "|")
			}
		}
	}()
}

func getEnvValue(name string) (*string, error) {
	value := os.Getenv(name)
	if value == "" {
		return &value, fmt.Errorf("%v is not set, %v environment value is required.", name, name)
	}
	return &value, nil
}
