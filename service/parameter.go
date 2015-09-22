package service

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
	if rval.Kind() != reflect.Ptr || rval.IsNil() {
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

func (ps Parameters) GetParam(key string) (*Parameter, error) {
	for _, p := range ps {
		if p.Key == key {
			return &p, nil
		}
	}
	return new(Parameter), fmt.Errorf("Parameter %v not found", key)
}

func SetValue(receiver interface{}, data interface{}, datatype string) error {
	err := validateReceiver(receiver, datatype)
	if err != nil {
		return err
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

type Value struct {
	Receiver interface{}
	Element
	DataType
	RefValue reflect.Value
}
