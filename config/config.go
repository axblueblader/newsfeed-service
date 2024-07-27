package config

import "os"

type Config struct {
	Port                string
	PostImageBucketName string
	CommentsLimit       int
}

var c *Config

func Env() *Config {
	return c
}

func Init() {
	// load from files or environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	c = &Config{
		Port:                port,
		PostImageBucketName: "post-image-bucket",
		CommentsLimit:       2,
	}
}
