package declaration

import (
	"github/mongolar/service"
)

type PrivateDeclaration interface {
	SetPrivateKey(*service.Service)
	GetPrivateKey(*service.Service)
}

func validPrivate(d *Declaration) bool {
	defer func() {
		if r := recover(); r != nil {
			return false
		}
	}
	var _ PrivateDeclaration = d
	return true
}
