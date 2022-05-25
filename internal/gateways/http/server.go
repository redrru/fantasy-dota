package http

import (
	"github.com/redrru/fantasy-dota/pkg/server"
)

var _ server.ServerInterface = (*Server)(nil)

type Server struct {
	usecase usecase
}

func NewServer(usecase usecase) *Server {
	return &Server{usecase: usecase}
}
