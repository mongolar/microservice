package service

type PrivateService struct {
	Service
	privatekey string
	hash       string
	leader     bool
}

const ETCDSERVICEKEYS = "/mongolar/service/private/keys"

func (ps *PrivateService) Serve(w http.ResponseWriter, r *http.Request) {
	if ps.validatePrivate(r) {
		ps.Service.Handler(w, r)
	} else {
		http.Error(w.Writer, "Forbidden", 403)
	}
}

func (ps *PrivateService) Register() {

}

func (ps *PrivateService) follow() {

}

func (ps *PrivateService) lead() {

}

func (ps *PrivateService) validatePrivate() bool {
	return false
}

func getPrivateHash() string {

}
