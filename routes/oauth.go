package routes

import (
	"github.com/labstack/echo"
	. "quote/api"
)

func SetQuoteInterfaces(e *echo.Echo) {

	e.GET("/api/quotes", GetApiQuotes)

}
