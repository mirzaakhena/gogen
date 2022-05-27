package config

type Config struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
}

type Server struct {
	Port int `json:"port,omitempty"`
}

type Database struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
	Host     string `json:"host,omitempty"`
	Database string `json:"database,omitempty"`
}
