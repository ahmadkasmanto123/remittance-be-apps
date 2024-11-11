package config

import (
	// "database/sql"

	"context"
	"os"

	"github.com/jackc/pgx/v5"
	// _ "github.com/lib/pq"
)

func GetDBRemit() (*pgx.Conn, error) {
	dbConfig := "host=" + os.Getenv("REMIT_HOST") + " user=" + os.Getenv("REMIT_USER") + " password=" + os.Getenv("REMIT_PASS") + " dbname=" + os.Getenv("REMIT_NAME") + " port=" + os.Getenv("REMIT_PORT") + " sslmode=disable TimeZone=Asia/Jakarta"
	db, err := pgx.Connect(context.Background(), dbConfig)
	return db, err
}

func GetDBLovePaycode() (*pgx.Conn, error) {
	dbConfig := "host=" + os.Getenv("PAYCODE_HOST") + " user=" + os.Getenv("PAYCODE_USER") + " password=" + os.Getenv("PAYCODE_PASS") + " dbname=" + os.Getenv("PAYCODE_NAME") + " port=" + os.Getenv("REMIT_PORT") + " sslmode=disable TimeZone=Asia/Jakarta"
	db, err := pgx.Connect(context.Background(), dbConfig)
	return db, err
}
