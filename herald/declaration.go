package herald

import (
	"github.com/mongolar/service"
	"github.com/mongolar/service/herald/vulcand"
)

var heraldTypes map[string]*Herald
var heraldType string
var DefaultHerald *Herald

func init(){
	heraldTypes = make(map[string]*Herald)
	AddHeraldType("vulcand",new(vulcand.Vulcand))
	flag.StringVar(&heraldType, "declare", "", "The herald to declare this service.")
}

func InitHerald(){
        if heraldType == "" {
                var err error
                heraldType, err = services.getEnvValue("SERVICE_DECLARATION")
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
	if private {
		if !IsPrivate(DefaultHerald) {
			err := fmt.Errorf("Private flag set, %v does not support private services.", heraldType)
                        log.Fatal(err)
		}
	}
}

type Herald interface {
	Init()
	Register(*service.Service) error
	UnRegister(*service.Service) error
	GetService(*service.Service) error
}

func AddHeraldType (key string, h Herald) {
	heraldTypes[key]h
}

func SetHerald(h *Herald){
	DefaultHerald = h
}

func Register(s *service.Service) error {
	return DefaultHerald.Register(s)
}

func GetService(s *service.Service) error {
	return DefaultHerald.GetService(s)
}

func UnRegister(*service.Service) error {
	return DefaultHerald.UnRegister(s)
}
