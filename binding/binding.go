package binding

import (
	"errors"
	"mime"

	// "mime/multipart"
	"net/http"
	"reflect"
)

const (
	urlEncodedContent    = "application/x-www-form-urlencoded"
	multipartFormContent = "multipart/form-data"
	jsonContent          = "application/json"
)

func Bind(req *http.Request, obj any) error {
	if err := ensurePointer(obj); err != nil {
		return err
	}
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if req.Method == http.MethodPost || req.Method == http.MethodPut || len(mediaType) > 0 {
		switch mediaType {
		case urlEncodedContent:
			return formBind(req, obj)
		case multipartFormContent:
			return multipartBind(req, obj, params)
		case jsonContent:
			return jsonBind(req, obj)
		}
	}
	return nil
}

func formBind(req *http.Request, obj any) error {
	formStruct := reflect.ValueOf(obj)
	_ = formStruct
	if err := req.ParseForm(); err != nil {
		return err
	}
	return nil
}

func multipartBind(req *http.Request, obj any, params map[string]string) error {
	_, _, _ = req, obj, params
	return nil
}

func jsonBind(req *http.Request, obj any) error {
	_, _ = req, obj
	return nil
}

func ensurePointer(obj any) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj must be pointer to the struct")
	}
	return nil
}

// func mapForm(formStruct reflect.Value, form map[string][]string, files map[string][]*multipart.FileHeader) error {
// 	_, _, _ = formStruct, form, files
// 	return nil
// }
