package routes

import (
	"github.com/labstack/echo"
	. "github.com/oldfritter/quote/api"
)

func SetQuoteInterfaces(e *echo.Echo) {

	e.GET("/api/quotes", GetApiQuotes)

}
