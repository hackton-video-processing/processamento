package config

const (
	User     = "MYSQL_USER"
	Password = "PASSWORD"
	Endpoint = "ENDPOINT"
	Port     = "MYSQL_PORT"
	DBName   = "DBNAME"

	_defaultUser     = "user"
	_defaultPassword = "password"
	_defaultEndpoint = "endpoit"
	_defaultPort     = "port"
	_defaultDBName   = "db_name"
)

type MySQLConfig struct {
	User     string
	Password string
	Endpoint string
	Port     string
	DBName   string
}

func NewMySQLConfig() MySQLConfig {
	return MySQLConfig{
		User:     GetString(User, _defaultUser),
		Password: GetString(Password, _defaultPassword),
		Endpoint: GetString(Endpoint, _defaultEndpoint),
		Port:     GetString(Port, _defaultPort),
		DBName:   GetString(DBName, _defaultDBName),
	}
}
