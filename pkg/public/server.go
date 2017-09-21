package public

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	server   http.Server
	router   *mux.Router
	listener net.Listener
}

var upgrader = websocket.Upgrader{}

func NewServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	r := mux.NewRouter()
	return &Server{
		server:   http.Server{Handler: r},
		router:   r,
		listener: l,
	}, err
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

}

func (s *Server) Serve() error {
	// register this for last so it can override any route
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./data/html")))
	s.router.HandleFunc("/ws", s.websocketHandler)
	return s.server.Serve(s.listener)
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *Server) Router() *mux.Router {
	return s.router
}
