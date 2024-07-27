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

	r := gin.Default()

	// middlewares
	r.Use(middlewares.BearerTokenAuth())

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

	// all routes in one file for quick reference
	r.POST("/posts/images", imageHandler.GenerateSignedUrl)
	r.POST("/posts", postsHandler.CreatePost)
	r.GET("/posts", postsHandler.RetrievePostWithComments)
	r.POST("/posts/:postID/comments", commentsHandler.CreateComment)
	r.DELETE("/comments/:commentID", commentsHandler.DeleteComment)

	r.GET("/", handlers.Healthcheck)
	r.GET("/health", handlers.Healthcheck)

	err = r.Run(":" + config.Env().Port)
	if err != nil {
		log.Fatal("Failed to start server", err)
		return
	}
}
