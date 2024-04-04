package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PageData struct {
	Title string
	Data  interface{}
}

func indexRoute(c echo.Context) error {
	// Setup page
	pd := &PageData{
		Title: "Home",
	}

	pd.Data = map[string]interface{}{
		"Message": "Hello, World!",
	}

	// Render
	return c.Render(http.StatusOK, "base:index", pd)
}

func signupRoute(c echo.Context) error {
	// Setup page
	pd := &PageData{
		Title: "Signup",
	}

	// Render
	return c.Render(http.StatusOK, "base:auth/signup", pd)
}

func loginRoute(c echo.Context) error {
	// Setup page
	pd := &PageData{
		Title: "Login",
	}

	// Render
	return c.Render(http.StatusOK, "base:auth/login", pd)
}
