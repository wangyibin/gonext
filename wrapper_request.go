package gonext

import (
	echoEngine "github.com/labstack/echo/engine"
	"github.com/wangyibin/gonext/engine"
	"io"
	"mime/multipart"
)

type requestWrapperFromEcho struct {
	echoEngine.Request
}

// IsTLS implements `engine.Request#TLS` function.
func (r *requestWrapperFromEcho) IsTLS() bool {
	return r.Request.IsTLS()
}

// Scheme implements `engine.Request#Scheme` function.
func (r *requestWrapperFromEcho) Scheme() string {
	return r.Request.Scheme()
}

// Host implements `engine.Request#Host` function.
func (r *requestWrapperFromEcho) Host() string {
	return r.Request.Host()
}

// URL implements `engine.Request#URL` function.
func (r *requestWrapperFromEcho) URL() engine.URL {
	return r.Request.URL()
}

// Header implements `engine.Request#URL` function.
func (r *requestWrapperFromEcho) Header() engine.Header {
	return &headerWrapperFromEcho{r.Request.Header()}
}

// func Proto() string {
// 	return r.request.Proto()
// }
//
// func ProtoMajor() int {
// 	return r.request.ProtoMajor()
// }
//
// func ProtoMinor() int {
// 	return r.request.ProtoMinor()
// }

// ContentLength implements `engine.Request#ContentLength` function.
func (r *requestWrapperFromEcho) ContentLength() int64 {
	return r.Request.ContentLength()
}

// UserAgent implements `engine.Request#UserAgent` function.
func (r *requestWrapperFromEcho) UserAgent() string {
	return r.Request.UserAgent()
}

// RemoteAddress implements `engine.Request#RemoteAddress` function.
func (r *requestWrapperFromEcho) RemoteAddress() string {
	return r.Request.RemoteAddress()
}

// Method implements `engine.Request#Method` function.
func (r *requestWrapperFromEcho) Method() string {
	return r.Request.Method()
}

// SetMethod implements `engine.Request#SetMethod` function.
func (r *requestWrapperFromEcho) SetMethod(method string) {
	r.Request.SetMethod(method)
}

// URI implements `engine.Request#URI` function.
func (r *requestWrapperFromEcho) URI() string {
	return r.Request.URI()
}

// Body implements `engine.Request#Body` function.
func (r *requestWrapperFromEcho) Body() io.Reader {
	return r.Request.Body()
}

// FormValue implements `engine.Request#FormValue` function.
func (r *requestWrapperFromEcho) FormValue(name string) string {
	return r.Request.FormValue(name)
}

// FormParams implements `engine.Request#FormParams` function.
func (r *requestWrapperFromEcho) FormParams() map[string][]string {
	return r.Request.FormParams()
}

// FormFile implements `engine.Request#FormFile` function.
func (r *requestWrapperFromEcho) FormFile(name string) (*multipart.FileHeader, error) {
	return r.Request.FormFile(name)
}

// MultipartForm implements `engine.Request#MultipartForm` function.
func (r *requestWrapperFromEcho) MultipartForm() (*multipart.Form, error) {
	return r.Request.MultipartForm()
}