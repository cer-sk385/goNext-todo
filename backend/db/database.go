package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB データベース接続を初期化
func InitDB() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	// データベース接続を最大30秒間再試行
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("データベース接続エラー（試行 %d/%d）: %v", i+1, maxRetries, err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = DB.Ping()
		if err == nil {
			log.Println("データベースに接続しました")
			return
		}

		log.Printf("データベースPingエラー（試行 %d/%d）: %v", i+1, maxRetries, err)
		time.Sleep(1 * time.Second)
	}

	log.Fatalf("データベース接続の最大試行回数を超えました。最後のエラー: %v", err)
}
