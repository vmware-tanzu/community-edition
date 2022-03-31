// Code generated by go-swagger; DO NOT EDIT.

package avi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// VerifyAccountHandlerFunc turns a function with the right signature into a verify account handler
type VerifyAccountHandlerFunc func(VerifyAccountParams) middleware.Responder

// Handle executing the request and returning a response
func (fn VerifyAccountHandlerFunc) Handle(params VerifyAccountParams) middleware.Responder {
	return fn(params)
}

// VerifyAccountHandler interface for that can handle valid verify account params
type VerifyAccountHandler interface {
	Handle(VerifyAccountParams) middleware.Responder
}

// NewVerifyAccount creates a new http.Handler for the verify account operation
func NewVerifyAccount(ctx *middleware.Context, handler VerifyAccountHandler) *VerifyAccount {
	return &VerifyAccount{Context: ctx, Handler: handler}
}

/* VerifyAccount swagger:route POST /api/avi avi verifyAccount

Validate Avi controller credentials

*/
type VerifyAccount struct {
	Context *middleware.Context
	Handler VerifyAccountHandler
}

func (o *VerifyAccount) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewVerifyAccountParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
