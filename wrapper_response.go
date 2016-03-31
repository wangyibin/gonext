package gonext

import (
	echoEngine "github.com/labstack/echo/engine"
	"github.com/wangyibin/gonext/engine"
	"io"
)

type responseWrapperFromEcho struct {
	echoEngine.Response
}

// Header implements `engine.Response#Header` function.
func (r *responseWrapperFromEcho) Header() engine.Header {
	return &headerWrapperFromEcho{r.Response.Header()}
}

// WriteHeader implements `engine.Response#WriteHeader` function.
func (r *responseWrapperFromEcho) WriteHeader(code int) {
	r.Response.WriteHeader(code)
}

// Write implements `engine.Response#Write` function.
func (r *responseWrapperFromEcho) Write(b []byte) (int, error) {
	return r.Response.Write(b)
}

// Status implements `engine.Response#Status` function.
func (r *responseWrapperFromEcho) Status() int {
	return r.Response.Status()
}

// Size implements `engine.Response#Size` function.
func (r *responseWrapperFromEcho) Size() int64 {
	return r.Response.Size()
}

// Committed implements `engine.Response#Committed` function.
func (r *responseWrapperFromEcho) Committed() bool {
	return r.Response.Committed()
}

// Writer implements `engine.Response#Writer` function.
func (r *responseWrapperFromEcho) Writer() io.Writer {
	return r.Response.Writer()
}

// SetWriter implements `engine.Response#SetWriter` function.
func (r *responseWrapperFromEcho) SetWriter(w io.Writer) {
	r.Response.SetWriter(w)
}