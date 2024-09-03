package config

const (
	DBHost     = "localhost"
	DBUser     = "offer_management_user"
	DBPassword = "offer_management_pass"
	DBName     = "offer_management_db"
	DBPort     = "5433"
)

func GetTestDBConfig() (string, string, string, string, string) {
	return DBHost, DBUser, DBPassword, DBName, DBPort
}
