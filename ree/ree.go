package ree

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	m map[string]HandlerFunc
}

type H map[string]any

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func New() *Engine {
	return &Engine{m: make(map[string]HandlerFunc)}
}

func (e *Engine) addRoute(method string, path string, handler HandlerFunc) {
	key := method + "-" + path
	e.m[key] = handler
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

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "404 NOT FOUND -- %s", request.URL)
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	context := newContext(writer, request)
	if v, ok := e.m[key]; ok {
		v(context)
	} else {
		notFoundHandler(writer, request)
	}
}

func (e *Engine) Start(addr string) error {
	return http.ListenAndServe(addr, e)
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Request:        request,
		ResponseWriter: writer,
	}
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

// HTML todo 完成html渲染
func (ctx *Context) HTML(code int, data interface{}) {

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

func (ctx *Context) BindJSON(data interface{}) error {
	bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &data)
	return err
}
