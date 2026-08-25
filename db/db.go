package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {

	connectionString := os.Getenv("DATABASE_URL")
	

	db, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Database connected successfully")

	return db, nil
}