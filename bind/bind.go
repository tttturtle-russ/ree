package bind

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
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
	case "application/json":
		return BindJSON(request, data)
	case "text/xml":
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
