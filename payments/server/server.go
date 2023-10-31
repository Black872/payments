package server

import (
	"context"
	"net"
	"net/http"
	"payments/server/handlers"

	"github.com/gorilla/mux"
	"github.com/zeebo/errs"
)

var serverError = errs.Class("web server")

type Server struct {
	listener net.Listener
	server   http.Server
}

func NewServer(ctx context.Context, listener net.Listener, h *handlers.Handlers) *Server {
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", h.Auth.SignUp).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", h.Auth.Login).Methods(http.MethodPost)

	return &Server{
		listener: listener,
		server: http.Server{
			Handler: router,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
		},
	}
}

func (s *Server) Run() (err error) {
	if err = s.server.Serve(s.listener); err != nil {
		return serverError.Wrap(err)
	}
	return nil
}

func (s *Server) Close(ctx context.Context) (err error) {
	return serverError.Wrap(s.server.Shutdown(ctx))
}
