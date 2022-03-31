// Code generated by go-swagger; DO NOT EDIT.

package provider

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetProviderHandlerFunc turns a function with the right signature into a get provider handler
type GetProviderHandlerFunc func(GetProviderParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetProviderHandlerFunc) Handle(params GetProviderParams) middleware.Responder {
	return fn(params)
}

// GetProviderHandler interface for that can handle valid get provider params
type GetProviderHandler interface {
	Handle(GetProviderParams) middleware.Responder
}

// NewGetProvider creates a new http.Handler for the get provider operation
func NewGetProvider(ctx *middleware.Context, handler GetProviderHandler) *GetProvider {
	return &GetProvider{Context: ctx, Handler: handler}
}

/* GetProvider swagger:route GET /api/providers provider getProvider

Get infrastructure provider given by the user via cli

*/
type GetProvider struct {
	Context *middleware.Context
	Handler GetProviderHandler
}

func (o *GetProvider) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetProviderParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
