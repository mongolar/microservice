package herald

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/mongolar/service"
)

var private

func init() {
	flag.Bool(&private, "private", "", "Is this service private.")
}

type PrivateHerald interface {
	InitPrivate(*service.Service)
	ValidatePrivate(*service.Service, *http.Request) (bool, error)
	SetPrivateRequest(*service.Service, *http.Request) error
}

func ValidatePrivate(s *service.Service, r *http.Request) (bool, error){
	return DefaultHerald.ValidatePrivate(s, r)
}

func SetPrivate(s *service.Service, r *http.Request) error{
	return DefaultHerald.SetPrivate(s, r)
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
