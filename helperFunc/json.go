package helperFunc

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

func RespondJSON(w http.ResponseWriter, statusCode int, errMsg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var response map[string]interface{}

	if errMsg != "" {
		response = map[string]interface{}{
			"error": errMsg,
		}
	} else {
		response = data.(map[string]interface{})
	}

	ToJSON(response, w)
}

func Clone(i interface{}) interface{} {
	t := reflect.TypeOf(i)

	if t.Kind() != reflect.Ptr {
		return reflect.New(t).Interface()
	}

	return reflect.New(t.Elem()).Interface()
}
