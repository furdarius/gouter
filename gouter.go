package gouter

import (
	"net/http"
)

// Params is map of HTTP request params, where key is param name.
type Params map[string]string

// HandlerFunc represent a function for HTTP request handling.
// Similar with http.HandlerFunc, but third parameter is 'Params'
type HandlerFunc func(http.ResponseWriter, *http.Request, Params)

const (
	methodsNum uint8 = 5
)

// Router is an implementation of http.Handler. It is used for processing of incoming url requests.
type Router struct {
	methods         map[string]*node
	notFoundHandler http.HandlerFunc
}

// New returns a new instance of Router.
func New() *Router {
	return &Router{
		methods:         make(map[string]*node, methodsNum),
		notFoundHandler: nil,
	}
}

// Get is alias of Add("GET", route, handle)
func (r *Router) Get(route string, handle HandlerFunc) {
	r.Add("GET", route, handle)
}

// Post is alias of Add("POST", route, handle)
func (r *Router) Post(route string, handle HandlerFunc) {
	r.Add("POST", route, handle)
}

// Put is alias of Add("PUT", route, handle)
func (r *Router) Put(route string, handle HandlerFunc) {
	r.Add("PUT", route, handle)
}

// Patch is alias of Add("PATCH", route, handle)
func (r *Router) Patch(route string, handle HandlerFunc) {
	r.Add("PATCH", route, handle)
}

// Delete is alias of Add("DELETE", route, handle)
func (r *Router) Delete(route string, handle HandlerFunc) {
	r.Add("DELETE", route, handle)
}

// Add register handler for defined route
func (r *Router) Add(method string, route string, handle HandlerFunc) {
	if len(route) < 1 {
		panic("Empty route is not allowed!")
	}

	if route[0] != '/' {
		panic("Route must start with '/' symbol, got \"" + route + "\"")
	}

	root := r.methods[method]

	if root == nil {
		root = newRootNode()

		r.methods[method] = root
	}

	root.insert(route, handle)
}

// SetNotFoundHandler alow user define own handler for 404 error.
func (r *Router) SetNotFoundHandler(handler http.HandlerFunc) {
	r.notFoundHandler = handler
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	route := request.URL.Path
	method := request.Method

	root := r.methods[method]

	if root == nil {
		if r.notFoundHandler == nil {
			http.NotFound(responseWriter, request)

			return
		}

		r.notFoundHandler(responseWriter, request)

		return
	}

	handler, params := root.lookup(route)

	if handler == nil {
		if r.notFoundHandler == nil {
			http.NotFound(responseWriter, request)

			return
		}

		r.notFoundHandler(responseWriter, request)

		return
	}

	handler(responseWriter, request, params)
}
