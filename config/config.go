package config

import "fmt"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func GetDbConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "TokoBelanja",
	}
}

func (config *DBConfig) GetDBURL() string {
	return fmt.Sprintf(
		"user=postgres password=p6H2HFVhOSggar5T host=db.unahekodberrlcmdukvl.supabase.co port=5432 dbname=postgres",
	)
}

// "host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
// config.Host, config.Port, config.User, config.DBName, config.Password,
