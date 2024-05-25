package bootstrap

import (
	"log"

	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDatabase(env *shared.Env) *gorm.DB {
	dsn := env.DSN

	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	shared.FailOnError(err, "Failed to connect to PostgreSQL")

	log.Println("Connection to PostgreSQL established.")

	return client
}

func AutoMigrate(client *gorm.DB) {
	err := client.AutoMigrate(&domain.User{}, &domain.Photo{})
	shared.FailOnError(err, "Failed to auto migrate the database")

	log.Println("Auto migration completed.")
}

func ClosePostgresDBConnection(client *gorm.DB) {
	sqlDB, err := client.DB()
	shared.FailOnError(err, "Failed to get database connection")

	err = sqlDB.Close()
	shared.FailOnError(err, "Failed to close the connection to PostgreSQL")

	log.Println("Connection to PostgreSQL closed.")
}
