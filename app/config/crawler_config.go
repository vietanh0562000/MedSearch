package config

type CrawlerConfig struct {
	baseURL  string
	startURL string
	dbURI    string
	dbName   string
}

func GetNewCrawlerConfig(baseURL string, startURL string, dbURI string, dbName string) *CrawlerConfig {
	var newCfg CrawlerConfig
	newCfg.baseURL = baseURL
	newCfg.startURL = startURL
	newCfg.dbURI = dbURI
	newCfg.dbName = dbName

	return &newCfg
}

func (c *CrawlerConfig) GetBaseURL() string {
	return c.baseURL
}

func (c *CrawlerConfig) GetStartURL() string {
	return c.startURL
}

func (c *CrawlerConfig) GetDbURI() string {
	return c.dbURI
}

func (c *CrawlerConfig) GetDbName() string {
	return c.dbName
}
