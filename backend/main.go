
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

const API_KEY = "3b43d97459a5d614321eb5f80f4978b0"

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

func getWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City is required"})
		return
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach weather API"})
		return
	}
	defer resp.Body.Close()

	var weatherData WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode API response"})
		return
	}

	if weatherData.Cod != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": weatherData.Name})
		return
	}

	weatherCondition := "No data"
	if len(weatherData.Weather) > 0 {
		weatherCondition = weatherData.Weather[0].Description
	}

	c.JSON(http.StatusOK, gin.H{
		"city":      weatherData.Name,
		"temp":      weatherData.Main.Temp,
		"condition": weatherCondition,
		"units":     "Celsius",
	})
}

func main() {
	r := gin.Default()

	// ✅ Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Allow all origins (change to frontend URL in production)
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.GET("/weather", getWeather)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Println("Server running on port:", port)
	r.Run(":" + port)
}
