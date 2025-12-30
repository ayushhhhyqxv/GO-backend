package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

var jwtSecret = []byte("dhfioAodAJIonKnmkKM")

type User struct {
	ID int `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Password string `json:"-"` // Hide the password in response! 
}

func main(){
	dsn:="host=localhost user=postgres password=test@123 dbname=testdb port=5432 sslmode=disable"

	database,err:= gorm.Open(postgres.Open(dsn),&gorm.Config{}) // gorm config allows to customize gorm how it behaves ! 

	if err!=nil {
		panic("Failed to connect to Database ! ")
	}

	db = database
	db.AutoMigrate(&User{})

	e:= echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.POST("/register",register)
	e.POST("/login",login)
	r := e.Group("/user")
	r.Use(authMiddleware)
	r.GET("/profile",profile)
	e.Logger.Fatal(e.Start(":8000"))

}

func register(c echo.Context) error {
	u:= new(User)
	if err:= c.Bind(u);err!=nil{
		return err
	}
	hash,err:=bcrypt.GenerateFromPassword([]byte(u.Password),bcrypt.DefaultCost)
	if err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error": "Failed to Hash Password "})
	}
	u.Password = string(hash)

	if err:=db.Create(&u).Error;err!=nil{
		return c.JSON(http.StatusBadRequest,echo.Map{"Error":"User Already Exists"})
	}

	return c.JSON(http.StatusOK,echo.Map{"Success":"User Registered Successfully"})
}

func login(c echo.Context) error {
	req:= new(User)

	if err:=c.Bind(req);err!=nil{
		return err
	}
	var check User 
	if err:=db.Where("email=?",req.Email).First(&check).Error;err!=nil{
		return c.JSON(http.StatusUnauthorized,echo.Map{"Error":"Invalid Email"})
	}
	// if check.Password != req.Password { 
    // return c.JSON(http.StatusUnauthorized, echo.Map{"Error": "Invalid Password"})
// }
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id":check.ID,
		"email":check.Email,
		"exp":time.Now().Add(time.Hour).Unix(),
	})

	t,err:=token.SignedString(jwtSecret)
	if err!=nil{
		return err
	}
	return c.JSON(http.StatusOK,echo.Map{
		"message":"Login Successful",
		"token": t ,
	})
}

func profile(c echo.Context) error {
	userID := c.Get("user_id")
	var check User 

	if err:=db.First(&check,userID).Error;err!=nil{
		return c.JSON(http.StatusNotFound,echo.Map{"Error":"User Not Found"})
	}
	return c.JSON(http.StatusOK,check)
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader:=c.Request().Header.Get("Authorization")
		if authHeader==""{
			return c.JSON(http.StatusUnauthorized,echo.Map{"Error":"Missing Token"})
		}
		tokenString := ""
		fmt.Sscanf(authHeader,"Bearer %s",&tokenString)

		token,err:=jwt.Parse(tokenString,func(token *jwt.Token)(interface{},error){
			if _,ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
				return nil,fmt.Errorf("Unexpected Signing Method")
			}
			return jwtSecret,nil
		})
		if err!=nil || !token.Valid{
			return c.JSON(http.StatusUnauthorized,echo.Map{"Error":"Invalid Token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Extract the user_id you put into the token during login
            userID := claims["user_id"]
            // "Set" it into the context so the profile function can "Get" it
            c.Set("user_id", userID)
        }
		
		return next(c)
	}
}
