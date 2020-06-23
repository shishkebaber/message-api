package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shishkebaber/message-api/model"
	"github.com/shishkebaber/message-api/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

var s *Server

func init() {
	stor := storage.NewMStorage()
	router := mux.NewRouter().StrictSlash(true)
	s = NewServer(stor, router)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code: %d. Got: %d\n", expected, actual)
	}
}

func createData() {
	s.storage.(*storage.MStorage).Messages["test"] = &model.Message{"test", "test payload", model.Metadata{"2020-12-12 19:02:33.123", "10eec78c-b4a3-11ea-b3de-0242ac130004"}}
}

func clearData() {
	s.storage.(*storage.MStorage).Messages = make(map[string]*model.Message)
}

func TestCreateMessage(t *testing.T) {
	clearData()
	var jsonStr = []byte(`{"name":"bbr", "payload":"payload test", "metadata":{"timestamp":"2020-12-12 19:02:33.123", "uuid":"10eec78c-b4a3-11ea-b3de-0242ac130004"}}`)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetMessages(t *testing.T) {
	clearData()
	keys := []string{"name", "payload"}
	values := []string{"test", "test payload"}

	createData()

	req, _ := http.NewRequest("GET", "/messages", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var rM []map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &rM)
	if err != nil {
		s.Logger.Error("Unable to unmarshal get result: ", err)
	}

	for i, v := range keys {
		if _, ok := rM[0][v]; !ok {
			t.Errorf("Expected existance of key %s. ", v)
		}
		if rM[0][v] != values[i] {
			t.Errorf("Expected value %s. Got: %s", values[i], rM[0][v])
		}
	}
}
