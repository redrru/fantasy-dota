package http

import (
	"github.com/redrru/fantasy-dota/pkg/server"
)

var _ server.ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}
