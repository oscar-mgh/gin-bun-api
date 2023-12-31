package db

import (
	// "context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	// "github.com/oscar-mgh/gin_bun/model"
	// "github.com/uptrace/bun"
	// "github.com/uptrace/bun/dialect/mysqldialect"
)

func Connect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DB_NAME")
	)

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	sqldb, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		panic(err)
	}

	// db := bun.NewDB(sqldb, mysqldialect.New())
	// err = db.ResetModel(context.TODO(), &model.GenreModel{})
	// if err != nil {
	// 	panic(err)
	// }
	// return db.DB

	return sqldb
}
