package gonext

import (
	"reflect"
	"strings"

	"github.com/labstack/echo"
)

// Router struct
type (
	Group struct {
		tag         string
		description string
		prefix      string
		echoGroup   *echo.Group
	}
// Middleware defines an interface for middleware via `Handle(Handler) Handler`
// function.
//	Middleware interface {
//		Handle(Handler) Handler
//	}

// MiddlewareFunc is an adapter to allow the use of `func(Handler) Handler` as
// middleware.
//	MiddlewareFunc func(Handler) Handler

// Handler defines an interface to server HTTP requests via `Handle(Context)`
// function.
	Handler interface {
		Handle(Context) error
	}

// HandlerFunc is an adapter to allow the use of `func(Context)` as an HTTP
// handler.
	HandlerFunc func(Context) error
)


// Group func
func (g *Group) Group(tag string, description string, prefix string) *Group {
	return &Group{tag: tag, description: description, prefix: g.prefix + prefix, echoGroup: g.echoGroup.Group(prefix)}
}

// Get func
func (g *Group) Get(path string, actions ...interface{}) {
	mount("GET", g, path, actions)
}

// Post func
func (g *Group) Post(path string, actions ...interface{}) {
	mount("POST", g, path, actions)
}

// Put func
func (g *Group) Put(path string, actions ...interface{}) {
	mount("PUT", g, path, actions)
}

// Delete func
func (g *Group) Delete(path string, actions ...interface{}) {
	mount("DELETE", g, path, actions)
}

// Mount func
func mount(method string, g *Group, path string, actions []interface{}) {
	var summary, description string
	var handlers []interface{}
	for _, a := range actions {
		if reflect.TypeOf(a).Kind() == reflect.String {
			if len(summary) == 0 {
				summary = a.(string)
			} else {
				description = a.(string)
			}
		} else {
			handlers = append(handlers, a)
		}
	}

	SwaggerTags[g.tag] = g.description
	fullPath := g.prefix + path
	MountSwaggerPath(&SwaggerPathDefine{Tag: g.tag, Method: method, Path: fullPath,
		Summary: summary, Description: description, Handlers: handlers})

	echoHandler := BuildEchoHandler(fullPath, handlers)
	switch strings.ToUpper(method) {
	case "GET":
		g.echoGroup.Get(path, echoHandler)
	case "POST":
		g.echoGroup.Post(path, echoHandler)
	case "PUT":
		g.echoGroup.Put(path, echoHandler)
	case "DELETE":
		g.echoGroup.Delete(path, echoHandler)
	}
}

// // AddHandlers func
// func (hdef *HandlerDef) AddHandlers(handlers ...interface{}) *HandlerDef {
// 	hdef.handlers = handlers
// 	return hdef
// }
//
// // Summary func
// func (hdef *HandlerDef) Summary(summary string) *HandlerDef {
// 	hdef.summary = summary
// 	return hdef
// }
//
// // Description func
// func (hdef *HandlerDef) Description(description string) *HandlerDef {
// 	hdef.description = description
// 	return hdef
// }

func containsIgnoreCase(s []string, e string) bool {
	for _, a := range s {
		if strings.ToUpper(a) == strings.ToUpper(e) {
			return true
		}
	}
	return false
}



func lowCamelStr(str string) string {
	for _, word := range []string{"ID", "URL", "URI"} {
		if word == str {
			return strings.ToLower(str)
		}
	}
	return strings.ToLower(string(str[0])) + string(str[1:])
}
