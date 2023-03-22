package config

import (
	"time"
)

var DefaultConfig = map[string]map[string]interface{}{
	// "rabbit_queues": RabbitMQDefaults,
	// "drs":      DefaultDRS,
}

func (c *Config) loadEnv() {
	c.MongoDB.DatabaseURL = c.envVariables("MONGODB_URL")
}

func InitConfig() Config {
	c := Config{}

	logger := &Logging{LogLevel: "debug"}

	mongo := &MongoDB{
		DatabaseURL:        "mongodb://dbtest:supra**@localhost:27017/?authSource=admin",
		MaxPoolSize:        20,
		ConnectTimeout:     30 * time.Second,
		DatabaseName:       "developmet_db",
		ThoughtsCollection: "thoughts_collection",
	}

	c.MongoDB = mongo
	c.Logging = logger

	return c
}
