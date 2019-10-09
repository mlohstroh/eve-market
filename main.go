package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	profileFunction("loadSDE", func() {
		err := loadSDE()
		if err != nil {
			panic(err)
		}
	})

	profileFunction("getPrices", func() {
		err := getPrices()
		if err != nil {
			fmt.Print(err)
		}
	})

	scheduler := NewScheduler(1 * time.Hour)
	scheduler.Schedule("FetchPrices", getPrices, 5*time.Minute)
	go scheduler.Run()

	router := gin.Default()
	router.GET("/items/", getItems)
	router.GET("/items/:id", getItem)

	router.Run(":3000")
}
