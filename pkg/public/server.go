package public

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	ErrCommandSend = errors.New("sending command")
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	maxQueueWaitDuration = 1 * time.Second
)

type Server struct {
	server   http.Server
	router   *mux.Router
	listener net.Listener

	toStream chan *Command
}

type Command struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"args"`
}

func NewCommand(name string) *Command {
	return &Command{
		Arguments: map[string]interface{}{}, // <---- how ugly is this?!
		Name:      name,
	}
}

func NewServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	r := mux.NewRouter()
	return &Server{
		server:   http.Server{Handler: r},
		router:   r,
		listener: l,
		toStream: make(chan *Command),
	}, err
}

func (s *Server) SendCommand(cmd *Command) error {
	select {
	case s.toStream <- cmd:
		return nil
	case <-time.After(maxQueueWaitDuration):
	}

	return ErrCommandSend
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	wg := sync.WaitGroup{}
	closeCh := make(chan bool)

	// trash any incoming msgs, they should not happen
	// increment this when it happens :D
	// 1
	wg.Add(1)
	go func() {
		wg.Done()
		defer close(closeCh)

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
		}
	}()

	log.Println("OBS connected")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-closeCh:
				ws.Close()
				return
			case msg := <-s.toStream:
				log.Println("command:", msg.Name)
				if err = ws.WriteJSON(msg); err != nil {
					log.Println("json:", err)
					return
				}
			}
		}
	}()

	wg.Wait()
}

func (s *Server) Serve() error {
	// register this for last so it can override any route
	s.router.HandleFunc("/ws", s.websocketHandler)
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./data/html")))
	return s.server.Serve(s.listener)
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *Server) Router() *mux.Router {
	return s.router
}
