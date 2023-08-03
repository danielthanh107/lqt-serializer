package lqtserializer

import (
	"encoding/json"
	"reflect"
)

func Struct2Map(s interface{}, m *map[string]interface{}) error {
	jStruct, err := json.Marshal(s)

	if err != nil {
		return err
	}

	err = json.Unmarshal(jStruct, &m)

	if err != nil {
		return err
	}

	return nil
}

func GetAllFieldsInStruct(s interface{}) []string {
	fields := []string{}

	t := reflect.TypeOf(s)

	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}

	return fields
}

func InterfaceToSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	// Keep the distinction between nil and empty slice sinput
	if s.IsNil() {
		return nil
	}

	if s.Kind() == reflect.Pointer {
		s = s.Elem()
	}

	if s.Kind() != reflect.Slice {
		panic("Interface{} is not slice")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func MergeMaps[K comparable, V any](m1 map[K]V, m2 map[K]V) map[K]V {
	merged := make(map[K]V)
	for key, value := range m1 {
		merged[key] = value
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}

func SliceMaps(m map[string]any, fields []string) map[string]any {
	sliced := map[string]any{}

	for _, field := range fields {
		value, ok := m[field]

		if ok {
			sliced[field] = value
		}
	}

	return sliced
}
