package declaration

import (
	"github.com/mongolar/service"
	"github.com/mongolar/service/declaration/vulcand"
)

var declarationTypes map[string]*ServiceDeclaration
var declaredType string
var Declaration *ServiceDeclaration

func init(){
	declarationTypes = make(map[string]*ServiceDeclaration)
	AddDeclarationType("vulcand",new(vulcand.Vulcand))
	flag.StringVar(&declaredType, "declare", "", "The declaration to declare this service.")
}

func InitDeclaration(){
        if declaredType == "" {
                var err error
                declaredType, err = services.getEnvValue("SERVICE_DECLARATION")
                if err != nil {
                        log.Fatal(err)
                }
        }
	if dec, ok := declarationTypes[declaredType]; ok {
		Declaration = dec
		return
	}
	err := fmt.Errorf("Use of an unregistered declaration type: %v")
        log.Fatal(err)
}

type ServiceDeclaration interface {
	Init()
	Register(*service.Service) error
	GetService(*service.Service) error
	UnRegister(*service.Service) error
}

func AddDeclarationType (key string, sd ServiceDeclaration) {
	declarationTypes[key]sd
}

func Register(s *service.Service) error {
	return Declaration.Register(s)
}

func GetService(s *service.Service) error {
	return Declaration.GetService(s)
}

func UnRegister(*service.Service) error {
	return Declaration.UnRegister(s)
}
