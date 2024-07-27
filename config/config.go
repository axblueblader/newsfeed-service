package config

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
	c = &Config{
		Port:                "8080",
		PostImageBucketName: "post-image-bucket",
		CommentsLimit:       2,
	}
}
