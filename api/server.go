package api

import (
	db "simplebank/db/sqlc"
	"simplebank/worker"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.elastic.co/apm/module/apmgin/v2"
)

type Server struct {
	store           db.Store
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
	redisClient     *redis.Client // Direct Redis access
}

func NewServer(store db.Store, taskDistributor worker.TaskDistributor, redisClient *redis.Client) *Server {
	server := &Server{
		store:           store,
		taskDistributor: taskDistributor,
		redisClient:     redisClient, // Direct Redis access
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
