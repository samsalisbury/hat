package hat

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type Server struct {
	root *Node
}

func NewServer(root interface{}) (*Server, error) {
	if rootNode, err := newNode(nil, nil, reflect.TypeOf(root)); err != nil {
		return nil, err
	} else {
		return &Server{rootNode}, nil
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if statusCode, resource, err := s.root.Render(r.URL.Path, r.Method, newPayload(r)); err != nil {
		writeError(w, err)
	} else {
		writeResponse(w, statusCode, resource)
	}
}

func writeError(w http.ResponseWriter, err error) {
	if httpErr, ok := err.(HTTPError); ok {
		writeResponse(w, httpErr.StatusCode(), httpErr.Err())
	}
}

func writeResponse(w http.ResponseWriter, statusCode int, resource interface{}) {
	if data, err := json.Marshal(resource); err != nil {
		writeError(w, HttpError(500, "Unable to marshal response into json:", err.Error()))
	} else {
		w.WriteHeader(statusCode)
		w.Write(data)
	}
}
