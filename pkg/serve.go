package devise

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/autonomy/devise/pkg/discoverer"
	"github.com/autonomy/devise/pkg/modifier"
	"github.com/autonomy/devise/pkg/renderer"
	"github.com/autonomy/devise/pkg/storage"
	"github.com/autonomy/devise/pkg/storage/datastore"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	homeTemplate = "home"
)

// Server represents the Devise server.
type Server struct {
	Storage datastore.Datastore
}

// ServeOptions is used to configure the server.
type ServeOptions struct {
	Storage      string
	BackendPort  string
	UIPort       string
	VaultAddress string
	Discoverers  []string
}

// Start starts the server.
func Start(opts *ServeOptions) {
	go discover(opts.Storage, opts.VaultAddress, opts.Discoverers)
	ui(opts.UIPort)
}

func discover(datastore, vaultAddress string, d []string) {
	modifiers := modifier.NewModifiers(&modifier.Options{VaultAddress: vaultAddress})
	renderer := renderer.NewRenderer(modifiers)
	storage := storage.NewDatastore(datastore)
	discoverers := discoverer.NewDiscoverers(d)
	for _, discoverer := range discoverers {
		go discoverer.Discover(storage, renderer)
	}
}

// Renderer implements echo.Renderer
type Renderer struct {
	templates *template.Template
}

// Render implements echo.Renderer
func (t *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func healthz(c echo.Context) (err error) {
	return
}

func home(c echo.Context) (err error) {
	return c.Render(http.StatusOK, homeTemplate, nil)
}

func ui(port string) {
	e := echo.New()
	r := &Renderer{
		templates: template.Must(template.ParseGlob("assets/templates/*.html")),
	}
	e.Renderer = r
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./assets")
	e.GET("/", home)
	e.GET("/healthz", healthz)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
