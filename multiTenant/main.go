package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type JwtCustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	db        *gorm.DB
	jwtsecret = "gniognsgnspkpqpcmqc"
)

func initDB() {
	dsn := "host=localhost user=postgres password=test@123 dbname=tokendb port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to Database")
	}
	db.AutoMigrate(&User{})
}

func register(c echo.Context) error {
	type Input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var input Input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Wrong Request Format"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
    return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to hash password"})
	}

	user := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "DataBase Write Error"})
	}

	return c.JSON(http.StatusOK, user)
}

func login(c echo.Context) error {
	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input Input

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"Error": "Invalid Request Body"})
	}

	var u User
	if err := db.Where("email=?", input.Email).First(&u).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"Error": "Invalid Email"})
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
    return c.JSON(http.StatusUnauthorized, echo.Map{"Error": "Wrong Password"})
	}

	claims := JwtCustomClaims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtsecret))

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"JW Token": signed})
}

func adminDashboard(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	if claims.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Only accessible to Admin")
	}

	return c.JSON(http.StatusOK, echo.Map{"Greet": "Admin Dashboard", "userID": claims.UserID})
}

func tenantDashboard(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	if claims.Role != "tenant" {
		return echo.NewHTTPError(http.StatusForbidden, "Only accessible to Tenant")
	}

	return c.JSON(http.StatusOK, echo.Map{"Greet": "Tenant Dashboard", "userID": claims.UserID})
}

func userDashboard(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	if claims.Role != "user" {
		return echo.NewHTTPError(http.StatusForbidden, "Only accessible to User")
	}

	return c.JSON(http.StatusOK, echo.Map{"Greet": "User Dashboard", "userID": claims.UserID})
}

func main() {
	initDB()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/register", register)
	e.POST("/login", login)

	config := echojwt.Config{
		SigningKey: []byte(jwtsecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
	}

	r := e.Group("/api")
	r.Use(echojwt.WithConfig(config))
	r.POST("/admin", adminDashboard)
	r.POST("/tenant", tenantDashboard)
	r.POST("/user", userDashboard)

	e.Logger.Fatal(e.Start(":8080"))
}

// Alternative to fetch directly from DB instead of playing with JWT's ! 
