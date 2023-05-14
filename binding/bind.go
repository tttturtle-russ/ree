package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

const (
	TypeJson      = "application/json"
	TypeXML       = "text/xml"
	TypeText      = "text/plain"
	TypeHTML      = "text/html"
	TypeForm      = "application/x-www-form-urlencoded"
	TypeMultipart = "multipart/form-data"
)

func filterContentType(mineType string) string {
	for i, v := range mineType {
		if v == ' ' || v == ';' {
			return mineType[:i]
		}
	}
	return mineType
}

func mapBindFiled(form url.Values, data any) map[string]any {
	var typeOf reflect.Type
	var valueOf reflect.Value
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		typeOf = reflect.TypeOf(data)
		valueOf = reflect.ValueOf(&data)
	} else {
		typeOf = reflect.TypeOf(*data.(*any))
		valueOf = reflect.ValueOf(data)
	}
	if valueOf.Elem().CanSet() {
		valueOf = valueOf.Elem()
	} else {
		return nil
	}
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		type_ := typeOf.Field(i)
		if field.CanSet() && type_.IsExported() {
			get := form.Get(type_.Name)
			field.Set(reflect.ValueOf(get))
		}
	}
	return nil
}

func Bind(request *http.Request, writer http.ResponseWriter, data any) error {
	contentType := filterContentType(request.Header.Get("Content-Type"))
	switch contentType {
	case TypeJson:
		return BindJSON(request, data)
	case TypeXML:
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
	if err := request.ParseMultipartForm(1 << 31); err != nil {
		return err
	}
	mapBindFiled(request.PostForm, data)
	return nil
}
