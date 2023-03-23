package config

import (
	"errors"
	"io/fs"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	API         *API     `mapstructure:"api"`
	Logging     *Logging `mapstructure:"logging"`
	MongoDB     *MongoDB `mapstructure:"mongo"`
	Grumpy      *Grumpy  `mapstructure:"grumplyclient"`
}

type Grumpy struct {
}

type API struct {
	Address string `mapstructure:"address"`
}

type MongoDB struct {
	DatabaseURL        string        `mapstructure:"database_url"`
	MaxPoolSize        int           `mapstructure:"max_pool_size"`
	ConnectTimeout     time.Duration `mapstructure:"connect_timeout"`
	DatabaseName       string        `mapstructure:"database_name"`
	ThoughtsCollection string        `mapstructure:"thoughts"`
	TopicsCollection   string        `mapstructure:"topics"`
}

type Logging struct {
	LogLevel string `mapstructure:"loglevel"`
}

type Worker struct {
	BunchSize      int           `mapstructure:"bunch_size"`
	TickerInterval time.Duration `mapstructure:"ticker_interval"`
	TaskTimeout    time.Duration `mapstructure:"task_timeout"`
	TaskConnDelay  time.Duration `mapstructure:"task_conn_delay"`
}

func loadDefaults(v *viper.Viper) {
	for key, val := range DefaultConfig {
		v.SetDefault(key, val)
	}
}

func LoadConf(env string) (Config, error) {
	var c Config

	var confFileName string
	if env == "prod" {
		confFileName = "production"
		err := c.readEnvironment(".env")
		if err != nil {
			return Config{}, err
		}
	} else {
		confFileName = "development"
		err := c.readEnvironment("../.env")
		if err != nil {
			return Config{}, err
		}
	}

	v := viper.New()
	v.SetConfigName(confFileName)
	v.AddConfigPath("../")
	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	loadDefaults(v)

	if err := v.Unmarshal(&c); err != nil {
		return Config{}, err
	}

	// Setup environment. Used for adjusting logic for different instances(dev & prod)
	c.Environment = env
	c.loadEnv()
	return c, nil
}

// envVariables assigns variables from Environment
func (c *Config) envVariables(key string) string {
	return os.Getenv(key)
}

// readEnvironment reads the first existing env file from the list
func (c *Config) readEnvironment(files ...string) error {
	for _, f := range files {
		err := godotenv.Load(f)
		if err == nil {
			return nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}
