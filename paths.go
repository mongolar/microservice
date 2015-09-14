package service

import (
	"fmt"
)

func (s *Service) backendPath() string {
	return fmt.Sprintf("%v/backend", s.basePath())
}
func (s *Service) serverPath() string {
	return fmt.Sprintf("%v/servers/%v.%v", s.basePath(), env.Host, env.Port)
}

func (s *Service) privateServiceKeyPath() string {
	return fmt.Sprintf("%v/privatekey", s.basePath())
}

func (s *Service) basePath() string {
	return fmt.Sprintf("%v/%v.%v", VULCANPATH, s.Title, s.Version)
}
