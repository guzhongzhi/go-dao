package utils

import (
	"encoding/json"
	"reflect"
)

func MapToJSON(callType reflect.Type) interface{} {
	if callType.NumIn() <= 1 {
		return "{}"
	}
	inType := callType.In(1)
	fields := LoopType("json", inType)
	js, _ := json.MarshalIndent(fields, "", "    ")
	return js
}

func LoopType(tagName string, inType reflect.Type) interface{} {
	fields := make(map[string]interface{})
	if inType.Kind() == reflect.Ptr {
		inType = inType.Elem()
	}
	if inType.Kind() == reflect.Struct {
		num := inType.NumField()
		for i := 0; i < num; i++ {
			f := inType.Field(i)
			name := f.Tag.Get(tagName)
			if name == "" {
				continue
			}

			if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct {
				fields[name] = LoopType(tagName, f.Type.Elem())
			} else if f.Type.Kind() == reflect.Struct {
				fields[name] = LoopType(tagName, f.Type)
			} else {
				fields[name] = f.Type.String()
			}
		}
	} else {
		return inType.String()
	}
	return fields
}
