package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ProjectID      uint      `json:"project_id"`
	AssignedTo     uint      `json:"assigned_to"`
	Status         string    `json:"status"`
	Priority       string    `json:"priority"`
	DueDate        time.Time `json:"due_date"`
	EstimatedHours float64   `json:"estimated_hours"`
	ActualHours    float64   `json:"actual_hours"`
	CreatedBy      uint      `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TaskCreateRequest struct {
	Title          string    `json:"title" binding:"required"`
	Description    string    `json:"description"`
	ProjectID      uint      `json:"project_id"`
	AssignedTo     uint      `json:"assigned_to"`
	Status         string    `json:"status"`
	Priority       string    `json:"priority"`
	DueDate        time.Time `json:"due_date"`
	EstimatedHours float64   `json:"estimated_hours"`
	ActualHours    float64   `json:"actual_hours"`
	CreatedBy      uint      `json:"created_by"`
}

var db *gorm.DB

func main() {
	initDB()

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		var dbStatus string
		if db != nil {
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					dbStatus = "OK"
				} else {
					dbStatus = "ERROR"
				}
			} else {
				dbStatus = "ERROR"
			}
		} else {
			dbStatus = "ERROR"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "Task Service OK",
			"database":  dbStatus,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Task routes
	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTask)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.GET("/tasks/stats", getTaskStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Task service running on port %s", port)
	log.Fatal(r.Run(":" + port))
}

func initDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	} else {
		log.Println("Successfully connected to database")
		db.AutoMigrate(&Task{})
	}
}

func getTasks(c *gin.Context) {
	var tasks []Task
	if db != nil {
		db.Find(&tasks)
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func getTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if db != nil {
		result := db.First(&task, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}
	}
	c.JSON(http.StatusOK, task)
}

func createTask(c *gin.Context) {
	var req TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	task := Task{
		Title:          req.Title,
		Description:    req.Description,
		ProjectID:      req.ProjectID,
		AssignedTo:     req.AssignedTo,
		Status:         req.Status,
		Priority:       req.Priority,
		DueDate:        req.DueDate,
		EstimatedHours: req.EstimatedHours,
		ActualHours:    req.ActualHours,
		CreatedBy:      req.CreatedBy,
	}

	if task.Status == "" {
		task.Status = "pending"
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	if db != nil {
		result := db.Create(&task)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания задачи: " + result.Error.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, task)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if db != nil {
		result := db.First(&task, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}

		var updateData TaskCreateRequest
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
			return
		}

		task.Title = updateData.Title
		task.Description = updateData.Description
		task.Status = updateData.Status
		task.Priority = updateData.Priority
		task.DueDate = updateData.DueDate
		task.EstimatedHours = updateData.EstimatedHours
		task.ActualHours = updateData.ActualHours
		task.UpdatedAt = time.Now()

		db.Save(&task)
	}

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	if db != nil {
		result := db.Delete(&Task{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления задачи"})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Задача успешно удалена"})
}

func checkTaskPermissions(c *gin.Context, task *Task, currentUserID uint, currentUserRole string) bool {
	if currentUserRole == "admin" {
		return true
	}
	if currentUserRole == "manager" {
		return task.AssignedTo == currentUserID || task.CreatedBy == currentUserID
	}
	return task.AssignedTo == currentUserID || task.CreatedBy == currentUserID
}

func getTaskStats(c *gin.Context) {
	var total int64
	if db != nil {
		db.Model(&Task{}).Count(&total)
	}
	c.JSON(http.StatusOK, gin.H{"total_tasks": total})
}
