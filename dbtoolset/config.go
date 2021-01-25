package dbtoolset

type DBConfig struct {
	Redis map[string]RedisConfig
	MySQL map[string]MySQLConfig
}

type RedisConfig struct {
	DSN string `json:"dsn"`
}

type MySQLConfig struct {
	DSN     string `json:"dsn"`
	ShowSQL bool   `json:"show_sql"`
}
