# Echo Template Renderer

This Go package provides a powerful and flexible template rendering solution for the Echo web framework. It allows the seamless integration of HTML templates with your Echo server, supporting both embedded file systems and custom template functions for a dynamic and responsive web application.

## Features

- **Embedded File System Support**: Directly integrates with Go's `embed` package, allowing you to bundle your HTML templates into your application binary.
- **Custom Template Functions**: Easily add custom template functions to enhance your templates with logic directly accessible within your HTML files.
- **Efficient Template Loading**: Templates are pre-loaded and cached for optimal performance, minimizing IO operations and processing time.
- **Flexible Rendering Options**: Supports rendering of layouts, pages, and shared components, offering great flexibility in organizing your template files.

## Installation

To install this package, use the go get command:


1. Import the Echo Template Renderer package into your Echo application.
2. Assign `renderer.New(cfg)` to the `Renderer` field of your Echo instance like so:

```go
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
```

3. Render a template in your route handler by calling `c.Render()` with the template name and optional page data:

```go
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
```

## Demo

To demo Echo Template Renderer's capabilities, navigate to the [`example`](/example) folder within this repository and run `go run main.go` in your terminal.

## Contribution

Contributions are welcome! If you'd like to improve Echo Template Renderer or suggest new features, feel free to fork the repository, make your changes, and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](/LICENSE) file for more details.
