package ree

import (
	"Ree/ree/bind"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	route *Route
}

type Route struct {
	handlers map[string]HandlerFunc
}

type H map[string]any

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func New() *Engine {
	return &Engine{route: newRoute()}
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Request:        request,
		ResponseWriter: writer,
	}
}

func newRoute() *Route {
	return &Route{handlers: make(map[string]HandlerFunc)}
}

func (r *Route) handle(ctx *Context) {
	key := ctx.Request.Method + "-" + ctx.Request.URL.Path
	if handlerFunc, ok := r.handlers[key]; ok {
		handlerFunc(ctx)
	} else {
		notFoundHandler(ctx)
	}
}

func (e *Engine) addRoute(method string, path string, handler HandlerFunc) {
	key := method + "-" + path
	e.route.handlers[key] = handler
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute(http.MethodGet, path, handler)
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.addRoute(http.MethodPost, path, handler)
}

func (e *Engine) PUT(path string, handler HandlerFunc) {
	e.addRoute(http.MethodPut, path, handler)
}

func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.addRoute(http.MethodDelete, path, handler)
}

func notFoundHandler(ctx *Context) {
	ctx.String(http.StatusNotFound, "404 NOT FOUND")
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := newContext(writer, request)
	e.route.handle(c)
}

func (e *Engine) Start(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (ctx *Context) JSON(code int, data interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.ResponseWriter)
	err := encoder.Encode(data)
	if err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), 500)
	}
}

// HTML
func (ctx *Context) HTML(code int, data interface{}) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	switch data.(type) {
	case *os.File:
		bytes, err := io.ReadAll(data.(*os.File))
		if err != nil {
			log.Println(err)
			return
		}
		_, err = ctx.ResponseWriter.Write(bytes)
		if err != nil {
			http.Error(ctx.ResponseWriter, err.Error(), 500)
		}
		return
	case string:
		_, err := ctx.ResponseWriter.Write([]byte(data.(string)))
		if err != nil {
			http.Error(ctx.ResponseWriter, err.Error(), 500)
		}
		return
	case []byte:
		_, err := ctx.ResponseWriter.Write(data.([]byte))
		if err != nil {
			http.Error(ctx.ResponseWriter, err.Error(), 500)
		}
		return
	default:
		log.Println("unsupported data type!")
	}
	return
}

func (ctx *Context) String(code int, data string) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	_, err := ctx.ResponseWriter.Write([]byte(data))
	if err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), 500)
	}
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

func (ctx *Context) Status(code int) {
	ctx.ResponseWriter.WriteHeader(code)
}

func (ctx *Context) ShouldBind(data any) error {
	return bind.ShouldBind(ctx.Request, ctx.ResponseWriter, data)
}

func (ctx *Context) BindJSON(data any) error {
	return bind.BindJSON(ctx.Request, data)
}

func (ctx *Context) BindXML(data any) error {
	return bind.BindXML(ctx.Request, data)
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Request.PostFormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Request.URL.Query().Get(key)
}
