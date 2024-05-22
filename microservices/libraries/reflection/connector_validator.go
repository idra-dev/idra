package reflection

import (
	"log"
	"reflect"
	"strings"
)

type ValidationError struct{
	Field string
	Message string
}

func ValidateStruct(data interface{}) []ValidationError {
	var validationErrors []ValidationError
	valueType := reflect.TypeOf(data)
	valueConnector := reflect.ValueOf(data)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		fieldTags := make(map[string]string)

		tag := field.Tag
		for _, tagName := range strings.Split(string(tag), " ") {
			parts := strings.Split(tagName, ":")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				value = strings.ReplaceAll(value,"\"", "")
				fieldTags[key] = value
				if key == "binding" && value == "required"{
					typeField := field.Type.String()
					switch typeField {
					case "string":
						valueField := valueConnector.Field(i).String()
						if len(strings.TrimSpace(valueField))==0{
							validationErrors = AppendError(validationErrors, field)
						}
					case "int", "int16", "int8", "int32", "int64":
						if valueConnector.Field(i).IsZero() {
							validationErrors = AppendError(validationErrors, field)
						}
					case "float32", "float64":
						if valueConnector.Field(i).IsZero() {
							validationErrors = AppendError(validationErrors, field)
						}
					}
				}
			}
		}
	}

	return validationErrors
}

func AppendError(validationErrors []ValidationError, field reflect.StructField) []ValidationError {
	validationErrors = append(validationErrors, ValidationError{
		Field:   field.Name,
		Message: "Missing required field.",
	})
	return validationErrors
}

func GetColumnName(pstruct interface{}, pfield interface{}) string {
	v := reflect.ValueOf(pstruct).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Addr().Interface() == pfield {
			return v.Type().Field(i).Tag.Get("db")
		}
	}
	log.Fatalln("field not in struct")
	return ""
}

func GetFieldTags(data interface{}) map[string]map[string]string {
	result := make(map[string]map[string]string)

	valueType := reflect.TypeOf(data)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		fieldTags := make(map[string]string)

		tag := field.Tag

		for _, tagName := range strings.Split(string(tag), " ") {
			parts := strings.Split(tagName, ":")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				fieldTags[key] = value
			}
		}

		result[field.Name] = fieldTags
	}

	return result
}

func HasRequiredTag(field reflect.StructField) bool {
	tag := field.Tag.Get("required")
	return tag == "true" || tag == "1"
}

