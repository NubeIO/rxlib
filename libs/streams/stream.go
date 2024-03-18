package streams

import (
	"fmt"
	"net"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

type PortStream struct {
	Key    string
	Ports  []string
	Value  float64
	Value2 float64
}

type Client struct {
	ID      string
	conn    net.Conn
	encoder *msgpack.Encoder
	decoder *msgpack.Decoder
}

func NewClient(conn net.Conn, id string) *Client {
	return &Client{
		ID:      id,
		conn:    conn,
		encoder: msgpack.NewEncoder(conn),
		decoder: msgpack.NewDecoder(conn),
	}
}

func (c *Client) SendMessage(port *PortStream) error {
	return c.encoder.Encode(port)
}

func (c *Client) ReceiveMessage() (*PortStream, error) {
	var port PortStream
	if err := c.decoder.Decode(&port); err != nil {
		return nil, err
	}
	return &port, nil
}

func (c *Client) ReceiveMessages(handler func(stream *PortStream)) {
	for {
		var receivedPort *PortStream
		if err := c.decoder.Decode(&receivedPort); err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}
		handler(receivedPort)
	}
}

func (c *Client) Close() error {
	fmt.Println("CLIENT CLOSE")
	return c.conn.Close()
}

type Server struct {
	listener net.Listener
	clients  map[string]*Client
	mu       sync.Mutex
}

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: listener,
		clients:  make(map[string]*Client),
	}, nil
}

func (s *Server) Accept() (*Client, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}
	client := NewClient(conn, fmt.Sprintf("%v", conn.RemoteAddr()))
	s.mu.Lock()
	s.clients[client.ID] = client
	s.mu.Unlock()
	return client, nil
}

func (s *Server) GetClient(id string) (*Client, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	client, ok := s.clients[id]
	return client, ok
}

func (s *Server) Close() error {
	return s.listener.Close()
}
