package parameter

import ()

var DataTypes map[string]reflect.Type = make(map[string]type)

func init(){
	AddDataType('int', 0)
	AddDataType('string', "")
	AddDataType('[]string', make([]string, 0))
	AddDataType('[]int', make([]int,0))
	AddDataType('map[string]string', make(map[string]string))
	AddDataType('map[string]int', make(map[string]int))
}

func AddDataType(key string, t interface{}){
	DataTypes[key] = reflect.TypeOf{t}
}

func ValidateDataType(key string, t interface{}) bool {
	if DataTypes[key] == reflect.TypeOf{t} {
		return true
	}
	return false
}

func ValidatePtrDataType(key string, t interface{}) bool {
        refval := reflect.ValueOf(t)
        if refval.Kind() != reflect.Ptr || refval.IsNil() {
                return false
        }
	return ValidateDataType(key, t)
}

func SetPointer(key string, receiver interface{}, sender interface{}) error {
	if !ValidatePtrDataType(key, receiver){
		return fmt.Errorf("Receiving value not right type")
	}
	if !ValidateDataType(key, sender){
		return fmt.Errorf("Sending value not right type")
	}
	rec := reflect.ValueOf(receiver).Elem()
	send := reflect.ValueOf(sender)
	// Value assignment will need to be further evaluated.
	rec.Set(send)
}
