package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
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

func ShouldBind(request *http.Request, writer http.ResponseWriter, data any) error {
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
