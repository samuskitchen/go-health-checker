package tools

// DbParams is a struct that contains the parameters to connect to the PostgresSQL database
type DbParams struct {
	Host           string
	Port           string
	User           string
	Password       string
	DbName         string
	SslMode        string
	MaxOpenCon     string
	MaxIdleCon     string
	MaxLifeTimeCon string
	MaxIdleTimeCon string
}
