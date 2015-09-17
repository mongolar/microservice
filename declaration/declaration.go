package declaration

import (
	"github.com/mongolar/service"
)

type ServiceDeclaration interface {
	Init()
	Register(*service.Service) error
	Get(*service.Service) error
	UnRegister(*service.Service) error
	SetPrivateKey(*service.Service)
	GetPrivateKey(*service.Service)
}
