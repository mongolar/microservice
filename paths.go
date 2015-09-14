package service

import (
	"fmt"
)

func (s *Service) backendPath() string {
	return fmt.Sprintf("%v/backend", s.basePath())
}
func (s *Service) serverPath() string {
	return fmt.Sprintf("%v/servers/%v.%v", s.basePath(), Env.Host, Env.Port)
}

func (s *Service) privateServiceKeyPath() string {
	return fmt.Sprintf("%v/privatekey", s.basePath())
}

func (s *Service) basePath() string {
	return fmt.Sprintf("%v/%v.%v", vulcanpath, s.Title, s.Version)
}