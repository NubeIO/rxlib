package rxlib

import (
	"reflect"
)

func (inst *RuntimeImpl) ToStringArray(interfaces interface{}) []string {
	if interfaces == nil {
		return nil
	}
	v := reflect.ValueOf(interfaces)
	if v.Kind() != reflect.Slice {
		return nil
	}
	var strings []string
	for i := 0; i < v.Len(); i++ {
		iface := v.Index(i).Interface()
		str, ok := iface.(string)
		if !ok {
			return nil
		}
		strings = append(strings, str)
	}
	return strings
}
