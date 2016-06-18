package gonext

import (
	"github.com/labstack/echo"
	netContext "golang.org/x/net/context"
	"time"
	"github.com/wangyibin/gonext/engine"
	"mime/multipart"
	"io"
)

type echoContextWrapper struct {
	c echo.Context
}

func NewGonextContextFromEcho(c echo.Context) Context {
	return &echoContextWrapper{c}
}

func (w *echoContextWrapper) Context() netContext.Context {
	return w.c.Context()
}

func (w *echoContextWrapper) SetContext(ctx netContext.Context) {
	w.c.SetContext(ctx)
}

func (w *echoContextWrapper) Deadline() (deadline time.Time, ok bool) {
	return w.c.Deadline()
}

func (w *echoContextWrapper) Done() <-chan struct{} {
	return w.c.Done()
}

func (w *echoContextWrapper) Err() error {
	return w.c.Err()
}

func (w *echoContextWrapper) Value(key interface{}) interface{} {
	return w.c.Value(key)
}

//func (w *echoContextWrapper) Handle(ctx Context) error {
//	return w.c.Handle(ctx)
//}

func (w *echoContextWrapper) Request() (r engine.Request) {
	return &requestWrapperFromEcho{w.c.Request()}
}

func (w *echoContextWrapper) Response() engine.Response {
	return &responseWrapperFromEcho{w.c.Response()}
}

func (w *echoContextWrapper) Path() string {
	return w.c.Path()
}

func (w *echoContextWrapper) P(i int) (string) {
	return w.c.P(i)
}

func (w *echoContextWrapper) Param(name string) (string) {
	return w.c.Param(name)
}

func (w *echoContextWrapper) ParamNames() []string {
	return w.c.ParamNames()
}

func (w *echoContextWrapper) QueryParam(name string) string {
	return w.c.QueryParam(name)
}

func (w *echoContextWrapper) QueryParams() map[string][]string {
	return w.c.QueryParams()
}

func (w *echoContextWrapper) FormValue(name string) string {
	return w.c.FormValue(name)
}

func (w *echoContextWrapper) FormParams() map[string][]string {
	return w.c.FormParams()
}

func (w *echoContextWrapper) FormFile(name string) (*multipart.FileHeader, error) {
	return w.c.FormFile(name)
}

func (w *echoContextWrapper) MultipartForm() (*multipart.Form, error) {
	return w.c.MultipartForm()
}

func (w *echoContextWrapper) Set(key string, val interface{}) {
	w.c.Set(key, val)
}

func (w *echoContextWrapper) Get(key string) interface{} {
	return w.c.Get(key)
}

func (w *echoContextWrapper) Bind(i interface{}) error {
	return w.c.Bind(i)
}

func (w *echoContextWrapper) Render(code int, name string, data interface{}) error {
	return w.c.Render(code, name, data)
}

func (w *echoContextWrapper) HTML(code int, html string) error {
	return w.c.HTML(code, html)
}

func (w *echoContextWrapper) String(code int, s string) error {
	return w.c.String(code, s)
}

func (w *echoContextWrapper) JSON(code int, i interface{}) error {
	return w.c.JSON(code, i)
}

func (w *echoContextWrapper) JSONBlob(code int, b []byte) error {
	return w.c.JSONBlob(code, b)
}

func (w *echoContextWrapper) JSONP(code int, callback string, i interface{}) error {
	return w.c.JSONP(code, callback, i)
}

func (w *echoContextWrapper) XML(code int, i interface{}) error {
	return w.c.XML(code, i)
}

func (w *echoContextWrapper) XMLBlob(code int, b []byte) error {
	return w.c.XMLBlob(code, b)
}

func (w *echoContextWrapper) File(file string) error {
	return w.c.File(file)
}

func (w *echoContextWrapper) Attachment(r io.ReadSeeker, name string) error {
	return w.c.Attachment(r, name)
}

func (w *echoContextWrapper) NoContent(code int) error {
	return w.c.NoContent(code)
}

func (w *echoContextWrapper) Redirect(code int, url string) error {
	return w.c.Redirect(code, url)
}

func (w *echoContextWrapper) Error(err error) {
	w.c.Error(err)
}

func (w *echoContextWrapper) ServeContent(content io.ReadSeeker, name string, modtime time.Time) error {
	return w.c.ServeContent(content, name, modtime)
}

func (w *echoContextWrapper) Reset(rq engine.Request, rs engine.Response) {
	request := rq.(*requestWrapperFromEcho)
	response := rs.(*responseWrapperFromEcho)
	w.c.Reset(request.Request, response.Response)
}
