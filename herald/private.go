package herald

import (
	"crypto/rand"
	"encoding/base64"
	"github/mongolar/service"
)

var private

func init() {
	flag.Bool(&private, "private", "", "Is this service private.")
}

type PrivateHerald interface {
	InitPrivate()
	SetPrivateKey(*service.Service) error
	GetPrivateKey(*service.Service) error
	ValidatePrivateKey(*service.Service, *http.Request) (bool, error)
}

func SetPrivateKey(s *service.Service){
	return DefaultHerald.SetPrivateKey(s)
}

func GetPrivateKey(s *service.Service){
	return DefaultHerald.GetPrivateKey(s)
}

func ValidatePrivateKey(s *service.Service, r *http.Request) (bool, error){
	return DefaultHerald.ValidatePrivateKey(s, r)
}

func IsPrivate(d *Herald) bool {
	defer func() {
		if r := recover(); r != nil {
			return false
		}
	}()
	var _ PrivateHerald = d
	return true
}

func NewPrivateKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		//TODO output err
		return newPrivateKey()
	}
	return base64.URLEncoding.EncodeToString(key)
}
