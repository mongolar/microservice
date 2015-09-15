package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"net/http"
	"os"
	"time"
)

func (s *Service) servePrivate(w http.ResponseWriter, r *http.Request) {
	if s.validatePrivate(r) {
		s.Handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}

func (s *Service) validatePrivate(r *http.Request) bool {
	key := r.Header.Get("PrivateServiceKey")
	if key != s.privateKey && key != s.privateKeyOld {
		return false
	}
	return true
}

func (s *Service) registerPrivateService() {
	s.privateKey = newPrivateKey()
	s.privateKeyOld = newPrivateKey()
	s.follow()
	s.lead()
}

func (s *Service) follow() {
	watchPrivateKey(s.privateKeyPath(), s.updatePrivateKeys)
}

func (s *Service) lead() {
	client := etcd.NewClient(Env.Machines())
	_, err := client.Create(s.privateKeyPath(), newPrivateKey(), Frequency)
	defer client.Close()
	if err != nil { //TODO: check if err is key already set
		return
	}
	go func() {
		for _ = range time.Tick(time.Duration(Frequency-1) * time.Second) {
			_, err := client.Set(s.privateKeyPath(), newPrivateKey(), Frequency)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}
	}()
}

func (s *Service) updatePrivateKeys(r *etcd.Response) {
	if !s.foreign {
		if r.Action == "expire" || r.Action == "delete" {
			s.lead()
			return
		}
	}
	if r.PrevNode != nil {
		s.privateKeyOld = r.PrevNode.Value
	}
	s.privateKey = r.Node.Value
}

func watchPrivateKey(key string, set func(*etcd.Response)) {
	client := etcd.NewClient(Env.Machines())
	defer client.Close()
	wc := make(chan *etcd.Response)
	go client.Watch(key, 0, false, wc, nil)
	go func() {
		for r := range wc {
			set(r)
		}
	}()
}

func newPrivateKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		//TODO output err
		return newPrivateKey()
	}
	return base64.URLEncoding.EncodeToString(key)
}
