package lqtserializer

import (
	"reflect"
)

type IListSerializer interface {
	GetData() *[]map[string]any
}

type ListSerializer struct {
	Records reflect.Value
}

func (serList *ListSerializer) GetData() *[]map[string]any {
	result := []map[string]any{}
	records := serList.Records

	for i := 0; i < records.Len(); i++ {
		record := records.Index(i)
		recordType := record.Type()

		// get data to map
		m := map[string]any{}

		// Check fields return
		fieldMethod, okFields := recordType.MethodByName("Fields")

		for j := 0; j < recordType.NumField(); j++ {
			field := recordType.Field(j)
			fieldValue := record.Field(j)

			if field.Name == "Model" {
				modelMap := map[string]any{}
				convertToMap(fieldValue.Interface(), &modelMap)

				m = MergeMaps(m, modelMap)
			} else if field.Tag.Get("valueIn") == "method" {
				methodName := "Get" + field.Name
				method := record.MethodByName(methodName)

				if !method.IsValid() && record.CanAddr() {
					method = record.Addr().MethodByName(methodName)

					if !method.IsValid() {
						panic("Not found method " + methodName + " in " + recordType.Name())
					}
				} else {
					panic("Not found method " + methodName + " in " + recordType.Name())
				}
				value := method.Call([]reflect.Value{})[0]

				if field.Tag.Get("json") != "" {
					m[field.Tag.Get("json")] = value.Interface()
				} else {
					m[field.Name] = value.Interface()
				}
			}
		}

		if okFields {
			fields := fieldMethod.Func.Call([]reflect.Value{record})[0].Interface().([]string)

			m = SliceMaps(m, fields)
		}

		result = append(result, m)
	}

	return &result
}

func New(ser interface{}, records interface{}) IListSerializer {
	rtSer := reflect.TypeOf(ser)

	if rtSer.Kind() != reflect.Pointer {
		panic("Serializer must be pointer to slice")
	}

	rtSer = rtSer.Elem()

	if rtSer.Kind() != reflect.Slice {
		panic("Serializer must be pointer to slice")
	}

	rvSer := reflect.ValueOf(ser).Elem()
	rtElemSer := rtSer.Elem()

	if rtElemSer.Kind() != reflect.Struct {
		panic("Element in Slice Serializer must be struct")
	}

	rvRecords := reflect.ValueOf(records)

	if rvRecords.Kind() == reflect.Pointer {
		rvRecords = rvRecords.Elem()
	}

	if rvRecords.Kind() != reflect.Slice {
		panic("Record must be slice")
	}

	for i := 0; i < rvRecords.Len(); i++ {
		serElemValue := reflect.New(rtElemSer).Elem()
		serModelField := serElemValue.FieldByName("Model")

		record := rvRecords.Index(i)
		serModelField.Set(record)

		rvSer = reflect.Append(rvSer, serElemValue)
	}

	return &ListSerializer{Records: rvSer}
}

// local method
func convertToMap(ser interface{}, m *map[string]any) error {
	err := Struct2Map(ser, m)

	if err != nil {
		return err
	}

	return nil
}
