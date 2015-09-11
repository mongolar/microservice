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
	Title               string            `json:"Title"`
	Version             string            `json:"Version"`
	Type                string            `json:"Type"`
	Private             bool              `json:"Private"`
	Requires            []Service         `json:"Requires,omitempty"`
	Parameters          []string          `json:"Parameters"`
	Method              string            `json:"Method"`
	Handler             http.HandlerFunc  `json:"-"`
	privateClientKeys   map[string]string `json:"-"`
	privateServerKeyOld string
	privateServerKey    string
}

func (s *Service) Serve() {
	s.bootstrap()
	if s.Private {
		http.ListenAndServe(fmt.Sprintf(":%v", env.Port), s.ServePrivate)
	} else {
		http.ListenAndServe(fmt.Sprintf(":%v", env.Port), s.Handler)
	}
}

func (s *Service) ServePrivate(w http.ResponseWriter, r *http.Request) {
	if s.validatePrivate(r) {
		s.Handler(w, r)
	} else {
		http.Error(w.Writer, "Forbidden", 403)
	}
}

func (s *Service) bootstrap() {
	s.shutdown()
	s.register()
	s.heartbeat()
	s.clientKeys()
}

func (s *Service) BackendPath() {
	return fmt.Sprintf("%v/backend", s.basePath())
}
func (s *Service) ServerPath() {
	return fmt.Sprintf("%v/servers/%v.%v", s.BackendPath(), env.Host, env.Port)
}

func (s *Service) PrivateServerKeyPath() string {
	return fmt.Sprintf("%v/privatekey", s.Service.BackendPath())
}

func (s *Service) basePath() {
	return fmt.Sprintf("%v/%v.%v", VULCANPATH, s.Title, s.Version)
}

func (s *Service) follow(client etcd.Client) {
	watchPrivateKey(key, s.updatePrivateServerKeys)
}

func (s *Service) lead(client etcd.Client) {
	// Attempt to set leadership for private key management
	resp, err := client.Create(s.PrivateServerKeyPath(), s.newPrivateServerKey(), 10)
	// If err from creation, management has already been established by another node
	if err != nil { //TODO: check if err is key already set
		// Take follower role and exit
		s.follow(client)
		return
	}
	// This is leadership so start go routine
	go func() {
		for _ = range time.Tick(9 * time.Second) {
			resp, err := client.Set(s.PrivateServerKeyPath(), s.newPrivateServerKey(), 10)
			// If set is successful update both keys
			if err == nil {
				s.updatePrivateServerKeys(resp)
			}
		}
	}()
}

func (s *Service) validatePrivateServer(r *http.Request) bool {
	key = r.Header.Get("PrivateServerKey")
	if key != s.privateKey && key != s.privateKeyOld {
		return false
	}
	return true
}

func (s *Service) updatePrivateServerKeys(r etcd.Response) {
	s.privateKeyOld = r.PrevNode.Value
	s.privateKey = r.Node.Value
}

func (s *Service) register() {
	client := etcd.NewClient(env.Machines)
	servicetype, _ := json.Marshal(s)
	//TODO: ERROR handling needs to be added
	client.Set(s.BackendPath(), string(servicetype), 0)
	serviceurl, _ := json.Marshal(env)
	client.Set(s.ServerPath(), string(serviceurl), 10)
	if s.Private {
		s.registerPrivateServer()
	}
}
func (s *Service) registerPrivateServer() {
	client := etcd.NewClient(env.Machines)
	pk, err := client.Get(s.PrivateServiceKeyPath(), false, false)
	if err != nil {
		//TODO check for key not set error
		s.lead(client)
	} else {
		s.follow(client)
	}
}

func (s *Service) unregister() error {
	client := etcd.NewClient(env.Machines)
	_, err := client.Delete(s.ServerPath, false)
	//TODO unregister private key
	if s.Private {
		_, err := client.Delete(s.ServerPath, false)
	}
	return err
}

func (s *Service) heartbeat() {
	go func() {
		for _ = range time.Tick(9 * time.Second) {
			s.register()
		}
	}()
}

func (s *Service) privateClientKeys() {
	for r := range s.Required {
		if r.Private {
			key := fmt.Sprintf("%v.%v", r.Title, r.Version)
			watchPrivateKey(key, s.updatePrivateKey)
		}
	}
}

func (s *Service) updatePrivateClientKey(r etcd.Response) {
	s.privatekeys[r.Node.Key] = r.Node.Value
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
			err := s.unregister()
			if err == nil {
				os.Exit(0)
			} else {
				//TODO: ERROR handling needs to be added
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}()
}

func watchPrivateKey(key string, set func(etcd.Response)) {
	client := etcd.NewClient(env.Machines)
	wc := make(chan etcd.Response)
	location := fmt.Sprintf("%v/%v", ETCDSERVICEKEYS, key)
	go client.Watch(location, 0, false, wc, nil)
	go func() {
		for r := range wc {
			set(r)
		}
	}()
}
