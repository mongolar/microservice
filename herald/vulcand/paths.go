package vulcand

import (
	"fmt"
)

const vulcanpath = "/vulcand/backends"

func backendPath(s *Service) string {
	return fmt.Sprintf("%v/backend", s.basePath(s))
}
func serverPath(s *Service) string {
	return fmt.Sprintf("%v/servers/%v.%v", s.basePath(s), Env.Host, Env.Port)
}

func privateKeyPath(s *Service) string {
	return fmt.Sprintf("%v/privatekey", s.basePath(s))
}

func basePath(s *Service) string {
	return fmt.Sprintf("%v/%v.%v", vulcanpath, s.Title, s.Version)
}
