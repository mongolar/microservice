package service

import (
	//"errors"
	"fmt"
	"log"
	"net/http"
	//"reflect"
)

func init() {

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
	Required    bool
	DataType    string
	Method      string
	OtherData   map[string]string
	pt          ParameterType
}

func (p *Parameter) GetValue(val interface{}, r *http.Request) error {
	return fmt.Errorf("")
}

type Parameters []Parameter

func (ps Parameters) Validate() {
	for _, p := range ps {
		if p.Type == "" {
			log.Fatal(fmt.Errorf("Parameter Type not set for %v", p.Key))
		}

	}

}

func (ps Parameters) GetParam(key string) (*Parameter, error) {
	for _, p := range ps {
		if p.Key == key {
			return &p, nil
		}
	}
	return new(Parameter), fmt.Errorf("Parameter %v not found", key)
}

func SetValue(receiver interface{}, data interface{}, declaredtype string) error {
	/*
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
	*/
	return nil
}
