package bootstrap

import (
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDatabase(env *shared.Env) *gorm.DB {
	dsn := env.DSN

	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	failOnError(err, "Failed to connect to PostgreSQL")

	log.Println("Connection to PostgreSQL established.")

	return client
}

func AutoMigrate(client *gorm.DB) {
	err := client.AutoMigrate(&domain.User{})
	failOnError(err, "Failed to auto migrate the database")

	log.Println("Auto migration completed.")
}

func ClosePostgresDBConnection(client *gorm.DB) {
	sqlDB, err := client.DB()
	failOnError(err, "Failed to get database connection")

	err = sqlDB.Close()
	failOnError(err, "Failed to close the connection to PostgreSQL")

	log.Println("Connection to PostgreSQL closed.")
}
