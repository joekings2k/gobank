package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/joekings2k/gobank/db/sqlc"
)

// server serves http requests for our banking service
type Server  struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.GET("/",server.checkHealth )
	router.POST("/accounts",server.createAccount)
	router.GET("/accounts/:id",server.getAccount )
	router.GET("/accounts",server.ListAccounts )
	router.PATCH("/accounts/:id",server.updateAccount)
	router.DELETE("/accounts/:id",server.deleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(addres string)error{
	return server.router.Run(addres)
}

func errorResponse (err error) gin.H{
	return  gin.H{"error":err.Error()}
}
