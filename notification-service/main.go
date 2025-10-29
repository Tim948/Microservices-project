package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Notification struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"user_id"`
	Title             string    `json:"title"`
	Message           string    `json:"message"`
	Type              string    `json:"type"`
	IsRead            bool      `json:"is_read"`
	RelatedEntityType string    `json:"related_entity_type"`
	RelatedEntityID   uint      `json:"related_entity_id"`
	CreatedAt         time.Time `json:"created_at"`
}

type UserActivity struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id"`
	ActivityType string    `json:"activity_type"`
	Description  string    `json:"description"`
	EntityType   string    `json:"entity_type"`
	EntityID     uint      `json:"entity_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
}

var db *gorm.DB

func main() {
	initDB()

	r := gin.Default()

	// Notification routes
	r.POST("/notifications", createNotification)
	r.GET("/notifications/user/:user_id", getUserNotifications)
	r.PUT("/notifications/:id/read", markAsRead)
	r.DELETE("/notifications/:id", deleteNotification)

	// Activity routes
	r.POST("/activities", logActivity)
	r.GET("/activities/user/:user_id", getUserActivities)
	r.GET("/activities/stats", getActivityStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Notification service running on port %s", port)
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

func createNotification(c *gin.Context) {
	var notification Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification.CreatedAt = time.Now()
	notification.IsRead = false

	result := db.Create(&notification)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

func getUserNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	unreadOnly := c.Query("unread_only") == "true"

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := db.Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	var notifications []Notification
	var total int64
	query.Model(&Notification{}).Count(&total)
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"page":          page,
		"limit":         limit,
		"total":         total,
	})
}

func markAsRead(c *gin.Context) {
	id := c.Param("id")

	result := db.Model(&Notification{}).Where("id = ?", id).Update("is_read", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func deleteNotification(c *gin.Context) {
	id := c.Param("id")

	result := db.Delete(&Notification{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}

func logActivity(c *gin.Context) {
	var activity UserActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity.CreatedAt = time.Now()
	activity.IPAddress = c.ClientIP()
	activity.UserAgent = c.GetHeader("User-Agent")

	result := db.Create(&activity)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func getUserActivities(c *gin.Context) {
	userID := c.Param("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var activities []UserActivity
	var total int64

	db.Where("user_id = ?", userID).Model(&UserActivity{}).Count(&total)
	result := db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&activities)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"page":       page,
		"limit":      limit,
		"total":      total,
	})
}

func getActivityStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	since := time.Now().AddDate(0, 0, -days)

	var stats struct {
		TotalActivities  int64            `json:"total_activities"`
		UniqueUsers      int64            `json:"unique_users"`
		ActivitiesByType map[string]int64 `json:"activities_by_type"`
	}

	db.Model(&UserActivity{}).Where("created_at >= ?", since).Count(&stats.TotalActivities)
	db.Model(&UserActivity{}).Where("created_at >= ?", since).Distinct("user_id").Count(&stats.UniqueUsers)

	// Activities by type
	var typeCounts []struct {
		ActivityType string
		Count        int64
	}
	db.Model(&UserActivity{}).Where("created_at >= ?", since).
		Select("activity_type, count(*) as count").
		Group("activity_type").Find(&typeCounts)

	stats.ActivitiesByType = make(map[string]int64)
	for _, tc := range typeCounts {
		stats.ActivitiesByType[tc.ActivityType] = tc.Count
	}

	c.JSON(http.StatusOK, stats)
}
