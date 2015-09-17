package service

import (
	"net/http"
)

var GetTypes = make(map[string]ParameterType)
var SetTypes = make(map[string]ParameterType)

func init() {
	AddGetType("form", GetFormValue)
	AddGetType("url", GetURLValue)
	AddGetType("json", GetJsonValue)
	AddGetType("query", GetQueryValue)
	AddSetType("form", SetFormValue)
	AddSetType("url", SetURLValue)
	AddSetType("json", SetJsonValue)
	AddSetType("query", SetQueryValue)
}

/*
Parameter defines a single parameter for the service to be called.

Title: Is a human readable title for the parameter, it will be used as a
value key for form values.

Type: The type is one of the following default types:
	"form" = normal form post value
	"url" = part of the url string (requires position)
	"json" = submitted as json
	"query" = as a query parameter in the url

Additional values can be added with the AddType function.

Position: is only relevant to url types, determines position in url.

Key: If set, this value overrides Title as key for value.

Required: Required value for service.

Description: A description of the parameter.

*/
type Parameter struct {
	Key         string
	Description string
	Type        string
	Position    string
	Required    bool
	DataType    string
}

func (p *Parameter) GetValue(val interface{}, r *http.Request) error {
	GetTypes[p.Type](val, p, r)
}

type ParameterType func(interface{}, Parameter, *http.Request) error

func AddGetType(key string, pt ParameterType) {
	GetTypes[key] = pt
}
func AddSetType(key string, pt ParameterType) {
	SetTypes[key] = pt
}

func GetFormValue(val interface{}, param Parameter, r *http.Request) error {

}
func SetFormValue(val interface{}, param Parameter, r *http.Request) error {
}

func GetURLValue(val interface{}, param Parameter, r *http.Request) error {
}
func SetURLValue(val interface{}, param Parameter, r *http.Request) error {
}

func GetJsonValue(val interface{}, param Parameter, r *http.Request) error {
}
func SetJsonValue(val interface{}, param Parameter, r *http.Request) error {
}

func GetQueryValue(val interface{}, param Parameter, r *http.Request) error {
}

func SetQueryValue(val interface{}, param Parameter, r *http.Request) error {
}

func SetValue(val interface{}, data interface{}, declaredtype string) error {
	rv := reflect.ValueOf(val)
	dv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Illegal value type, must be pointer and not nil")
	}
	realval := rv.Elem()
	dataval := rv.Elem()
	realtype := reflect.TypeOf(realval)
	datatype := reflect.TypeOf(data)
	if realtype != datatype {
		return fmt.Errorf("Type mismatch attempt to set %v as %v",
			reflect.TypeOf(data), reflect.TypeOf(realval))
	}
	if realtype.String() != declaredtype {
		return fmt.Errorf("Type mismatch declared Parameter DataType %v does not match value type %v.",
			declaredtype, realtype.String())
	}
	realval = dataval
	return nil

}
