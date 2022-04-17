package middleware

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ParamKey represents the key to retrieve the request parameters from the context
type ParamKey string

// Handler accepts a handler to make it compatible with http.HandlerFunc
// Example: r.GET("/", httprouterwrapper.Handler(http.HandlerFunc(controller.Index)))
func Handler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := newContextWithParams(r.Context(), p)
		h.ServeHTTP(w, r.WithContext(ctx)) //TODO: cancel somewhere?
	}
}

func newContextWithParams(ctx context.Context, p httprouter.Params) context.Context {
	var key ParamKey = "params"
	return context.WithValue(ctx, key, paramsToMap(p))
}

// ParamsToMap converts [httprouter.Params] type to a map of string->string
func paramsToMap(ps httprouter.Params) map[string]string {
	params := make(map[string]string, len(ps))
	for _, p := range ps {
		params[p.Key] = p.Value
	}
	return params
}
