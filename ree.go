package ree

import (
	"Ree/ree/binding"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	*Routers
	route   *Route
	routers []*Routers
}

type Route struct {
	handlers map[string]HandlerFunc
	trie     map[string]*node
}

type H map[string]any

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Method         string
	Path           string
	StatusCode     int
	Params         map[string]string
	data           map[string]any
}

func New() *Engine {
	e := &Engine{route: newRoute()}
	e.Routers = &Routers{engine: e}
	e.routers = []*Routers{e.Routers}
	return e
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Request:        request,
		ResponseWriter: writer,
		data:           make(map[string]any),
		Path:           request.URL.Path,
		Method:         request.Method,
	}
}

func newRoute() *Route {
	return &Route{
		handlers: make(map[string]HandlerFunc),
		trie:     make(map[string]*node),
	}
}

func (r *Route) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Path)
	if n != nil {
		ctx.Params = params
		key := ctx.Method + "-" + n.path
		handler := r.handlers[key]
		handler(ctx)
	} else {
		http.NotFound(ctx.ResponseWriter, ctx.Request)
	}
}

func (e *Engine) addRoute(method string, path string, handler HandlerFunc) {
	e.route.addRoute(method, path, handler)
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

// ServeHTTP 将引擎变为一个http.Handler，将每个请求都用route处理
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := newContext(writer, request)
	e.route.handle(c)
}

// Start 启动http引擎，同时实现优雅的退出
func (e *Engine) Start(addr string) error {
	e.route.printRoute()
	srv := &http.Server{
		Addr:    addr,
		Handler: e,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()
	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

// JSON 以json格式返回数据
func (ctx *Context) JSON(code int, data interface{}) {
	ctx.SetHeader("Content-Type", binding.TypeJson)
	ctx.Status(code)
	err := json.NewEncoder(ctx.ResponseWriter).Encode(data)
	if err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), 500)
	}
}

// HTML 渲染html页面
func (ctx *Context) HTML(code int, data interface{}) {
	ctx.SetHeader("Content-Type", binding.TypeHTML)
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

// String 将数据以string格式返回
func (ctx *Context) String(code int, data string) {
	ctx.SetHeader("Content-Type", binding.TypeText)
	ctx.Status(code)
	_, err := ctx.ResponseWriter.Write([]byte(data))
	if err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), 500)
	}
}

// XML 将数据以xml格式返回
func (ctx *Context) XML(code int, data any) {
	ctx.SetHeader("Content-Type", binding.TypeXML)
	ctx.Status(code)
	err := xml.NewEncoder(ctx.ResponseWriter).Encode(&data)
	if err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), 500)
	}
}

// SetHeader 在ResponseWriter中设置一对请求头
func (ctx *Context) SetHeader(key, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// Status 在ResponseWriter中设置状态码
func (ctx *Context) Status(code int) {
	ctx.ResponseWriter.WriteHeader(code)
}

// Bind 根据请求头自动选择绑定类型
// 目前支持 JSON XML 类型的请求头
func (ctx *Context) Bind(data any) error {
	return binding.Bind(ctx.Request, ctx.ResponseWriter, data)
}

func (ctx *Context) BindJSON(data any) error {
	return binding.BindJSON(ctx.Request, data)
}

func (ctx *Context) BindXML(data any) error {
	return binding.BindXML(ctx.Request, data)
}

func (ctx *Context) BindFORM(data any) error {
	return binding.BindFORM(ctx.Request, data)
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Request.PostFormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

func (ctx *Context) Set(key string, value any) {
	if ctx.data == nil {
		ctx.data = make(map[string]any)
	}
	ctx.data[key] = value
}

func (ctx *Context) Get(key string) (any, bool) {
	if ctx.data == nil {
		return nil, false
	}
	value, ok := ctx.data[key]
	return value, ok
}

func (ctx *Context) Param(key string) string {
	return ctx.Params[key]
}
