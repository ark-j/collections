package binding

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"reflect"
)

const (
	urlEncodedContent    = "application/x-www-form-urlencoded"
	multipartFormContent = "multipart/form-data"
	jsonContent          = "application/json"
)

var MaxMemory int64 = 1024 * 1024 * 10

func Bind(req *http.Request, obj any) error {
	if err := ensurePointer(obj); err != nil {
		return err
	}
	mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if req.Method == http.MethodPatch ||
		req.Method == http.MethodPost ||
		req.Method == http.MethodPut ||
		len(mediaType) > 0 {
		switch mediaType {
		case urlEncodedContent:
			return Form(req, obj)
		case multipartFormContent:
			return MultipartForm(req, obj)
		case jsonContent:
			return JSON(req, obj)
		}
	}
	return errors.New("request method should be either PATCH, PUT, POST")
}

// Form is middleware to deserialize form-urlencoded data from the request.
// It gets data from the form-urlencoded body, if present, or from the
// query string. It uses the http.Request.ParseForm() method
// to perform deserialization, then reflection is used to map each field
// into the struct with the proper type. Structs with primitive slice types
// (bool, float, int, string) can support deserialization of repeated form
// keys, for example: key=val1&key=val2&key=val3
// An interface pointer can be added as a second argument in order
// to map the struct to a specific interface.
func Form(req *http.Request, formStruct interface{}) error {
	if err := ensurePointer(formStruct); err != nil {
		return err
	}
	formStructV := reflect.ValueOf(formStruct)
	err := req.ParseForm()
	// Format validation of the request body or the URL would add considerable overhead,
	// and ParseForm does not complain when URL encoding is off.
	// Because an empty request body or url can also mean absence of all needed values,
	// it is not in all cases a bad request, so let's return 422.
	if err != nil {
		return err
	}
	return mapForm(formStructV, req.Form, nil)
}

// MultipartForm works much like Form, except it can parse multipart forms
// and handle file uploads. Like the other deserialization middleware handlers,
// you can pass in an interface to make the interface available for injection
// into other handlers later.
func MultipartForm(req *http.Request, obj any) error {
	formStructV := reflect.ValueOf(obj)
	if req.MultipartForm == nil {
		if reader, err := req.MultipartReader(); err != nil {
			return err
		} else {
			form, err := reader.ReadForm(MaxMemory)
			if err != nil {
				return err
			}
			if req.Form == nil {
				if err := req.ParseForm(); err != nil {
					return err
				}
			}

			for k, v := range form.Value {
				req.Form[k] = append(req.Form[k], v...)
			}

			req.MultipartForm = form
		}
	}
	return mapForm(formStructV, req.MultipartForm.Value, req.MultipartForm.File)
}

// JSON binds object to json payload
func JSON(req *http.Request, obj any) error {
	return json.NewDecoder(req.Body).Decode(obj)
}
