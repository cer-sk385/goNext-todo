package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var db *sql.DB

func initDB() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
	
	// データベース接続を最大30秒間再試行
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("データベース接続エラー（試行 %d/%d）: %v", i+1, maxRetries, err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("データベースに接続しました")
			return
		}
		
		log.Printf("データベースPingエラー（試行 %d/%d）: %v", i+1, maxRetries, err)
		time.Sleep(1 * time.Second)
	}
	
	log.Fatalf("データベース接続の最大試行回数を超えました。最後のエラー: %v", err)
}

func main() {
	initDB()
	defer db.Close()

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
	r.GET("/api/todos", getAllTodos)
	r.GET("/api/todos/:id", getTodo)
	r.POST("/api/todos", createTodo)
	r.PUT("/api/todos/:id", updateTodo)
	r.DELETE("/api/todos/:id", deleteTodo)

	r.Run(":8080")
}

// すべてのTODOを取得
func getAllTodos(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY id")
	if err != nil {
		log.Printf("クエリエラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの取得に失敗しました"})
		return
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Printf("スキャンエラー: %v", err)
			continue
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

// 特定のTODOを取得
func getTodo(c *gin.Context) {
	id := c.Param("id")
	var todo Todo

	err := db.QueryRow("SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1", id).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "TODOが見つかりません"})
		} else {
			log.Printf("クエリエラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの取得に失敗しました"})
		}
		return
	}

	c.JSON(http.StatusOK, todo)
}

// 新しいTODOを作成
func createTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストボディ"})
		return
	}

	if todo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "タイトルが必要です"})
		return
	}

	err := db.QueryRow(`
		INSERT INTO todos (title, description, completed) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`,
		todo.Title, todo.Description, todo.Completed).
		Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		log.Printf("挿入エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// TODOを更新
func updateTodo(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストボディ"})
		return
	}

	// 更新時間を現在時刻に設定
	now := time.Now()

	// TODOが存在するか確認
	var existingID int
	err := db.QueryRow("SELECT id FROM todos WHERE id = $1", id).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "TODOが見つかりません"})
		} else {
			log.Printf("クエリエラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの更新に失敗しました"})
		}
		return
	}

	// TODOを更新
	_, err = db.Exec(`
		UPDATE todos SET title = $1, description = $2, completed = $3, updated_at = $4 
		WHERE id = $5`,
		todo.Title, todo.Description, todo.Completed, now, id)

	if err != nil {
		log.Printf("更新エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの更新に失敗しました"})
		return
	}

	// 更新後のTODOを取得
	err = db.QueryRow(`
		SELECT id, title, description, completed, created_at, updated_at 
		FROM todos WHERE id = $1`, id).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		log.Printf("クエリエラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新されたTODOの取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// TODOを削除
func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		log.Printf("削除エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの削除に失敗しました"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("行数取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "削除の確認に失敗しました"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "TODOが見つかりません"})
		return
	}

	c.Status(http.StatusNoContent)
}
