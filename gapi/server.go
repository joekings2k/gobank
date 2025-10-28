package gapi

import (
	"fmt"

	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/pb"
	"github.com/joekings2k/gobank/token"
	"github.com/joekings2k/gobank/util"
)

// server serves http requests for our banking service
type Server  struct {
	pb.UnimplementedGobankServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error ) {
	tokenMaker,err := token.NewPasetoMaker(config.TokenSymmeticKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err) 
	}
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}
	
	
	
	return server, nil
}


