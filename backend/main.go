package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"backend/db"
	"backend/handlers"
)

func main() {
	// データベース接続を初期化
	db.InitDB()

	// プログラム終了時にデータベース接続を閉じる
	defer db.DB.Close()

	// rはginのインスタンス
	// gin.Default()はデフォルトの設定を使用してginのインスタンスを作成
	r := gin.Default()

	// CORSミドルウェアを追加
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// ヘルスチェック
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// TODOのCRUDエンドポイント
	r.GET("/api/todos", handlers.GetAllTodos)
	r.GET("/api/todos/:id", handlers.GetTodo)
	r.POST("/api/todos", handlers.CreateTodo)
	r.PUT("/api/todos/:id", handlers.UpdateTodo)
	r.DELETE("/api/todos/:id", handlers.DeleteTodo)

	r.Run(":8080")
}
