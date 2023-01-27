package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hamidreza-abooei/ie-project/db"
	"github.com/hamidreza-abooei/ie-project/monitor"
	"github.com/labstack/echo/v4"
)

func main() {
	// Setup Database
	d := db.Setup("ie-project.db")
	st := db.NewStore(d)
	mnt := monitor.NewMonitor(st, nil, 10)
	sch, _ := monitor.NewScheduler(mnt)
	sch.DoWithIntervals(time.Minute * 3)

	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello world")
	})

	if err := e.Start("127.0.0.1:1358"); err != nil {
		log.Fatalf("cannot start the echo http server: %s", err)
	}

}
