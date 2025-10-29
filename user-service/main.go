package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserCreateRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

var (
	db          *gorm.DB
	redisClient *redis.Client
)

func main() {
	initDB()
	initRedis()

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		_, redisErr := redisClient.Ping(c).Result()
		redisStatus := "OK"
		if redisErr != nil {
			redisStatus = "ERROR"
		}

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
			"status":    "User Service OK",
			"database":  dbStatus,
			"redis":     redisStatus,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// User routes
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUser)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
	r.GET("/users/stats", getUserStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("User service running on port %s", port)
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
		db.AutoMigrate(&User{})
	}
}

func initRedis() {
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
}

func getUsers(c *gin.Context) {
	var users []User
	if db != nil {
		db.Find(&users)
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if db != nil {
		result := db.First(&user, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
	}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	user := User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if db != nil {
		result := db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя: " + result.Error.Error()})
			return
		}
		// Invalidate cache
		redisClient.Del(c, "users:all")
	}

	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if db != nil {
		result := db.First(&user, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		var updateData UserCreateRequest
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
			return
		}

		user.Username = updateData.Username
		user.Email = updateData.Email
		user.FirstName = updateData.FirstName
		user.LastName = updateData.LastName
		user.Role = updateData.Role
		user.UpdatedAt = time.Now()

		db.Save(&user)
		// Invalidate cache
		redisClient.Del(c, "users:all")
	}

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	if db != nil {
		result := db.Delete(&User{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления пользователя"})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		// Invalidate cache
		redisClient.Del(c, "users:all")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}

func getUserStats(c *gin.Context) {
	var total int64
	if db != nil {
		db.Model(&User{}).Count(&total)
	}
	c.JSON(http.StatusOK, gin.H{"total_users": total})
}
