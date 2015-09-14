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
	env = Environment{}
	flag.StringVar(&env.Port, "port", "", "The microservice port.")
	flag.StringVar(&env.Host, "host", "", "The microservice host.")
	flag.StringVar(&env.machines, "etcd", "", "The etcd machines.")
}

type Environment struct {
	Port     string `json:"-"`
	Host     string `json:"-"`
	URL      string `json:"URL"`
	machines string `json:"-"`
}

func (e *Environment) bootstrap() {

	if e.Port == "" {
		log.Fatal(errors.New("Port parameter is a required flag."))
	}
	if e.Host == "" {
		var err error
		e.Host, err = getEnvValue("MICRO_SERVICES_HOST")
		if err != nil {
			log.Fatal(err)
		}
	}
	if e.machines == "" {
		var err error
		e.machines, err = getEnvValue("ETCD_MACHINES")
		if err != nil {
			log.Fatal(err)
		}
	}
	e.URL = fmt.Sprintf("http://%v:%v", e.Host, e.Port)
	env.refresh()
}

func (e *Environment) Machines() []string {
	return strings.Split(env.machines, "|")
}

func (e *Environment) refresh() {

	go func() {
		for _ = range time.Tick(10 * time.Second) {
			etcdmachines, err := getEnvValue("ETCD_MACHINES")
			if err != nil || etcdmachines != "" {
				//TODO: ERROR handling needs to be added
				fmt.Println(err)
			} else {
				env.machines = etcdmachines
			}
		}
	}()
}

func getEnvValue(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return value, fmt.Errorf("%v is not set, %v environment value is required.", name, name)
	}
	return value, nil
}
