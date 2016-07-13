package gonext

import (
	"golang.org/x/net/context"
	"github.com/wangyibin/gonext/engine"
	"mime/multipart"
	"io"
	"time"
)

type (
// Context represents the context of the current HTTP request. It holds request and
// response objects, path, path parameters, data and registered handler.
	Context interface {
		// Context returns `net/context.Context`.
		Context() context.Context

		// SetContext sets `net/context.Context`.
		SetContext(context.Context)

		// Deadline returns the time when work done on behalf of this context
		// should be canceled.  Deadline returns ok==false when no deadline is
		// set.  Successive calls to Deadline return the same results.
		Deadline() (deadline time.Time, ok bool)

		// Done returns a channel that's closed when work done on behalf of this
		// context should be canceled.  Done may return nil if this context can
		// never be canceled.  Successive calls to Done return the same value.
		Done() <-chan struct{}

		// Err returns a non-nil error value after Done is closed.  Err returns
		// Canceled if the context was canceled or DeadlineExceeded if the
		// context's deadline passed.  No other values for Err are defined.
		// After Done is closed, successive calls to Err return the same value.
		Err() error

		// Value returns the value associated with this context for key, or nil
		// if no value is associated with key.  Successive calls to Value with
		// the same key returns the same result.
		Value(key interface{}) interface{}

		// Request returns `engine.Request` interface.
		Request() engine.Request

		// Request returns `engine.Response` interface.
		Response() engine.Response

		// Path returns the registered path for the handler.
		Path() string

		// SetPath sets the registered path for the handler.
		SetPath(string)

		// P returns path parameter by index.
		P(int) string

		// Param returns path parameter by name.
		Param(string) string

		// ParamNames returns path parameter names.
		ParamNames() []string

		// SetParamNames sets path parameter names.
		SetParamNames(...string)

		// ParamValues returns path parameter values.
		ParamValues() []string

		// SetParamValues sets path parameter values.
		SetParamValues(...string)

		// QueryParam returns the query param for the provided name. It is an alias
		// for `engine.URL#QueryParam()`.
		QueryParam(string) string

		// QueryParams returns the query parameters as map.
		// It is an alias for `engine.URL#QueryParams()`.
		QueryParams() map[string][]string

		// FormValue returns the form field value for the provided name. It is an
		// alias for `engine.Request#FormValue()`.
		FormValue(string) string

		// FormParams returns the form parameters as map.
		// It is an alias for `engine.Request#FormParams()`.
		FormParams() map[string][]string

		// FormFile returns the multipart form file for the provided name. It is an
		// alias for `engine.Request#FormFile()`.
		FormFile(string) (*multipart.FileHeader, error)

		// MultipartForm returns the multipart form.
		// It is an alias for `engine.Request#MultipartForm()`.
		MultipartForm() (*multipart.Form, error)

		// Cookie returns the named cookie provided in the request.
		// It is an alias for `engine.Request#Cookie()`.
		Cookie(string) (engine.Cookie, error)

		// SetCookie adds a `Set-Cookie` header in HTTP response.
		// It is an alias for `engine.Response#SetCookie()`.
		SetCookie(engine.Cookie)

		// Cookies returns the HTTP cookies sent with the request.
		// It is an alias for `engine.Request#Cookies()`.
		Cookies() []engine.Cookie

		// Get retrieves data from the context.
		Get(string) interface{}

		// Set saves data in the context.
		Set(string, interface{})

		// Bind binds the request body into provided type `i`. The default binder
		// does it based on Content-Type header.
		Bind(interface{}) error

		// Render renders a template with data and sends a text/html response with status
		// code. Templates can be registered using `Echo.SetRenderer()`.
		Render(int, string, interface{}) error

		// HTML sends an HTTP response with status code.
		HTML(int, string) error

		// String sends a string response with status code.
		String(int, string) error

		// JSON sends a JSON response with status code.
		JSON(int, interface{}) error

		// JSONBlob sends a JSON blob response with status code.
		JSONBlob(int, []byte) error

		// JSONP sends a JSONP response with status code. It uses `callback` to construct
		// the JSONP payload.
		JSONP(int, string, interface{}) error

		// XML sends an XML response with status code.
		XML(int, interface{}) error

		// XMLBlob sends a XML blob response with status code.
		XMLBlob(int, []byte) error

		// File sends a response with the content of the file.
		File(string) error

		// Attachment sends a response from `io.ReaderSeeker` as attachment, prompting
		// client to save the file.
		Attachment(io.ReadSeeker, string) error

		// NoContent sends a response with no body and a status code.
		NoContent(int) error

		// Redirect redirects the request with status code.
		Redirect(int, string) error

		// Error invokes the registered HTTP error handler. Generally used by middleware.
		Error(err error)

		// Handler returns the matched handler by router.
		Handler() HandlerFunc

		// SetHandler sets the matched handler by router.
		SetHandler(HandlerFunc)

		// Logger returns the `Logger` instance.
		//Logger() log.Logger

		// Echo returns the `Echo` instance.
		//Echo() *Echo

		// ServeContent sends static content from `io.Reader` and handles caching
		// via `If-Modified-Since` request header. It automatically sets `Content-Type`
		// and `Last-Modified` response headers.
		ServeContent(io.ReadSeeker, string, time.Time) error

		// Reset resets the context after request completes. It must be called along
		// with `Echo#AcquireContext()` and `Echo#ReleaseContext()`.
		// See `Echo#ServeHTTP()`
		Reset(engine.Request, engine.Response)
	}

)

