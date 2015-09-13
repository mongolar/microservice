package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service struct {
	Title               string       `json:"Title"`
	Version             string       `json:"Version"`
	Type                string       `json:"Type"`
	Private             bool         `json:"Private"`
	Requires            []Service    `json:"Requires,omitempty"`
	Parameters          []string     `json:"Parameters"`
	Method              string       `json:"Method"`
	Handler             http.Handler `json:"-"`
	privateClientKeys   map[string]string
	privateServerKeyOld string
	privateServerKey    string
}

func New() *Service {
	service := new(Service)
	v := viper.New()
	v.SetConfigName("Service")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = v.Marshal(service)
	if err != nil {
		log.Fatal(err)
	}
	return service
}

func (s *Service) Serve() {
	s.bootstrap()
	if s.Private {
		http.ListenAndServe(fmt.Sprintf(":%v", env.Port), http.HandlerFunc(s.servePrivate))
	} else {
		http.ListenAndServe(fmt.Sprintf(":%v", env.Port), s.Handler)
	}
}

func (s *Service) servePrivate(w http.ResponseWriter, r *http.Request) {
	if s.validatePrivateServer(r) {
		s.Handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}

func (s *Service) bootstrap() {
	s.shutdown()
	s.registerPrivateServer()
	s.heartbeat()
	s.watchPrivateClientKeys()
}

func (s *Service) BackendPath() string {
	return fmt.Sprintf("%v/backend", s.basePath())
}
func (s *Service) ServerPath() string {
	return fmt.Sprintf("%v/servers/%v.%v", s.basePath(), env.Host, env.Port)
}

func (s *Service) PrivateServerKeyPath() string {
	return fmt.Sprintf("%v/privatekey", s.basePath())
}

func (s *Service) basePath() string {
	return fmt.Sprintf("%v/%v.%v", VULCANPATH, s.Title, s.Version)
}

func (s *Service) validatePrivateServer(r *http.Request) bool {
	key := r.Header.Get("PrivateServerKey")
	if key != s.privateServerKey && key != s.privateServerKeyOld {
		return false
	}
	return true
}

func (s *Service) updatePrivateServerKeys(r *etcd.Response) {
	if r.Action == "expire" || r.Action == "delete" {
		client := etcd.NewClient(env.Machines)
		s.lead(client)
		return
	}
	if r.PrevNode != nil {
		s.privateServerKeyOld = r.PrevNode.Value
	}
	s.privateServerKey = r.Node.Value
}

func (s *Service) register() {
	client := etcd.NewClient(env.Machines)
	servicetype, _ := json.Marshal(s)
	//TODO: ERROR handling needs to be added
	client.Set(s.BackendPath(), string(servicetype), 0)
	serviceurl, _ := json.Marshal(env)
	client.Set(s.ServerPath(), string(serviceurl), 10)
}
func (s *Service) registerPrivateServer() {
	fmt.Println(s.PrivateServerKeyPath())
	client := etcd.NewClient(env.Machines)
	_, err := client.Get(s.PrivateServerKeyPath(), false, false)
	s.follow(client)
	if err != nil {
		//TODO check for key not set error
		s.lead(client)
	}
}

func (s *Service) unregister() error {
	client := etcd.NewClient(env.Machines)
	_, err := client.Delete(s.ServerPath(), false)
	//TODO unregister private key
	if s.Private {
		_, err = client.Delete(s.ServerPath(), false)
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

func (s *Service) follow(client *etcd.Client) {
	watchPrivateKey(s.PrivateServerKeyPath(), s.updatePrivateServerKeys)
}

func (s *Service) lead(client *etcd.Client) {
	// Attempt to set leadership for private key management
	_, err := client.Create(s.PrivateServerKeyPath(), newPrivateServerKey(), 10)
	// If err from creation, management has already been established by another node
	if err != nil { //TODO: check if err is key already set
		// Already follower exit
		return
	}
	// This is leadership so start go routine
	go func() {
		for _ = range time.Tick(9 * time.Second) {
			_, err := client.Set(s.PrivateServerKeyPath(), newPrivateServerKey(), 10)
			// If set is successful update both keys
			if err == nil {
				//TODO: handle error
				fmt.Println(err)
			}
		}
	}()
}

func (s *Service) watchPrivateClientKeys() {
	for _, r := range s.Requires {
		if r.Private {
			key := fmt.Sprintf("%v.%v", r.Title, r.Version)
			watchPrivateKey(key, s.updatePrivateClientKey)
		}
	}
}

func (s *Service) updatePrivateClientKey(r *etcd.Response) {
	s.privateClientKeys[r.Node.Key] = r.Node.Value
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

func watchPrivateKey(key string, set func(*etcd.Response)) {
	client := etcd.NewClient(env.Machines)
	wc := make(chan *etcd.Response)
	go client.Watch(key, 0, false, wc, nil)
	go func() {
		for r := range wc {
			set(r)
		}
	}()
}

func newPrivateServerKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(key)
}
