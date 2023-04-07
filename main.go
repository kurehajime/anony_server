package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

func main() {
	e := echo.New()
	s := Server{}

	RegisterHandlers(e, s)

	e.Logger.Fatal(e.Start(":8080"))
}
