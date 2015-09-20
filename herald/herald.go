package herald

import (
	"flag"
	"fmt"
	"github.com/mongolar/microservice/environment"
	"github.com/mongolar/microservice/herald/vulcand"
	"github.com/mongolar/microservice/service"
	"log"
)

var heraldTypes map[string]Herald
var heraldType string
var DefaultHerald Herald

func init() {
	heraldTypes = make(map[string]Herald)
	AddHeraldType("vulcand", vulcand.Vulcand{})
	flag.StringVar(&heraldType, "declare", "", "The herald to declare this service.")
}

func Init() {
	if heraldType == "" {
		var err error
		heraldType, err = environment.GetEnvValue("SERVICE_DECLARATION")
		if err != nil {
			log.Fatal(err)
		}
	}
	if dec, ok := heraldTypes[heraldType]; ok {
		DefaultHerald = dec
		DefaultHerald.Init()
		return
	}
	err := fmt.Errorf("Use of an unregistered herald type: %v")
	log.Fatal(err)
}

type Herald interface {
	Init()
	Register(*service.Service) error
	UnRegister(*service.Service) error
	GetService(*service.Service) error
}

func AddHeraldType(key string, h Herald) {
	heraldTypes[key] = h
}

func Set(h Herald) {
	DefaultHerald = h
}

func Register(s *service.Service) error {
	return DefaultHerald.Register(s)
}

func GetServiceS(title string, version string) (*service.Service, error) {
	service := &service.Service{Title: title, Version: version}
	err := GetService(service)
	return service, err
}

func GetService(s *service.Service) error {
	return DefaultHerald.GetService(s)
}

func UnRegister(s *service.Service) error {
	return DefaultHerald.UnRegister(s)
}
