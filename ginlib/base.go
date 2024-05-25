package ginlib

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Server struct {
	engine *gin.Engine
	port   int
	url    string
}

type Opts struct {
	Port  string
	HTTPS bool
}

func NewServer(opts *Opts) *Server {
	i, _ := strconv.Atoi(opts.Port)
	return &Server{
		engine: gin.Default(),
		url:    fmt.Sprintf("127.0.0.1:%s", opts.Port),
		port:   i,
	}
}

// AddGetRoute adds a GET route to the server
func (s *Server) AddGetRoute(path string, handler func(*gin.Context)) {
	s.engine.GET(path, handler)
}

// AddPostRoute adds a POST route to the server
func (s *Server) AddPostRoute(path string, handler func(*gin.Context)) {
	s.engine.POST(path, handler)
}

// Run starts the server and tries different ports if the current url is in use
func (s *Server) Run() error {
	return s.engine.Run(s.url)
}

// GetPort returns the url that the server is running on
func (s *Server) GetPort() int {
	return s.port
}
