package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct{}

func (h Server) Query(ctx echo.Context) error {
	req := new(Req)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, Res{
		Answer: strings.ToUpper(req.Query),
	})
}
func (h Server) OpenApi(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/yaml; charset=utf-8")
	return ctx.File("openapi.yaml")
}
func (h Server) AiPlugin(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/json; charset=utf-8")
	return ctx.File("ai-plugin.json")
}

func main() {
	e := echo.New()
	s := Server{}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/.well-known/ai-plugin.json", s.AiPlugin)
	e.GET("/openapi.yaml", s.OpenApi)

	RegisterHandlers(e, s)

	e.Logger.Fatal(e.Start(":3333"))
}
