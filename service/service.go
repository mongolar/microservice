// Copyright © 2015 Jason Smith <jasonrichardsmith@gmail.com>.
//
// Use of this source code is governed by the GPL-3
// license that can be found in the LICENSE file.

package service

import (
	"flag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

var serviceFile string
var DefaultService *Service

func init() {
	DefaultService = New()
	flag.StringVar(&serviceFile, "service", "Service.yaml", "Full path to service file.")
}

// Service definition
type Service struct {
	Title      string     `json:"Title"`
	Version    string     `json:"Version"`
	Type       string     `json:"Type"`
	Private    bool       `json:"Private"`
	Requires   []Service  `json:"Requires,omitempty"`
	Parameters Parameters `json:"Parameters"`
	Method     string     `json:"Method"`
}

// Get a new Service and set the default Handler to the DefaultServerMux
func New() *Service {
	service := new(Service)
	return service
}

func (s *Service) Init() {
	s.Marshal()
}

func (s *Service) Marshal() {
	s.MarshalF(serviceFile)
	return
}

func (s *Service) ValidParameters(r *http.Request) bool {
	return true
}

func (s *Service) MarshalF(file string) {
	_, err := os.Stat(file)
	if err == nil {
		v := viper.New()
		v.SetConfigName(file)
		err := v.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = v.Marshal(s)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Fatal(err)
}
