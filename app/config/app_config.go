package config

type AppConfig struct {
	port   string
	dbURI  string
	dbName string
}

func GetNewAppConfig(port string, dbURI string, dbName string) *AppConfig {
	var newCfg AppConfig
	newCfg.dbURI = dbURI
	newCfg.dbName = dbName
	newCfg.port = port

	return &newCfg
}

func (c *AppConfig) GetDbURI() string {
	return c.dbURI
}

func (c *AppConfig) GetDbName() string {
	return c.dbName
}

func (c *AppConfig) GetPort() string {
	return c.port
}
