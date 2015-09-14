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

var Env Environment

func init() {
	Env = Environment{}
	flag.StringVar(&Env.Port, "port", "", "The microservice port.")
	flag.StringVar(&Env.Host, "host", "", "The microservice host.")
	flag.StringVar(&Env.machines, "etcd", "", "The etcd machines.")
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
		} else {
			e.refreshEtcdMachines()
		}
	}
	e.URL = fmt.Sprintf("http://%v:%v", e.Host, e.Port)
}

func (e *Environment) Machines() []string {
	return strings.Split(e.machines, "|")
}

func (e *Environment) refreshEtcdMachines() {
	go func() {
		for _ = range time.Tick(10 * time.Second) {
			etcdmachines, err := getEnvValue("ETCD_MACHINES")
			if err != nil || etcdmachines == "" {
				if err != nil {
					fmt.Fprint(os.Stderr, err)
				}
				if etcdmachines == "" {
					fmt.Fprint(os.Stderr, "ETCD_MACHINES not set.")
				}
			} else {
				e.machines = etcdmachines
			}
		}
	}()
}

func getEnvValue(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return value, fmt.Errorf("%v is not set, %v environment value is required", name, name)
	}
	return value, nil
}
