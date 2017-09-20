package public

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	server   http.Server
	router   *mux.Router
	listener net.Listener
}

func NewServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	r := mux.NewRouter()
	return &Server{
		server:   http.Server{Handler: r},
		router:   r,
		listener: l,
	}, err
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *Server) Router() *mux.Router {
	return s.router
}
