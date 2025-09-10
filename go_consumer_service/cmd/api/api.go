package api

import (
	"go_consumer_service/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	listenAddr string
	dbClient   *mongo.Client
}

// NewAPIServer creates a new API server
func NewAPIServer(listenAddr string, dbClient *mongo.Client) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		dbClient:   dbClient,
	}
}

func (s *APIServer) Start() error {
	router := gin.Default()

	rootController := controller.NewRootController()
	rootController.RegisterRoutes(router)
	// Add routes here
	return router.Run(s.listenAddr)
}
