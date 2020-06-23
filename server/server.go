package server

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shishkebaber/message-api/model"
	"github.com/shishkebaber/message-api/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	storage storage.Storage
	router  *mux.Router
	Logger  *logrus.Logger
	v       *model.Validation
}

func NewServer(stor storage.Storage, r *mux.Router) *Server {
	l := logrus.New()
	v := model.NewValidation()
	s := &Server{stor, r, l, v}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) routes() {
	getR := s.router.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/messages", s.handleMessagesGet())
	getR.HandleFunc("/messages/{name}", s.handleMessageSingleGet())

	postR := s.router.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/messages", s.validateCreateRequest(s.handleCreateMessage()))
}

func (s *Server) handleMessagesGet() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Getting messages")
		rw.Header().Add("Content-Type", "application/json")
		result, err := s.storage.FetchAllMessages()
		if err != nil {
			s.Logger.Error("Error during get: ", err)
			s.respond(rw, r, struct{}{}, http.StatusInternalServerError)
		}
		s.respond(rw, r, result, http.StatusOK)
	}
}

func (s *Server) handleMessageSingleGet() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Getting message")
		rw.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)

		name := vars["name"]

		result, err := s.storage.FetchMessage(name)
		if err != nil {
			s.Logger.Error("Error during get single: ", err)
			s.respond(rw, r, &GenericError{Message: err.Error()}, http.StatusNotFound)
		}
		s.respond(rw, r, result, http.StatusOK)
	}
}

func (s *Server) validateCreateRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		s.Logger.Info("Validating message")
		message := &model.Message{}
		err := json.NewDecoder(r.Body).Decode(message)
		if err != nil {
			s.Logger.Error("Error during validation: ", err)
			s.respond(rw, r, &GenericError{Message: err.Error()}, http.StatusBadRequest)
			return
		}

		errs := s.v.Validate(message)
		if len(errs) != 0 {
			s.respond(rw, r, &ValidationError{Messages: errs.Errors()}, http.StatusUnprocessableEntity)
			return
		}

		ctx := context.WithValue(r.Context(), model.MessageKey{}, message)
		r = r.WithContext(ctx)

		h(rw, r)
	}
}

func (s *Server) handleCreateMessage() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Creating message")
		input := r.Context().Value(model.MessageKey{}).(*model.Message)
		err := s.storage.CreateMessage(input)
		if err != nil {
			s.Logger.Error("Error during creating: ", err)
			s.respond(rw, r, &GenericError{Message: err.Error()}, http.StatusBadRequest)
		}
		s.respond(rw, r, nil, http.StatusOK)
	}
}

func (s *Server) respond(rw http.ResponseWriter, r *http.Request, data interface{}, status int) {
	rw.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(rw).Encode(data)
		if err != nil {
			s.Logger.Error("Error during encoding structure to Json")
		}
	}
}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}
