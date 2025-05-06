package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func ConnectToDB() {
	dsn := os.Getenv("DB_URL")

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	DB = bun.NewDB(sqldb, pgdialect.New())

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}
