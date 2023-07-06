package gapi

import (
	"fmt"

	"github.com/morka17/shiny_bank/v1/pb"
	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"github.com/morka17/shiny_bank/v1/src/token"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/morka17/shiny_bank/v1/src/worker"
)

// Server serves gRPC requests for banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot creare token maker: %v", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}
	
	
	return server, nil
}
