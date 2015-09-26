package parameter

import (
	//"errors"
	"fmt"
	"log"
	"net/http"
	//"reflect"
)

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
	Position    string
}

func (p *Parameter) Value() reflect.Value {
	dt = DataTypes[p.Key]
	return New(dt)
}

func (p *Parameter) Set(val interface{}, r *http.Request) error {
	err := validateReceiver(val, p.DataType)
	if err != nil {
		return err
	}
	pt, err := getParameterType(p.Type)
	if err != nil {
		return err
	}
	return pt.Set(val, r, p)
}

func (p *Parameter) Get(val interface{}, r *http.Request) error {
	err := validateReceiver(val, p.DataType)
	if err != nil {
		return err
	}
	pt, err := getParameterType(p.Type)
	if err != nil {
		return err
	}
	return pt.Get(val, r, p)
}

func validateReceiver(val interface{}, datatype string) err {
	refval := reflect.ValueOf(val)
	if refval.Kind() != reflect.Ptr || refval.IsNil() {
		return errors.New("Illegal value type, must be pointer and not nil")
	}
	e := rv.Elem()
	valtype := reflect.TypeOf(e).String()
	if valtype != datatype {
		return fmt.Errorf("Illegal value type, %v does not match %v", valtype, datatype)
	}
	return nil
}

type Parameters []Parameter

func (ps Parameters) Validate() {
	for _, p := range ps {
		if p.Type == "" {
			log.Fatal(fmt.Errorf("Parameter Type not set for %v", p.Key))
		}
		if p.DataType == "" {
			log.Fatal(fmt.Errorf("Data Type not set for %v", p.Key))
		}

	}

}

func (ps Parameters) GetRequired() []string {
	required := make([]string, 0)
	for _, p := range ps {
		if p.Required {
			required = append(required, p.Key)
		}
	}
	return required
}

func (ps Parameters) Get(key string) (*Parameter, error) {
	for _, p := range ps {
		if p.Key == key {
			return &p, nil
		}
	}
	return new(Parameter), fmt.Errorf("Parameter %v not found", key)
}
