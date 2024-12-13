package config

// ServerConfig the server config
type ServerConfig struct {
	Name string
}

// AppConfig the app config
type AppConfig struct {
	Token   string
	Columns int32
}

// Config the config struct
type Config struct {
	Server ServerConfig
	App    AppConfig
}

// Module the module
type Module struct {
	cfg *Config
}
