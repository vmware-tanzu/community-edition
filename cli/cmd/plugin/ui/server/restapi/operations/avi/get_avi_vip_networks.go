// Code generated by go-swagger; DO NOT EDIT.

package avi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetAviVipNetworksHandlerFunc turns a function with the right signature into a get avi vip networks handler
type GetAviVipNetworksHandlerFunc func(GetAviVipNetworksParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAviVipNetworksHandlerFunc) Handle(params GetAviVipNetworksParams) middleware.Responder {
	return fn(params)
}

// GetAviVipNetworksHandler interface for that can handle valid get avi vip networks params
type GetAviVipNetworksHandler interface {
	Handle(GetAviVipNetworksParams) middleware.Responder
}

// NewGetAviVipNetworks creates a new http.Handler for the get avi vip networks operation
func NewGetAviVipNetworks(ctx *middleware.Context, handler GetAviVipNetworksHandler) *GetAviVipNetworks {
	return &GetAviVipNetworks{Context: ctx, Handler: handler}
}

/* GetAviVipNetworks swagger:route GET /api/avi/vipnetworks avi getAviVipNetworks

Retrieve all Avi networks

*/
type GetAviVipNetworks struct {
	Context *middleware.Context
	Handler GetAviVipNetworksHandler
}

func (o *GetAviVipNetworks) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetAviVipNetworksParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
