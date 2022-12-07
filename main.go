package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	// Echo instance
	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/discord", discord)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Exotico")
	})
	// Start server
	log.Println("http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func discord(c echo.Context) error {
	var jsonData = []byte(`{
		"username": "Agogo",
		"avatar_url": "https://images-ext-2.discordapp.net/external/od190cVs9THpP2sqpSZMiV-lQK5XuzN47uSGvsgkEfI/%3Fsize%3D4096%26ignore%3Dtrue/https/cdn.discordapp.com/avatars/665223957082931210/a6bcf1655eccbc2c7b43d75851abf725.png",
		"content":"exotico"
	}`)
	req, err := http.NewRequest("POST", goDotEnvVariable("URL"), bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error mi fafa")
	}

	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		return c.String(http.StatusInternalServerError, "Error mí fafa")
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	return c.String(http.StatusOK, "Lo mando mí fafa")
}
