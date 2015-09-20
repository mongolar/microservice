package vulcand

import (
	"github.com/mongolar/service"
)

var etcdmachines string

func init() {
	flag.StringVar(&etcdmachines, "etcdv", "", "The etcd machines for Vulcand.")
}

type Vulcand struct{}

func (v *Vulcand) Init() {
	if etcdmachines == "" {
		var err error
		etcdmachines, err = service.getEnvValue("ETCD_MACHINES")
		if err != nil {
			log.Fatal(err)
		} else {
			refreshEtcdMachines()
		}
	}

}
func (v *Vulcand) Register(s *service.Service) error {
	client := etcd.NewClient(Machines())
	defer client.Close()
	servicetype, err := json.Marshal(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, err)
	}
	//TODO: ERROR handling needs to be added
	_, err = client.Set(backendPath(s), string(servicetype), 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, err)
	}
	serviceurl, err := json.Marshal(service.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, err)
	}
	heartbeat(s)
}

func heartbeat(s *service.Service) {
	go func() {
		for _ = range time.Tick(time.Duration(service.Frequency-1) * time.Second) {
			v.Register(s)
		}
	}()
}

func setServer(s *service.Service) {
	client := etcd.NewClient(Machines())
	defer client.Close()
	_, err = client.Set(serverPath(s), string(serviceurl), Frequency)
	if err != nil {
		fmt.Fprintf(os.Stderr, err)
	}
}

func (v *Vulcand) GetService(s *service.Service) error {
	client := etcd.NewClient(Machines())
	defer client.Close()
	raw, err := client.RawGet(backendPath(s), false, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw.Body, s)
	return err
}

func (v *Vulcand) UnRegister(s *service.Service) error {
	client := etcd.NewClient(Machines())
	_, err := client.Delete(serverPath(s), false)
	client.Close()
	return err
}

func (v *Vulcand) InitPrivate(s *service.Service) {
	getPrivateKey(s)
}

func (v *Vulcand) ValidatePrivate(s *service.Service, r *http.Request) (bool, error) {
	rkey := r.Header.Get("PrivateServiceKey")
	pkey, err := getPrivateKey(s)
	if err != nil {
		return false, nil
	}
	if rkey == pkey {
		return false, nil
	}
	return true

}

func getPrivateKey(s *service.Service) (string, error) {
	client := etcd.NewClient(Machines())
	defer client.Close()
	val, err := client.Get(privateKeyPath(s))
	if err != nil {
		code := err.Error()
		if "100" == code[0:3] {
			resp, err := client.Create(privateKeyPath(), newPrivateKey(), service.Frequency)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			} else {
				go lead()
			}
			return resp.Node.Value, err
		} else {
			return err
		}
		val, err := client.Create(privateKeyPath(s))
		fmt.Fprint(os.Stderr, err)
	}
	return resp.Node.Value, err
}

func lead() {
	client := etcd.NewClient(Machines())
	defer client.Close()
	for _ = range time.Tick(time.Duration(service.Frequency-1) * time.Second) {
		_, err := client.Set(privateKeyPath(), newPrivateKey(), service.Frequency)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
}

func Machines() []string {
	return strings.Split(etcdmachines, "|")
}

func refreshEtcdMachines() {
	go func() {
		for _ = range time.Tick(service.Frequency * time.Second) {
			machines, err := getEnvValue("ETCD_MACHINES")
			if err != nil || machines == "" {
				if err != nil {
					fmt.Fprintf(os.Stderr, err)
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
