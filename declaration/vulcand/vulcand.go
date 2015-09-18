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
		etcdmachines, err = services.getEnvValue("ETCD_MACHINES")
		if err != nil {
			log.Fatal(err)
		} else {
			refreshEtcdMachines()
		}
	}

}
func (v *Vulcand) Register(s *service.Service) error {

}

func (v *Vulcand) GetService(s *service.Service) error {
	client := etcd.NewClient(Machines())
	defer client.Close()
	raw, err := client.RawGet(s.backendPath(), false, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw.Body, s)
	return err
}
func (v *Vulcand) UnRegister(s *service.Service) error {

}

func (v *Vulcand) SetPrivateKey(s *service.Service) error {

}
func (v *Vulcand) GetPrivateKey(s *service.Service) error {

}

func (v *Vulcand) ValidatePrivateKey(s *service.Service, r *http.Request) error {

}
func Machines() []string {
	return strings.Split(etcdmachines, "|")
}

func refreshEtcdMachines() {
	go func() {
		for _ = range time.Tick(10 * time.Second) {
			machines, err := getEnvValue("ETCD_MACHINES")
			if err != nil || machines == "" {
				if err != nil {
					fmt.Fprint(os.Stderr, err)
				}
				if machines == "" {
					fmt.Fprint(os.Stderr, "ETCD_MACHINES not set.")
				}
			} else {
				etcdmachines = machines
			}
		}
	}()
}
