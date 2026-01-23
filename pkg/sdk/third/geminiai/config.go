package geminiai

type Config struct {
	Proxy string
	Key   string
}

type FactoryConfig struct {
	Configs []Config
}
