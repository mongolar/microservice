package vulcand

import (
	"fmt"
	"github.com/mongolar/microservice/environment"
	"github.com/mongolar/microservice/service"
)

const vulcanpath = "/vulcand/backends"

func backendPath(s *service.Service) string {
	return fmt.Sprintf("%v/backend", basePath(s))
}
func serverPath(s *service.Service) string {
	return fmt.Sprintf("%v/servers/%v.%v", basePath(s), environment.Host(), environment.Port())
}

func basePath(s *service.Service) string {
	return fmt.Sprintf("%v/%v.%v", vulcanpath, s.Title, s.Version)
}
