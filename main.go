package main

import (
	"encoding/json"
	"fmt"
	"github.com/vkozubal/in-memory-kv/store"
	"io/ioutil"
	"net/http"
)

type RequestType string

const (
	GET    RequestType = "GET"
	SET    RequestType = "SET"
	DELETE RequestType = "DELETE"
	EXISTS RequestType = "EXISTS"
)

type KVRequest struct {
	Key    string      `json:"key,omitempty"`
	Value  string      `json:"value,omitempty"`
	Method RequestType `json:"method,omitempty"`
}

type KVReponse struct {
	Key    string      `json:"key,omitempty"`
	Value  string      `json:"value,omitempty"`
	Method RequestType `json:"method,omitempty"`
	Result string      `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

var registry store.Registry

func init() {
	var err error
	if registry, err = store.NewRegistry(); err != nil {
		panic(err)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	request, err := parseRequest(r)
	if err != nil {
		errMessage := fmt.Sprintf("Couldn't parse JSON format. %#v\n", err)
		data := KVReponse{Error: errMessage}

		if err = writeResponse(err, data, w); err != nil {
			panic(err)
		}
		return
	}

	response := handleRequest(request)
	err = writeResponse(err, response, w)
	if err != nil {
		panic(err)
	}
}

func writeResponse(err error, response KVReponse, w http.ResponseWriter) error {
	bytes, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	fmt.Printf("response: %#v\n", response)
	return nil
}

func parseRequest(r *http.Request) (KVRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return KVRequest{}, err
	}
	request := KVRequest{}
	if err = json.Unmarshal(body, &request); err != nil {
		return KVRequest{}, err
	}
	fmt.Printf("requst: %#v\n", request)
	return request, err
}

func handleRequest(request KVRequest) KVReponse {
	switch request.Method {
	case GET:
		val, err := registry.Get(request.Key)
		if err != nil {
			return errorResponseFrom(GET, request, err)
		}
		return KVReponse{Method: GET, Key: request.Key, Value: val}
	case SET:
		err := registry.Set(request.Key, request.Value)
		if err != nil {
			return errorResponseFrom(SET, request, err)
		}
		return KVReponse{Method: SET, Key: request.Key, Result: "success"}
	case DELETE:
		err := registry.Delete(request.Key)
		if err != nil {
			return errorResponseFrom(DELETE, request, err)
		}
		return KVReponse{Method: DELETE, Key: request.Key, Result: "success"}
	case EXISTS:
		if found := registry.Exists(request.Key); found {
			return KVReponse{Method: EXISTS, Key: request.Key, Result: "found"}
		}
		return KVReponse{Method: EXISTS, Key: request.Key, Error: "not found"}
	}

	return KVReponse{Method: request.Method, Key: request.Key, Error: "method not found"}
}

func errorResponseFrom(method RequestType, request KVRequest, err error) KVReponse {
	return KVReponse{Method: method, Key: request.Key, Error: err.Error()}
}

func main() {
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("Couldn't start a server %#v\n", err)
	}
}
