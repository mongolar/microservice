package service

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service struct {
	Title   string           `json:"-"`
	Version string           `json:"-"`
	Type    string           `json:"Type"`
	backend string           `json:"-"`
	server  string           `json:"-"`
	Handler http.HandlerFunc `json:"-"`
}

func (s *Service) Serve() {
	base := fmt.Sprintf("%v/%v.%v", vulcanpath, s.Title, s.Version)
	s.backend = fmt.Sprintf("%v/backend", base)
	s.server = fmt.Sprintf("%v/servers/%v.%v", base, env.Host, env.Port)
	s.shutdown()
	s.register()
	go s.heartbeat()
	go env.refresh()
	http.ListenAndServe(fmt.Sprintf(":%v", env.Port), s.Handler)
}

func (s *Service) register() {
	client := etcd.NewClient(env.Machines)
	servicetype, _ := json.Marshal(s)
	client.Set(s.backend, string(servicetype), 0)
	serviceurl, _ := json.Marshal(env)
	client.Set(s.server, string(serviceurl), 10)
}

func (s *Service) unregister() {
	client := etcd.NewClient(env.Machines)
	client.Delete(s.server, false)
}

func (s *Service) heartbeat() {
	for _ = range time.Tick(9 * time.Second) {
		s.register()
	}
}

func (s *Service) shutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for _ = range c {
			s.unregister()
			os.Exit(0)
		}
	}()
}
