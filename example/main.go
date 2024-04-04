package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	renderer "github.com/rdbell/echo-template-renderer"
)

var (
	cacheBuster int64
	//go:embed assets/*
	assets embed.FS
	//go:embed views/*
	views embed.FS
)

func init() {
	cacheBuster = time.Now().Unix()
}

var funcMap = map[string]interface{}{
	"cacheBuster": func() int64 {
		return cacheBuster
	},
	"appName": func() string {
		return "Echo Template Renderer"
	},
	"add": func(a, b int) float64 {
		return float64(a + b)
	},
	"dollarFormat": func(amount float64) string {
		return fmt.Sprintf("$%.2f", amount)
	},
}

func main() {
	// Start a new Echo instance
	e := echo.New()

	// Define some routes
	e.GET("/", indexRoute)
	e.GET("/auth/login", loginRoute)
	e.GET("/auth/signup", signupRoute)

	// Static assets
	e.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(assets))))

	// Error handler
	e.HTTPErrorHandler = httpErrorHandler

	// Set up the renderer with views and functions
	cfg := &renderer.Config{
		Views: views,
		Funcs: funcMap,
	}
	e.Renderer = renderer.New(cfg)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
