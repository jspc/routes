package routes

import (
	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestRoutes_Route(t *testing.T) {
	r := New()
	r.Routes = map[string]fasthttp.RequestHandler{
		normaliseRoute("/"):                    func(ctx *fasthttp.RequestCtx) { fmt.Fprintf(ctx, "root") },
		normaliseRoute("/onething"):            func(ctx *fasthttp.RequestCtx) { fmt.Fprintf(ctx, "onething") },
		normaliseRoute("/a/couple/of/things"):  func(ctx *fasthttp.RequestCtx) { fmt.Fprintf(ctx, "a couple of things") },
		normaliseRoute("/a/:param/to/be/used"): func(ctx *fasthttp.RequestCtx) { fmt.Fprintf(ctx, "%+v", ctx.UserValue("param")) },
	}

	for _, test := range []struct {
		name   string
		url    string
		expect string
	}{
		{"Simple request, slash", "/", "root"},
		{"Simple request, a path", "/onething/", "onething"},
		{"Many path elems", "a/couple/of/things/", "a couple of things"},
		{"A templated value", "/a/thing/to/be/used/", "thing"},
		{"Undefined, simple path", "/nonesuch", "404 - no such route /nonesuch"},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := fasthttp.AcquireRequest()
			req.SetRequestURI(test.url)
			req.Header.SetMethod("GET")

			resp := fasthttp.AcquireResponse()

			c := &fasthttp.RequestCtx{
				Request:  *req,
				Response: *resp,
			}

			r.Route(c)

			received := string(c.Response.Body())
			if test.expect != received {
				t.Errorf("expected %q, received %q", test.expect, received)
			}
		})
	}
}

func TestRoutes_Add(t *testing.T) {
	for _, test := range []struct {
		name   string
		input  string
		expect string
	}{
		{"Empty string", "", "/"},
		{"Missing starting slash", "something/", "/something/"},
		{"Missing ending slash", "/something", "/something/"},
		{"Missing starting _and_ ending slashes", "something", "/something/"},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := New()
			r.Add(test.input, func(*fasthttp.RequestCtx) {})

			received := keys(r.Routes)[0]
			if test.expect != received {
				t.Errorf("expected %q, received %q", test.expect, received)
			}
		})
	}
}

func keys(m map[string]fasthttp.RequestHandler) (s []string) {
	s = make([]string, len(m))

	i := 0
	for k, _ := range m {
		s[i] = k
		i++
	}

	return
}
