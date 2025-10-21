package httputil

import (
	"github.com/HatefBarari/microblog-shared/pkg/validator"
	"github.com/labstack/echo/v4"
)

func BindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(400, "invalid body")
	}
	if err := validator.Struct(req); err != nil {
		return echo.NewHTTPError(422, err.Error())
	}
	return nil
}