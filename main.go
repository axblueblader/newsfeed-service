package main

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"newsfeed-service/config"
	"newsfeed-service/handlers"
	"newsfeed-service/middlewares"
	"newsfeed-service/models"
	"newsfeed-service/services"
	"newsfeed-service/storage"
)

func main() {
	// load some environment variables here
	config.Init()

	// init tracing agent
	// tracer.init()

	postgresDB := embeddedpostgres.NewDatabase()
	err := postgresDB.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(postgres *embeddedpostgres.EmbeddedPostgres) {
		err := postgres.Stop()
		if err != nil {
			log.Println(err)
		}
	}(postgresDB)

	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"))

	err = db.AutoMigrate(models.Post{}, models.Comment{})
	if err != nil {
		log.Fatal("Failed to auto migrate schemas", err)
		return
	}

	// database access
	postsDB := storage.NewPostDB(db)
	commentsDB := storage.NewCommentDB(db)

	// external services
	objectStorage := storage.NewObjectStorage("i am an aws client")

	// internal services
	postService := services.NewPostService(postsDB)
	commentService := services.NewCommentService(commentsDB)

	// handlers
	commentsHandler := handlers.CommentsHandler{
		CommentService: commentService,
	}
	postsHandler := handlers.PostsHandler{
		PostService: postService,
	}
	imageHandler := handlers.ImagesHandler{
		ObjectStorage: objectStorage,
	}

	r := gin.Default()

	// health check and non auth public endpoints
	r.GET("/", handlers.Healthcheck)
	r.GET("/health", handlers.Healthcheck)

	// all routes in one file for quick reference
	authed := r.Group("/")
	authed.Use(middlewares.BearerTokenAuth())
	{
		authed.POST("/posts/images", imageHandler.GenerateSignedUrl)
		authed.POST("/posts", postsHandler.CreatePost)
		authed.GET("/posts", postsHandler.RetrievePostWithComments)
		authed.POST("/posts/:postID/comments", commentsHandler.CreateComment)
		authed.DELETE("/comments/:commentID", commentsHandler.DeleteComment)
	}

	err = r.Run(":" + config.Env().Port)
	if err != nil {
		log.Fatal("Failed to start server", err)
		return
	}
}
