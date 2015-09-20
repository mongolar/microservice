package vulcand

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/mongolar/microservice/environment"
	"github.com/mongolar/microservice/service"
	"log"
	"os"
	"strings"
	"time"
)

var etcdmachines string

func init() {
	flag.StringVar(&etcdmachines, "etcdv", "", "The etcd machines for Vulcand.")
}

type Vulcand struct{}

func (v Vulcand) Init() {
	if etcdmachines == "" {
		var err error
		etcdmachines, err = environment.GetEnvValue("ETCD_MACHINES")
		if err != nil {
			log.Fatal(err)
		} else {
			refreshEtcdMachines()
		}
	}

}
func (v Vulcand) Register(s *service.Service) error {
	client := etcd.NewClient(Machines())
	servicetype, err := json.Marshal(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	//TODO: ERROR handling needs to be added
	_, err = client.Set(backendPath(s), string(servicetype), 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	heartbeat(s)
	return err
}

func heartbeat(s *service.Service) {
	go func() {
		for _ = range time.Tick(time.Duration(environment.Frequency()-1) * time.Second) {
			setServer(s)
		}
	}()
}

func setServer(s *service.Service) {
	serviceurl, err := json.Marshal(environment.URL())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	client := etcd.NewClient(Machines())
	_, err = client.Set(serverPath(s), string(serviceurl), environment.Frequency())
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func (v Vulcand) GetService(s *service.Service) error {
	client := etcd.NewClient(Machines())
	raw, err := client.RawGet(backendPath(s), false, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw.Body, s)
	return err
}

func (v Vulcand) UnRegister(s *service.Service) error {
	client := etcd.NewClient(Machines())
	_, err := client.Delete(serverPath(s), false)
	return err
}

func Machines() []string {
	return strings.Split(etcdmachines, "|")
}

func refreshEtcdMachines() {
	go func() {
		for _ = range time.Tick(time.Duration(environment.Frequency()) * time.Second) {
			machines, err := environment.GetEnvValue("ETCD_MACHINES")
			if err != nil || machines == "" {
				if err != nil {
					fmt.Fprint(os.Stderr, err)
				}
				if machines == "" {
					fmt.Fprintf(os.Stderr, "ETCD_MACHINES not set.")
				}
			} else {
				etcdmachines = machines
			}
		}
	}()
}
