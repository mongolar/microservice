package request

import (
	"github.com/mongolar/microservice/environment"
	"github.com/mongolar/microservice/service"
	"net/http"
	"reflect"
)

type Request struct {
	Service service.Service
	Request *http.Request
	Values  map[string]reflect.Value
}

func New() *Request {
	r := new(Request)
	return r
}

func NewSender(s service.Service) *Request {
	r := Request{Service: s}
	url := fmt.Printf("%v/%v.%v.%v", environment.IntServiceURL(), r.Service.Domain, r.Service.Title, r.Title.Version)
	r.Request = http.NewRequest(r.Service.Method, url, nil)
	r.Values = make(map[string]reflect.Value)
	r.getValues()
	return &r
}

func NewReceiver(s service.Service, re *http.Request) *Request {
	r := Request{s, re, make(map[string]reflect.Value)}
	r.getRValuesRequest()
	return &r
}

func (r *Request) getRValues() {
	for _, p := range r.Service.Parameters {
		r.Values[p.Key] = p.Value()
	}
}

func (r *Request) getValuesFromRequest() error {
	for k, p := range r.Service.Parameters {
		val = r.Values[k]
		param, err = r.Service.Parameters.Get(key)
		if err != nil {
			return err
		}
		param.GetValue(val, r.Request)
	}
}

func (r *Request) getRValuesRequest() error {
	for _, p := range r.Service.Parameters {
		val := p.Value()
		err := p.GetValue(val, r.Request)
		if err != nil {
			return err
		}
		r.Values[p.Key] = val
	}
	return nil
}

func (r *Request) SetValue(key string, val interface{}) error {
	p, err := r.Service.Parameters.Get(key)
	if err != nil {
		return err
	}
	r.Values[key] = reflect.ValueOf(val)
}

func (r *Request) GetValue(key string, val interface{}) error {
	if pt, ok := Values[key]; ok {
		if parameter.ValidatePtrDataType(key, val) {
			baseval := reflect.ValueOf(val).Elem()
			baseval.Set(p.Values[key].Elem())
			return nil
		}
		return fmt.Errorf("Invalid Data type")
	}
	param, err = r.Service.Parameters.Get(key)
	if err != nil {
		return err
	}
	newval := param.Value()
	err := param.GetValue(newval, r.Request)
	return err
}

func (r *Request) Call() error {

}
