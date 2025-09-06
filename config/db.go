package config

type DBConfig struct {
	DBPath string `env:"PATH" envDefault:"/var/lib/ipxe-manager/bolt.db"`
}
