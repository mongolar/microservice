package declaration

import (
	"github.com/mongolar/service"
	"github.com/mongolar/service/declaration/vulcand"
)

var declarationTypes map[string]*ServiceDeclaration
var declaredType string
var DeclarationType

func init(){
	declarationTypes = make(map[string]*ServiceDeclaration)
	AddDeclarationType("vulcand",new(vulcand.Vulcand))
	flag.StringVar(&declaredType, "declare", "", "The declaration to declare this service.")
}


type ServiceDeclaration interface {
	Init()
	Register(*service.Service) error
	Get(*service.Service) error
	UnRegister(*service.Service) error
}

func AddDeclarationType (key string, sd ServiceDeclaration) {
	declarationTypes[key]sd
}
