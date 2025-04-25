package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"backend/db"
	"backend/models"
)

// GetAllTodos すべてのTODOを取得
func GetAllTodos(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY id")
	if err != nil {
		log.Printf("クエリエラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの取得に失敗しました"})
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Printf("スキャンエラー: %v", err)
			continue
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodo 特定のTODOを取得
func GetTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	err := db.DB.QueryRow("SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1", id).
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

// CreateTodo 新しいTODOを作成
func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストボディ"})
		return
	}

	if todo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "タイトルが必要です"})
		return
	}

	err := db.DB.QueryRow(`
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

// UpdateTodo TODOを更新
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストボディ"})
		return
	}

	// 更新時間を現在時刻に設定
	now := time.Now()

	// TODOが存在するか確認
	var existingID int
	err := db.DB.QueryRow("SELECT id FROM todos WHERE id = $1", id).Scan(&existingID)
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
	_, err = db.DB.Exec(`
		UPDATE todos SET title = $1, description = $2, completed = $3, updated_at = $4 
		WHERE id = $5`,
		todo.Title, todo.Description, todo.Completed, now, id)

	if err != nil {
		log.Printf("更新エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TODOの更新に失敗しました"})
		return
	}

	// 更新後のTODOを取得
	err = db.DB.QueryRow(`
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

// DeleteTodo TODOを削除
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	result, err := db.DB.Exec("DELETE FROM todos WHERE id = $1", id)
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
