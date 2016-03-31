package gonext

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/echo/engine/standard"
)

type (
	Engine struct {
		echo *echo.Echo
	}

)

var (
	DefaultEngine *Engine
)

// The one and only init to the whole package
func init() {
	DefaultEngine = new()
}
func new() (*Engine) {
	e := echo.New()
	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	return &Engine{echo: echo.New()}
}

func Static(prefix, root string) {
	DefaultEngine.echo.Static(prefix, root)
}

// Run starts the HTTP server.
func Run(addr string) {
	DefaultEngine.echo.Run(standard.New(addr))
}

// Group creates a new router group with prefix and optional group-level middleware.
func Group(tag string, description string, prefix string, m ...Middleware) (*Group) {
	return &Group{tag: tag, description: description, prefix: prefix,
		echoGroup: DefaultEngine.echo.Group(prefix, m...)}
}