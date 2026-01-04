package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oloomoses/opinions-hub/internal/database"
	"github.com/oloomoses/opinions-hub/internal/handlers"
	"github.com/oloomoses/opinions-hub/internal/middleware"
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

	userRepo := repository.NewUserRepo(dbConn)
	userHandler := handlers.NewUserHandler(userRepo)

	r.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	r.GET("/health", handlers.Health)
	r.Static("/uploads", "./uploads")

	protected := r.Group("api/v1")
	protected.Use(middleware.LoginRequired())

	{
		protected.GET("/profile", userHandler.GetUserProfile)
		protected.GET("/opinions", opinionHandler.AllOpinions)
		protected.POST("/opinion", opinionHandler.CreateOpinion)
		protected.PATCH("/opinion/:id", opinionHandler.UpdateOpinion)
		protected.DELETE("opinion/:id", opinionHandler.DeleteOpinion)
		protected.POST("opinion/:opinion_id", opinionHandler.CreateOpinion)

	}

	protected.POST("/users/:id/follow", userHandler.FollowUser)
	protected.DELETE("/users/:id/follow", userHandler.Unfollow)
	protected.GET("/users/:id/following", userHandler.GetFollowing)
	protected.GET("/users/:id/followers", userHandler.GetFollowers)

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.LoginUser)
	return r
}
