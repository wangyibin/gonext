package gonext

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/log"
)

type (
	Engine struct {
		echo *echo.Echo
	}

	HTTPError struct {
		Code    int
		Message string
	}

)

var (
	DefaultEngine *Engine
)

// The one and only init to the whole package
func init() {
	DefaultEngine = new()
	DefaultEngine.echo.Get("/api-docs", getdoc())
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
func NewGroup(tag string, description string, prefix string) (*Group) {
	return &Group{tag: tag, description: description, prefix: prefix,
		echoGroup: DefaultEngine.echo.Group(prefix)}
}

func SetHTTPErrorHandler(errorHandler func(error, Context)) {
	DefaultEngine.echo.SetHTTPErrorHandler(func (err error, eCtx echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			err = &HTTPError{Code: he.Code, Message: he.Message}
		}
		errorHandler(err, NewGonextContextFromEcho(eCtx))
	})
}
// Error makes it compatible with `error` interface.
func (e *HTTPError) Error() string {
	return e.Message
}
func Logger() log.Logger {
	return DefaultEngine.echo.Logger()
}
func NewHTTPError(code int, msg ...string) *HTTPError {
	return echo.NewHTTPError(code, msg...)
}