package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/token"
	"github.com/joekings2k/gobank/util"
)

// server serves http requests for our banking service
type Server  struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
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
	
	if v,ok :=binding.Validator.Engine().(*validator.Validate);ok {
		v.RegisterValidation("currency",validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// check health
	router.GET("/",server.checkHealth )
 
	//user routes 
	router.POST("/users",server.createUser)
	router.POST("/users/login",server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts",server.createAccount)
	authRoutes.GET("/accounts/:id",server.getAccount )
	authRoutes.GET("/accounts",server.ListAccounts )
	authRoutes.PATCH("/accounts/:id",server.updateAccount)
	authRoutes.DELETE("/accounts/:id",server.deleteAccount)

	// transfers
	authRoutes.POST("/transfers", server.createTransfer)
	server.router = router
}

func (server *Server) Start(addres string)error{
	return server.router.Run(addres)
}

func errorResponse (err error) gin.H{
	return  gin.H{"error":err.Error()}
}
