package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/chneau/limiter"

	"github.com/labstack/echo"
)

type ParamDto struct {
	Count int `json:"count" query:"count"`
}

func main() {
	goRoutine := flag.String("GO_ROUTINE", os.Getenv("GO_ROUTINE"), "GO_ROUTINE")
	goRoutineCount, _ := strconv.ParseInt(*goRoutine, 10, 64)
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	var limitEvent limiter.Limiter = limiter.New(int(goRoutineCount))
	e.GET("/test1", func(c echo.Context) error {

		limitEvent.Execute(func() {
			var v ParamDto
			if err := c.Bind(&v); err != nil {
				return
			}
			i := 0
			for index := 0; index < v.Count; index++ {
				for index2 := 0; index2 < 100000000; index2++ {
					i += index2
				}
			}
			fmt.Println(i)
		})
		return c.String(http.StatusOK, "success")
	})

	e.GET("/test2", func(c echo.Context) error {
		var v ParamDto
		if err := c.Bind(&v); err != nil {
			return c.String(http.StatusOK, "failure")
		}
		i := 0
		for index := 0; index < v.Count; index++ {
			for index2 := 0; index2 < 100000000; index2++ {
				i += index2
			}
		}
		fmt.Println(i)
		return c.String(http.StatusOK, "success")
	})

	e.GET("/test3", func(c echo.Context) error {
		go func(cA echo.Context) {
			var v ParamDto
			if err := cA.Bind(&v); err != nil {
				return
			}
			i := 0
			for index := 0; index < v.Count; index++ {
				for index2 := 0; index2 < 100000000; index2++ {
					i += index2
				}
			}
			fmt.Println(i)
		}(c)
		return c.String(http.StatusOK, "success")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
