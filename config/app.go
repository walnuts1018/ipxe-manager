package config

type AppConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
}
