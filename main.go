package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello world")
	})

	if err := e.Start("127.0.0.1:1358"); err != nil {
		log.Fatalf("cannot start the echo http server: %s", err)
	}

}
