package binding

import (
	"errors"
	"mime/multipart"
	"reflect"
	"strconv"

	"github.com/google/go-cmp/cmp"
)

func ensurePointer(obj any) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj must be pointer to the struct")
	}
	return nil
}

// Takes values from the form data and puts them into a struct
func mapForm(
	formStruct reflect.Value,
	form map[string][]string,
	formfile map[string][]*multipart.FileHeader,
) error {
	if formStruct.Kind() == reflect.Ptr {
		formStruct = formStruct.Elem()
	}

	typ := formStruct.Type()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := formStruct.Field(i)

		if typeField.Type.Kind() == reflect.Ptr && typeField.Anonymous {
			structField.Set(reflect.New(typeField.Type.Elem()))
			if err := mapForm(structField.Elem(), form, formfile); err != nil {
				return err
			}
			if cmp.Diff(
				structField.Elem().Interface(),
				reflect.Zero(structField.Elem().Type()).Interface(),
			) == "" {
				structField.Set(reflect.Zero(structField.Type()))
			}
		} else if typeField.Type.Kind() == reflect.Struct {
			if err := mapForm(structField, form, formfile); err != nil {
				return err
			}
		}

		inputFieldName := parseFormName(typeField.Name, typeField.Tag.Get("form"))
		if len(inputFieldName) == 0 || !structField.CanSet() {
			continue
		}

		inputValue, exists := form[inputFieldName]
		if exists {
			numElems := len(inputValue)
			if structField.Kind() == reflect.Slice && numElems > 0 {
				sliceOf := structField.Type().Elem().Kind()
				slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
				for j := 0; j < numElems; j++ {
					if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j), inputFieldName); err != nil {
						return err
					}
				}
				formStruct.Field(i).Set(slice)
			} else {
				if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField, inputFieldName); err != nil {
					return err
				}
			}
			continue
		}

		inputFile, exists := formfile[inputFieldName]
		if !exists {
			continue
		}
		fhType := reflect.TypeOf((*multipart.FileHeader)(nil))
		numElems := len(inputFile)
		if structField.Kind() == reflect.Slice && numElems > 0 &&
			structField.Type().Elem() == fhType {
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				slice.Index(i).Set(reflect.ValueOf(inputFile[i]))
			}
			structField.Set(slice)
		} else if structField.Type() == fhType {
			structField.Set(reflect.ValueOf(inputFile[0]))
		}
	}
	return nil
}

// This sets the value in a struct of an indeterminate type to the
// matching value from the request (via Form middleware) in the
// same type, so that not all deserialized values have to be strings.
// Supported types are string, int, float, bool, and ptr of these types.
func setWithProperType(
	valueKind reflect.Kind,
	val string,
	structField reflect.Value,
	nameInTag string,
) error {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		} else {
			structField.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val == "" {
			val = "0"
		}
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		} else {
			structField.SetUint(uintVal)
		}
	case reflect.Bool:
		if val == "on" {
			structField.SetBool(true)
			break
		}

		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		} else if boolVal {
			structField.SetBool(true)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	case reflect.Ptr:
		newVal := reflect.New(structField.Type().Elem())
		if err := setWithProperType(structField.Type().Elem().Kind(), val, newVal.Elem(), nameInTag); err != nil {
			return err
		}
		structField.Set(newVal)
	}
	return nil
}

// Checks if actual ie tag is not zero if so returns it
// If tag is not present it directly return Fieldname
func parseFormName(raw, actual string) string {
	if len(actual) > 0 {
		return actual
	}
	return raw
}
