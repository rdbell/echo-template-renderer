package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

// httpErrorHandler is a custom HTTP error responder.
func httpErrorHandler(err error, c echo.Context) {
	var e struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// Log error
	c.Logger().Error(err)

	// Make the error pretty
	re := regexp.MustCompile(`code=.*, message=`)
	readableError := re.ReplaceAllString(err.Error(), "")
	e.Message = readableError

	// Attempt to catch error code
	// Default error code 500 for uncaught errors
	var he *echo.HTTPError
	if ok := errors.As(err, &he); ok {
		e.Code = he.Code
	} else {
		e.Code = http.StatusInternalServerError
	}

	// Setup page
	pd := &PageData{
		Title: fmt.Sprintf("Notice - %d", e.Code),
		Data:  e,
	}

	// Render
	_ = c.Render(e.Code, "base:error", pd)
}
