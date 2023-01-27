package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt"
	"github.com/hamidreza-abooei/ie-project/db"
	"github.com/hamidreza-abooei/ie-project/monitor"
	"github.com/labstack/echo/v4"
)

// require validator to add "required" tag to every struct field in the package
func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Handler type for GET/POST
type Handler struct {
	st  *db.Store
	sch *monitor.Scheduler
}

// NewHandler creates a new handler with given store instance
func NewHandler(st *db.Store, sch *monitor.Scheduler) *Handler {
	return &Handler{st: st, sch: sch}
}

func extractID(c echo.Context) uint {
	e := c.Get("user").(*jwt.Token)
	claims := e.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))
	return id
}
