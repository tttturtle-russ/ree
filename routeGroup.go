package ree

type Routers struct {
	prefix string
	parent *Routers
	engine *Engine
}

func (r *Routers) Group(prefix string) *Routers {
	return &Routers{
		prefix: r.prefix + prefix,
		parent: r,
		engine: r.engine,
	}
}

func (r *Routers) GET(path string, handler HandlerFunc) {
	r.engine.addRoute("GET", r.prefix+path, handler)
}
