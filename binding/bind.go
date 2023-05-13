package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"reflect"
)

const (
	TYPEJSON      = "application/json"
	TYPEXML       = "text/xml"
	TYPETEXT      = "text/plain"
	TYPEHTML      = "text/html"
	TYPEFORM      = "application/x-www-form-urlencoded"
	TYPEMULTIPART = "multipart/form-data"
)

func filterContentType(mineType string) string {
	for i, v := range mineType {
		if v == ' ' || v == ';' {
			return mineType[:i]
		}
	}
	return mineType
}

func mapBindFiled(data any) map[string]any {
	m := make(map[string]any)
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		data = &data
	}
	valueOf := reflect.ValueOf(data)
	if valueOf.Elem().CanSet() {
		valueOf = valueOf.Elem()
	} else {
		return nil
	}
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		if field.CanSet() {

		}
	}
}

func Bind(request *http.Request, writer http.ResponseWriter, data any) error {
	contentType := filterContentType(request.Header.Get("Content-Type"))
	switch contentType {
	case TYPEJSON:
		return BindJSON(request, data)
	case TYPEXML:
		return BindXML(request, data)
	default:
		return errors.New("unsupported mine type")
	}
}

func BindJSON(request *http.Request, data interface{}) error {
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &data)
	return err
}

func BindXML(request *http.Request, data any) error {
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(bytes, &data)
	return err
}

func BindFORM(request *http.Request, data any) error {
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	if err = request.ParseForm(); err != nil {
		return err
	}

}
