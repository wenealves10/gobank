package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/wenealves10/gobank/db/sqlc"
	"github.com/wenealves10/gobank/token"
	"github.com/wenealves10/gobank/utils"
)

type Server struct {
	config       utils.Config
	store        db.Store
	tokenCreator token.TokenCreator
	router       *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	tokenCreator, err := token.NewPasetoTokenCreator(config.TokenPassetoKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token creator: %w", err)
	}

	server := &Server{
		store:        store,
		tokenCreator: tokenCreator,
		config:       config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenCreator))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
