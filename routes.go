package routes

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

// Routes represents the fasthttp aware routes
// and configuration that determine which route to choose
type Routes struct {
	Routes  map[string]fasthttp.RequestHandler
	Catcher fasthttp.RequestHandler
}

// New is a friendly, convenience function for returning
// an instance of routes.Routes that can be used in client code
func New() *Routes {
	return &Routes{
		Routes: make(map[string]fasthttp.RequestHandler),
		Catcher: func(ctx *fasthttp.RequestCtx) {
			ctx.SetStatusCode(404)
			fmt.Fprintf(ctx, "404 - no such route %s", string(ctx.Path()))
		},
	}
}

// Add takes a pattern, a function, and adds them to its self
// so requests can be routed correctly.
//
// A pattern can be a full url, or can use parameters.
// Params in URLs look like:
//   /users/:user/address
// This would match on:
//   /users/12345/address
// (For instance)
//
// Add() does no checking for existing routes; it is the responsibility
// of the developer to ensure there are no duplicates. The last function
// assigned to a patter will be used.
func (r *Routes) Add(pattern string, f fasthttp.RequestHandler) {
	r.Routes[normaliseRoute(pattern)] = f
}

// Route will send a fasthttp request to the correct function based
// on the path in the request.
// Parameters, as defined in a route, are accessed by ctx.userValue(param)
func (r Routes) Route(ctx *fasthttp.RequestCtx) {
	path := normaliseRoute(string(ctx.Path()))
	pathSplit := strings.Split(path, "/")

	for spec, f := range r.Routes {
		if spec == path {
			f(ctx)

			return
		}

		params := make(map[string]string)
		specSplit := strings.Split(spec, "/")

		if len(pathSplit) != len(specSplit) {
			continue
		}

		for idx, elem := range specSplit {
			pathElem := pathSplit[idx]

			if strings.HasPrefix(elem, ":") {
				params[elem] = pathElem
			} else if elem != pathElem {
				goto BADROUTE
			}
		}

		for k, v := range params {
			ctx.SetUserValue(stripTemplateChars(k), v)
		}

		f(ctx)

		return

	BADROUTE:
	}

	r.Catcher(ctx)

	return

}

func normaliseRoute(s string) string {
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	if strings.HasSuffix(s, "/") {
		return s
	}

	return s + "/"
}

func stripTemplateChars(s string) string {
	return strings.TrimPrefix(s, ":")
}
