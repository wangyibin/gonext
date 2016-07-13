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
	EchoContext echo.Context
}

func NewGonextContextFromEcho(c echo.Context) Context {
	return &echoContextWrapper{c}
}

func (w *echoContextWrapper) Context() netContext.Context {
	return w.EchoContext.Context()
}

func (w *echoContextWrapper) SetContext(ctx netContext.Context) {
	w.EchoContext.SetContext(ctx)
}

func (w *echoContextWrapper) Deadline() (deadline time.Time, ok bool) {
	return w.EchoContext.Deadline()
}

func (w *echoContextWrapper) Done() <-chan struct{} {
	return w.EchoContext.Done()
}

func (w *echoContextWrapper) Err() error {
	return w.EchoContext.Err()
}

func (w *echoContextWrapper) Value(key interface{}) interface{} {
	return w.EchoContext.Value(key)
}

//func (w *echoContextWrapper) Handle(ctx Context) error {
//	return w.c.Handle(ctx)
//}

func (w *echoContextWrapper) Request() (r engine.Request) {
	return &requestWrapperFromEcho{w.EchoContext.Request()}
}

func (w *echoContextWrapper) Response() engine.Response {
	return &responseWrapperFromEcho{w.EchoContext.Response()}
}

func (w *echoContextWrapper) Path() string {
	return w.EchoContext.Path()
}

func (w *echoContextWrapper) SetPath(arg string) {
	w.EchoContext.SetPath(arg)
}

func (w *echoContextWrapper) P(i int) (string) {
	return w.EchoContext.P(i)
}

func (w *echoContextWrapper) Param(name string) (string) {
	return w.EchoContext.Param(name)
}

func (w *echoContextWrapper) ParamNames() []string {
	return w.EchoContext.ParamNames()
}

func (w *echoContextWrapper) SetParamNames(args ...string) {
	w.EchoContext.SetParamNames(args...)
}

func (w *echoContextWrapper) ParamValues() []string {
	return w.EchoContext.ParamNames()
}

func (w *echoContextWrapper) SetParamValues(args ...string) {
	w.EchoContext.SetParamValues(args...)
}

func (w *echoContextWrapper) QueryParam(name string) string {
	return w.EchoContext.QueryParam(name)
}

func (w *echoContextWrapper) QueryParams() map[string][]string {
	return w.EchoContext.QueryParams()
}

func (w *echoContextWrapper) FormValue(name string) string {
	return w.EchoContext.FormValue(name)
}

func (w *echoContextWrapper) FormParams() map[string][]string {
	return w.EchoContext.FormParams()
}

func (w *echoContextWrapper) FormFile(name string) (*multipart.FileHeader, error) {
	return w.EchoContext.FormFile(name)
}

func (w *echoContextWrapper) MultipartForm() (*multipart.Form, error) {
	return w.EchoContext.MultipartForm()
}

func (w *echoContextWrapper) Cookie(key string) (engine.Cookie, error) {
	return w.EchoContext.Cookie(key)
}

func (w *echoContextWrapper) SetCookie(cookie engine.Cookie) {
	w.EchoContext.SetCookie(cookie)
}

func (w *echoContextWrapper) Cookies() []engine.Cookie {
	var cookies []engine.Cookie
	for _, c := range w.EchoContext.Cookies() {
		cookies = append(cookies, c)
	}
	return cookies
}

func (w *echoContextWrapper) Set(key string, val interface{}) {
	w.EchoContext.Set(key, val)
}

func (w *echoContextWrapper) Get(key string) interface{} {
	return w.EchoContext.Get(key)
}

func (w *echoContextWrapper) Bind(i interface{}) error {
	return w.EchoContext.Bind(i)
}

func (w *echoContextWrapper) Render(code int, name string, data interface{}) error {
	return w.EchoContext.Render(code, name, data)
}

func (w *echoContextWrapper) HTML(code int, html string) error {
	return w.EchoContext.HTML(code, html)
}

func (w *echoContextWrapper) String(code int, s string) error {
	return w.EchoContext.String(code, s)
}

func (w *echoContextWrapper) JSON(code int, i interface{}) error {
	return w.EchoContext.JSON(code, i)
}

func (w *echoContextWrapper) JSONBlob(code int, b []byte) error {
	return w.EchoContext.JSONBlob(code, b)
}

func (w *echoContextWrapper) JSONP(code int, callback string, i interface{}) error {
	return w.EchoContext.JSONP(code, callback, i)
}

func (w *echoContextWrapper) XML(code int, i interface{}) error {
	return w.EchoContext.XML(code, i)
}

func (w *echoContextWrapper) XMLBlob(code int, b []byte) error {
	return w.EchoContext.XMLBlob(code, b)
}

func (w *echoContextWrapper) File(file string) error {
	return w.EchoContext.File(file)
}

func (w *echoContextWrapper) Attachment(r io.ReadSeeker, name string) error {
	return w.EchoContext.Attachment(r, name)
}

func (w *echoContextWrapper) NoContent(code int) error {
	return w.EchoContext.NoContent(code)
}

func (w *echoContextWrapper) Redirect(code int, url string) error {
	return w.EchoContext.Redirect(code, url)
}

func (w *echoContextWrapper) Error(err error) {
	w.EchoContext.Error(err)
}

// Handler returns the matched handler by router.
func (w *echoContextWrapper) Handler() HandlerFunc {
	return func(c Context) error {
		echoC := c.(*echoContextWrapper).EchoContext
		return w.EchoContext.Handler()(echoC)
	}
}

// SetHandler sets the matched handler by router.
func (w *echoContextWrapper) SetHandler(h HandlerFunc) {
	w.EchoContext.SetHandler(func (c echo.Context) error {
		return h(NewGonextContextFromEcho(c))
	})
}

func (w *echoContextWrapper) ServeContent(content io.ReadSeeker, name string, modtime time.Time) error {
	return w.EchoContext.ServeContent(content, name, modtime)
}

func (w *echoContextWrapper) Reset(rq engine.Request, rs engine.Response) {
	request := rq.(*requestWrapperFromEcho)
	response := rs.(*responseWrapperFromEcho)
	w.EchoContext.Reset(request.Request, response.Response)
}
