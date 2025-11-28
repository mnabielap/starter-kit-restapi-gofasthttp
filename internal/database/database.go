package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/config"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB establishes a connection to the database based on configuration
func ConnectDB(cfg *config.Config) *gorm.DB {
	var dialector gorm.Dialector

	switch cfg.DBDriver {
	case "postgres":
		dialector = postgres.Open(cfg.DBSource)
	case "sqlite":
		dialector = sqlite.Open(cfg.DBSource)
	default:
		log.Fatalf("❌ Unsupported DB_DRIVER: %s", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Printf("✅ Connected to Database (%s)", cfg.DBDriver)

	// AutoMigrate creates tables based on structs
	log.Println("🔄 Running Auto Migration...")
	err = db.AutoMigrate(&model.User{}, &model.Token{})
	if err != nil {
		log.Fatalf("❌ Auto Migration failed: %v", err)
	}
	log.Println("✅ Auto Migration completed.")

	return db
}