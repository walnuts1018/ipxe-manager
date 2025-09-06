package config

type IpxeConfig struct {
	ScriptDir string `env:"SCRIPT_DIR" envDefault:"/etc/ipxe"`
	DefaultOS string `env:"DEFAULT_OS" envDefault:"windows"`
}
