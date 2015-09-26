package environment

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

var DefaultEnvironment Environment

func init() {
	DefaultEnvironment = Environment{}
	flag.StringVar(&DefaultEnvironment.Port, "port", "", "The microservice port.")
	flag.StringVar(&DefaultEnvironment.Host, "host", "", "The microservice host.")
	flag.Uint64Var(&DefaultEnvironment.Frequency, "frequency", 10, "The frequency at which the service updates statuses.")
	flag.StringVar(&DefaultEnvironment.IntServiceURL, "int", "", "The internal service to service url.")
}

type Environment struct {
	Port          string `json:"-"`
	Host          string `json:"-"`
	Frequency     uint64 `json:"-"`
	URL           string `json:"URL"`
	IntServiceURL string `json:"-"`
}

func (e *Environment) Init() {
	if e.Port == "" {
		log.Fatal(errors.New("Port parameter is a required flag."))
	}
	if e.Host == "" {
		var err error
		e.Host, err = GetEnvValue("MICRO_SERVICES_HOST")
		if err != nil {
			log.Fatal(err)
		}
	}
	if e.IntServiceURL == "" {
		var err error
		e.Host, err = GetEnvValue("MICRO_SERVICES_INT_URL")
		if err != nil {
			log.Fatal(err)
		}
	}
	e.URL = fmt.Sprintf("http://%v:%v", e.Host, e.Port)
}

func Init() {
	DefaultEnvironment.Init()
}

func Port() string {
	return DefaultEnvironment.Port
}

func Host() string {
	return DefaultEnvironment.Host
}

func Frequency() uint64 {
	return DefaultEnvironment.Frequency
}

func URL() string {
	return DefaultEnvironment.URL
}

func IntServiceURL() string {
	return DefaultEnvironment.IntServiceURL
}

func GetEnvValue(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return value, fmt.Errorf("%v is not set, %v environment value is required", name, name)
	}
	return value, nil
}
