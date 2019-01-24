package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/chneau/limiter"

	"github.com/labstack/echo"
)

func main() {
	goRoutine := flag.String("GO_ROUTINE", os.Getenv("GO_ROUTINE"), "GO_ROUTINE")
	taskCount := flag.String("TASK_COUNT", os.Getenv("TASK_COUNT"), "TASK_COUNT")

	goRoutineCount, _ := strconv.ParseInt(*goRoutine, 10, 64)
	count, _ := strconv.ParseInt(*taskCount, 10, 64)
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	var limitEvent limiter.Limiter = limiter.New(int(goRoutineCount))
	e.GET("/test1", func(c echo.Context) error {
		limitEvent.Execute(func() {
			for index := 0; index < int(count); index++ {
				time.Sleep(20 * time.Second)
			}
		})
		return c.String(http.StatusOK, "success")
	})

	e.GET("/test2", func(c echo.Context) error {
		for index := 0; index < int(count); index++ {
			time.Sleep(20 * time.Second)
		}
		return c.String(http.StatusOK, "success")
	})

	e.GET("/test3", func(c echo.Context) error {
		go func(number int64) {
			for index := 0; index < int(number); index++ {
				time.Sleep(20 * time.Second)
			}
		}(count)
		return c.String(http.StatusOK, "success")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
