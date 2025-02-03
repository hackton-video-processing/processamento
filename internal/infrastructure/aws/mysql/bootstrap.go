package mysql

import (
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func BootstrapMySQLRepository(config config.AppConfig) (*Repository, error) {
	if config.Env.IsLocal() {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err := db.AutoMigrate(&ProcessMySQL{}, &File{}); err != nil {
			return nil, fmt.Errorf("erro ao rodar AutoMigrate: %w", err)
		}

		return NewMySQLRepository(db), nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Endpoint,
		config.MySQL.Port,
		config.MySQL.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Erro ao obter a instância de *sql.DB:", err)
		return nil, fmt.Errorf("erro ao obter a instância de *sql.DB: %s", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar ping no banco de dados: %s", err)
	}
	defer sqlDB.Close()

	return NewMySQLRepository(db), nil
}
