package main

import (
	"encoding/json"
	"fmt"
	"github.com/vkozubal/in-memory-kv/store"
	"io/ioutil"
	"net/http"
)

var registry store.Registry

func init() {
	registry = store.Registry{}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Caught an error %#v\n", err)
	}
	request := KVRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Printf("Caught an error %#v\n", err)
	}

	fmt.Printf("requst: %#v\n", request)

	response := handleRequest(request)
	bytes, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Caught an error %#v\n", err)
	}

	fmt.Printf("response: %#v\n", response)

	_, err = w.Write(bytes)
	if err != nil {
		fmt.Printf("Caught an error %#v\n", err)
	}
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

	return KVReponse{} // todo handle
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

type RequestType string

const (
	GET    RequestType = "GET"
	SET    RequestType = "SET"
	DELETE RequestType = "DELETE"
	EXISTS RequestType = "EXISTS"
)

type KVRequest struct {
	Key    string      `json:"key"`
	Value  string      `json:"value"`
	Method RequestType `json:"method"`
}

type KVReponse struct {
	Key    string      `json:"key"`
	Method RequestType `json:"method"`
	Result string      `json:"result"`
	Error  string      `json:"error"`
	Value  string      `json:"value"`
}
