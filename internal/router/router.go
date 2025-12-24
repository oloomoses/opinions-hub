package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oloomoses/opinions-hub/internal/database"
	"github.com/oloomoses/opinions-hub/internal/handlers"
	"github.com/oloomoses/opinions-hub/internal/repository"
)

func New() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)

	dbConn, err := database.Connect()

	if err != nil {
		log.Fatal("database connection failed")
	}

	r := gin.New()

	opinionRepo := repository.NewOpinionRepo(dbConn)
	opinionHandler := handlers.NewOpinionHandler(opinionRepo)

	r.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	r.GET("/health", handlers.Health)
	r.Static("/uploads", "./uploads")

	api := r.Group("api/v1")

	{
		api.GET("/opinions", opinionHandler.AllOpinions)
		api.POST("/opinion", opinionHandler.CreateOpinion)
		api.PATCH("/opinion/:id", opinionHandler.UpdateOpinion)
		api.DELETE("opinion/:id", opinionHandler.DeleteOpinion)

	}

	return r
}
