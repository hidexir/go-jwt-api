package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	fmt.Println(username, password)

	//// Throws unauthorized error
	//if username != "jon" || password != "shhh!" {
	//	return echo.ErrUnauthorized
	//}

	// Create token
	token := jwt.New(jwt.SigningMethodRS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_type"] = "tier1"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	keyData, _ := ioutil.ReadFile("jwtRS256.key")
	fmt.Println(string(keyData))
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)

	if err != nil {
		fmt.Println(err, "43")
		return err
	}

	t, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err, "49")
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	keyData, _ := ioutil.ReadFile("sample_key.pub")
	key, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: "RS256",
	}))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":1323"))
}
