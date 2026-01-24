package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)

type User struct {
	ID uint `gorm:"primaryKey:autoIncrement"`
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

func initDB() *gorm.DB{
	dsn:= os.Getenv("POSTGRES_DSN")
	db,err:= gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatalf("Failed to connect to Database: %v",err)
	}
	db.AutoMigrate(&User{})
	return db
}

func initRedis() *redis.Client {
	rdb:= redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})
	_,err := rdb.Ping(ctx).Result()
	if err!=nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return rdb;
}

func getinfo(c echo.Context) error{
	// val,err:= rdb.Get(ctx,"users").Result()
	// if err==redis.Nil {
	// 	var users []User
	// 	if err:=db.Find(&users).Error;err!=nil{
	// 		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Data Not Found"})
	// 	}

	// 	data,_ := json.Marshal(users)
	// 	rdb.Set(ctx,"users",data,10*time.Minute)
	// 	return c.JSON(http.StatusOK,echo.Map{"Success":"Data saved to redis"})
	// }else if err!=nil {
	// 	return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Data Not Found"})
	// }

	// var users []User
	// if err:= json.Unmarshal([]byte(val),&users);err!=nil{
	// 	return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Failed to Parse"})
	// }
	// return c.JSON(http.StatusOK,users)

	var users []User
	if err:=db.Find(&users).Error;err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Success":"DB error"})
	}

	data,_ := json.Marshal(users)
		rdb.Set(ctx,"users",data,10*time.Minute)
		return c.JSON(http.StatusOK,users)

}

func createUser(c echo.Context) error{
	u:= new(User)
	if err:= c.Bind(u);err!=nil{
		return c.JSON(http.StatusBadGateway,echo.Map{"Error":"Failed to Bind Data"}) 
	}
	var lastUser User

	if err:=db.Order("id desc").First(&lastUser).Error;err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Database Fetching Error"}) 
	}

	u.ID = lastUser.ID + 1

	if err:=db.Create(&u).Error;err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Database Insertion Error"})  
	}

	var users []User
	if err:=db.Find(&users).Error;err==nil{
		data,_:=json.Marshal(users)
		rdb.Set(ctx,"users",data,10*time.Minute) // Impossible to append cache !
	}

	data,_:=json.Marshal(u)
	rdb.Set(ctx,fmt.Sprintf("user:%d",u.ID),data,10*time.Minute)
	return c.JSON(http.StatusCreated,u)
}

func main() {

	if err := godotenv.Load();err!=nil{
		log.Printf("No Prior Credentials Passed !")
	}
	db = initDB()
	rdb = initRedis()
	e:=echo.New()
	e.GET("/users",getinfo)
	e.POST("/create",createUser)

	e.Logger.Fatal(e.Start(":8080"))	


}