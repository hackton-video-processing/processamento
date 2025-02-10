package config

const (
	User     = "MYSQL_USER"
	Password = "PASSWORD"
	Endpoint = "ENDPOINT"
	Port     = "MYSQL_PORT"
	DBName   = "DBNAME"

	_defaultUser     = "admin"
	_defaultPassword = "XPK39b73"
	_defaultEndpoint = "video-processing-api-database.chcgmmmie0nu.us-east-1.rds.amazonaws.com"
	_defaultPort     = "3306"
	_defaultDBName   = "video_processing_api_database"
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
