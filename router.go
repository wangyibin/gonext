package apidoc

import (
	"strings"

	"github.com/labstack/echo"
)

// Router struct
type Router struct {
	e *echo.Echo
	// resources []*Resource
}

// Group struct
type Group struct {
	tag         string
	description string
	prefix      string
	echoGroup   *echo.Group
}

// HandlerDef struct
type HandlerDef struct {
	method      string
	path        string
	h1          interface{}
	h2          interface{}
	h3          interface{}
	summary     string
	description string
	group       *Group
}

// NewRouter func
func NewRouter(e *echo.Echo) *Router {
	r := &Router{e: e}
	e.Get("/api-docs", getdoc(r))
	return r
}

// Group func
func (r *Router) Group(tag string, description string, prefix string, m ...echo.Middleware) *Group {
	return &Group{tag: tag, description: description, prefix: prefix, echoGroup: r.e.Group(prefix, m...)}
}

// Group func
func (g *Group) Group(tag string, description string, prefix string, m ...echo.Middleware) *Group {
	return &Group{tag: tag, description: description, prefix: g.prefix + prefix, echoGroup: g.echoGroup.Group(prefix, m...)}
}

// Get func
func (g *Group) Get(path string) *HandlerDef {
	return &HandlerDef{method: "GET", path: path, group: g}
}

// Post func
func (g *Group) Post(path string) *HandlerDef {
	return &HandlerDef{method: "POST", path: path, group: g}
}

// AddHandler func
func (hdef *HandlerDef) AddHandler(handler interface{}) *HandlerDef {
	if hdef.h1 == nil {
		hdef.h1 = handler
	} else if hdef.h2 == nil {
		hdef.h2 = handler
	} else if hdef.h3 == nil {
		hdef.h3 = handler
	} else {
		panic("Only can support 3 handler at most")
	}
	return hdef
}

// Mount func
func (hdef *HandlerDef) Mount() {
	g := hdef.group
	SwaggerTags[g.tag] = g.description
	fullPath := g.prefix + hdef.path
	MountSwaggerPath(&SwaggerPathDefine{Tag: g.tag, Method: hdef.method,
		Path: fullPath, Handler1: hdef.h1, Handler2: hdef.h2, Handler3: hdef.h3})

	echoHandler := BuildEchoHandler(fullPath, hdef.h1, hdef.h2, hdef.h3)
	switch strings.ToUpper(hdef.method) {
	case "GET":
		g.echoGroup.Get(hdef.path, echoHandler)
	case "POST":
		g.echoGroup.Post(hdef.path, echoHandler)
	}
}

func containsIgnoreCase(s []string, e string) bool {
	for _, a := range s {
		if strings.ToUpper(a) == strings.ToUpper(e) {
			return true
		}
	}
	return false
}

// PathNames func
func PathNames(path string) []string {
	pnames := []string{} // Param names
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)
		} else if path[i] == '*' {
			pnames = append(pnames, "_*")
		}
	}
	return pnames
}

func lowCamelStr(str string) string {
	return strings.ToLower(string(str[0])) + string(str[1:])
}
