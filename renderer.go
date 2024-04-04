package renderer

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// Config is a struct for template renderer configuration.
type Config struct {
	Views embed.FS
	Funcs template.FuncMap
}

// Renderer struct for template rendering.
type Renderer struct {
	templates map[string]*template.Template
}

// New creates a new instance of the Renderer.
func New(config *Config) echo.Renderer {
	renderer := &Renderer{
		templates: make(map[string]*template.Template),
	}

	loadTemplates(renderer, config.Views, config.Funcs)

	return renderer
}

// loadTemplates pre-computes templates in the views folder.
func loadTemplates(renderer *Renderer, views embed.FS, funcMap template.FuncMap) {
	templatePaths := []struct {
		prefix   string
		fileList *[]string
	}{
		{"views/layouts", new([]string)},
		{"views/pages", new([]string)},
		{"views/shared", new([]string)},
	}

	// Walk the views folder and add all layouts/pages/shared to be combined into templates
	for _, tp := range templatePaths {
		if err := templatesFromPath(tp.prefix, tp.fileList, views); err != nil {
			log.Fatal(err)
		}
	}

	// Loop through layouts and pages to combine them into templates
	layouts, pages, shared := *templatePaths[0].fileList, *templatePaths[1].fileList, *templatePaths[2].fileList
	for _, layout := range layouts {
		for _, page := range pages {
			layoutBase, pageBase := filepath.Base(layout), strings.Replace(page, "views/pages/", "", 1)

			// Find indexes safely
			layoutIndex := strings.Index(layoutBase, ".html.tmpl")
			pageIndex := strings.Index(pageBase, ".html.tmpl")
			if layoutIndex == -1 {
				log.Fatalf("File %s does not have a '.html.tmpl' extension", layoutBase)
			}
			if pageIndex == -1 {
				log.Fatalf("File %s does not have a '.html.tmpl' extension", pageBase)
			}

			// Shorten names for template map
			// "layouts/base.html.tmpl" -> "base"
			// "pages/auth/login.html.tmpl" -> "auth/login"
			layoutShort, pageShort := layoutBase[:layoutIndex], pageBase[:pageIndex]
			layoutContent, _ := views.ReadFile(layout)
			pageContent, _ := views.ReadFile(page)
			combinedTemplate := string(layoutContent) + string(pageContent)

			// Combine shared templates
			for _, sharedName := range shared {
				sharedContent, _ := views.ReadFile(sharedName)
				combinedTemplate += string(sharedContent)
			}

			// Build a name like "base:auth/login" to be used when calling Render()
			name := layoutShort + ":" + pageShort

			// Add template to renderer
			renderer.templates[name] = template.Must(template.New(pageBase).Funcs(funcMap).Parse(combinedTemplate))
		}
	}
}

// templateFromPath recursively loads template file data from a given path.
func templatesFromPath(path string, destination *[]string, e embed.FS) error {
	return fs.WalkDir(e, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Skip non tmpl files
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Add file
		*destination = append(*destination, path)

		return nil
	})
}

// Render renders templates for http responses.
// Use like this: return c.Render(http.StatusOK, "base:index", pageData)
// See examples in the example folder.
func (t *Renderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}

	return tmpl.Execute(w, data)
}
