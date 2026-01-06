package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `json:"name"`
	APIKey    string `gorm:"uniqueIndex;size;255" json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type RateLimit struct {
	Key string `gorm:"primaryKey;size:255"`
	Count int
	WindowStart time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func initDB() *gorm.DB {
	dsn:= "host=localhost user=postgres password=test@123 dbname=rate port=5432 sslmode=disable"
	db,err:= gorm.Open(postgres.Open(dsn))
	if err!=nil{
		log.Fatal("Failed to connect database ! ")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&RateLimit{})
	return db
}

func generateAPIkey() string {
	bytes := make([]byte,16)
	_,_ = rand.Read(bytes)

	return "API-"+hex.EncodeToString(bytes)
}

func Ratelimiter(db *gorm.DB,limit int,window time.Duration,exceedStatus int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-API-Key")
			if key==""{
				return c.JSON(http.StatusUnauthorized,echo.Map{
					"Error":"API key not available",
				})
			}

			now := time.Now().UTC()
			var rl RateLimit
			err:= db.First(&rl,"key=?",key).Error 
			if err!=nil && err!= gorm.ErrRecordNotFound {
				return c.JSON(http.StatusInternalServerError,echo.Map{
					"Error":"Database Error",
				})
			}

			if err==gorm.ErrRecordNotFound || now.Sub(rl.WindowStart) >= window {
				rl = RateLimit{
					Key: key,
					Count: 1,
					WindowStart: now,
				}
				db.Save(&rl)
				return next(c)
			}

			if rl.Count>=limit {
				retryAfter:= int(window.Seconds())-int(now.Sub(rl.WindowStart).Seconds())
				c.Response().Header().Set("Retry After: ",fmt.Sprint(retryAfter))
				return c.JSON(exceedStatus,echo.Map{
					"Error":"Rate Limit Exceeded",
					"Retry-After": retryAfter,
					"Status": exceedStatus,
				})
			}
			rl.Count++
			db.Save(&rl)
			return next(c)

		}
	}
}

func main() {
	db := initDB()
	e := echo.New()

	e.POST("/signup",func(c echo.Context) error {
		type SignupRequest struct {
			Name string `json:"name"`
		}
		req := new(SignupRequest)
		if err:=c.Bind(req);err!=nil || strings.TrimSpace(req.Name)==""{
			return c.JSON(http.StatusBadRequest,echo.Map{
				"Error":"Name is Required",
			})
		}
		apiKey := generateAPIkey()
		user := User {
			Name : req.Name,
			APIKey : apiKey,
		}
		db.Create(&user)
		return c.JSON(http.StatusOK,echo.Map{
			"Success":"User created Sucessfully",
			"APIKey":apiKey,
		})
	})
	
	rateLimit := Ratelimiter(db,15,15*time.Second,http.StatusForbidden)

	e.GET("/data",func(c echo.Context) error {
		return c.JSON(http.StatusOK,echo.Map{
			"Success":"Welcome to fetch point",
			"time":time.Now().Format(time.RFC3339),
		})
	},rateLimit)

	log.Print("Port started on 8080")
	e.Start(":8080")
}
