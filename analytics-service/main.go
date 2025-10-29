package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Analytics struct {
	TotalUsers         int64   `json:"total_users"`
	TotalTasks         int64   `json:"total_tasks"`
	TotalProjects      int64   `json:"total_projects"`
	CompletedTasks     int64   `json:"completed_tasks"`
	PendingTasks       int64   `json:"pending_tasks"`
	OverdueTasks       int64   `json:"overdue_tasks"`
	TaskCompletionRate float64 `json:"task_completion_rate"`
	AvgTaskHours       float64 `json:"avg_task_hours"`
	RecentActivities   int64   `json:"recent_activities"`
}

type ProjectStats struct {
	ProjectID      uint    `json:"project_id"`
	ProjectName    string  `json:"project_name"`
	TotalTasks     int64   `json:"total_tasks"`
	Completed      int64   `json:"completed"`
	Pending        int64   `json:"pending"`
	CompletionRate float64 `json:"completion_rate"`
}

var (
	db          *gorm.DB
	redisClient *redis.Client
)

func main() {
	initDB()
	initRedis()

	r := gin.Default()

	r.GET("/analytics/overview", getAnalyticsOverview)
	r.GET("/analytics/project-stats", getProjectStats)
	r.GET("/analytics/user-activity", getUserActivityStats)
	r.GET("/analytics/task-trends", getTaskTrends)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("Analytics service running on port %s", port)
	log.Fatal(r.Run(":" + port))
}

func initDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})
}

func getAnalyticsOverview(c *gin.Context) {
	var analytics Analytics

	// Базовые статистики
	db.Table("users").Count(&analytics.TotalUsers)
	db.Table("tasks").Count(&analytics.TotalTasks)
	db.Table("projects").Count(&analytics.TotalProjects)
	db.Table("tasks").Where("status = ?", "completed").Count(&analytics.CompletedTasks)
	db.Table("tasks").Where("status = ?", "pending").Count(&analytics.PendingTasks)
	db.Table("tasks").Where("due_date < ? AND status != ?", time.Now(), "completed").Count(&analytics.OverdueTasks)

	// Расчет процента завершения
	if analytics.TotalTasks > 0 {
		analytics.TaskCompletionRate = float64(analytics.CompletedTasks) / float64(analytics.TotalTasks) * 100
	}

	// Среднее время выполнения задач
	db.Table("tasks").Select("COALESCE(AVG(estimated_hours), 0)").Scan(&analytics.AvgTaskHours)

	// Активность за последние 7 дней
	since := time.Now().AddDate(0, 0, -7)
	db.Table("user_activities").Where("created_at >= ?", since).Count(&analytics.RecentActivities)

	c.JSON(http.StatusOK, analytics)
}

func getProjectStats(c *gin.Context) {
	var projectStats []ProjectStats

	rows, err := db.Table("projects p").
		Select("p.id as project_id, p.name as project_name, " +
			"COUNT(t.id) as total_tasks, " +
			"SUM(CASE WHEN t.status = 'completed' THEN 1 ELSE 0 END) as completed, " +
			"SUM(CASE WHEN t.status = 'pending' THEN 1 ELSE 0 END) as pending").
		Joins("LEFT JOIN tasks t ON t.project_id = p.id").
		Group("p.id, p.name").
		Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ps ProjectStats
		var total, completed int64

		rows.Scan(&ps.ProjectID, &ps.ProjectName, &total, &completed, &ps.Pending)

		ps.TotalTasks = total
		ps.Completed = completed

		if total > 0 {
			ps.CompletionRate = float64(completed) / float64(total) * 100
		}

		projectStats = append(projectStats, ps)
	}

	c.JSON(http.StatusOK, projectStats)
}

func getUserActivityStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	since := time.Now().AddDate(0, 0, -days)

	var result struct {
		DailyActivities []struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		} `json:"daily_activities"`
		TopUsers []struct {
			UserID   uint   `json:"user_id"`
			Username string `json:"username"`
			Count    int64  `json:"count"`
		} `json:"top_users"`
	}

	// Активность по дням
	db.Table("user_activities").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Order("date").
		Scan(&result.DailyActivities)

	// Топ пользователей по активности
	db.Table("user_activities ua").
		Select("ua.user_id, u.username, COUNT(*) as count").
		Joins("LEFT JOIN users u ON u.id = ua.user_id").
		Where("ua.created_at >= ?", since).
		Group("ua.user_id, u.username").
		Order("count DESC").
		Limit(10).
		Scan(&result.TopUsers)

	c.JSON(http.StatusOK, result)
}

func getTaskTrends(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	since := time.Now().AddDate(0, 0, -days)

	var trends struct {
		TasksByStatus []struct {
			Status string `json:"status"`
			Count  int64  `json:"count"`
		} `json:"tasks_by_status"`
		TasksByPriority []struct {
			Priority string `json:"priority"`
			Count    int64  `json:"count"`
		} `json:"tasks_by_priority"`
		WeeklyCompletion struct {
			Week      string `json:"week"`
			Created   int64  `json:"created"`
			Completed int64  `json:"completed"`
		} `json:"weekly_completion"`
	}

	// Задачи по статусам
	db.Table("tasks").
		Select("status, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("status").
		Scan(&trends.TasksByStatus)

	// Задачи по приоритету
	db.Table("tasks").
		Select("priority, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("priority").
		Scan(&trends.TasksByPriority)

	c.JSON(http.StatusOK, trends)
}
