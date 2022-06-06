package config

type Config struct {
	Servers map[string]Server `json:"servers"`
}

type Server struct {
	Address  string   `json:"address,omitempty"`
	Database Database `json:"database"`
}
type Database struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
	Host     string `json:"host,omitempty"`
	Database string `json:"database,omitempty"`
}
