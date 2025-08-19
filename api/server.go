package api

import (
	db "simplebank/db/sqlc"
	"simplebank/worker"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
)

type Server struct {
	store           db.Store
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, taskDistributor worker.TaskDistributor) *Server {
	server := &Server{
		store:           store,
		taskDistributor: taskDistributor,
	}
	router := gin.Default()
	router.Use(apmgin.Middleware(router))
	// TODO: add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/users", server.createUser)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
