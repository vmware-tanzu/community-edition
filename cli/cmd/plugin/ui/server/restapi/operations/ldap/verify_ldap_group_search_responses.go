// Code generated by go-swagger; DO NOT EDIT.

package ldap

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// VerifyLdapGroupSearchOKCode is the HTTP code returned for type VerifyLdapGroupSearchOK
const VerifyLdapGroupSearchOKCode int = 200

/*VerifyLdapGroupSearchOK Verified LDAP credentials successfully

swagger:response verifyLdapGroupSearchOK
*/
type VerifyLdapGroupSearchOK struct {

	/*
	  In: Body
	*/
	Payload *models.LdapTestResult `json:"body,omitempty"`
}

// NewVerifyLdapGroupSearchOK creates VerifyLdapGroupSearchOK with default headers values
func NewVerifyLdapGroupSearchOK() *VerifyLdapGroupSearchOK {

	return &VerifyLdapGroupSearchOK{}
}

// WithPayload adds the payload to the verify ldap group search o k response
func (o *VerifyLdapGroupSearchOK) WithPayload(payload *models.LdapTestResult) *VerifyLdapGroupSearchOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the verify ldap group search o k response
func (o *VerifyLdapGroupSearchOK) SetPayload(payload *models.LdapTestResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *VerifyLdapGroupSearchOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// VerifyLdapGroupSearchBadRequestCode is the HTTP code returned for type VerifyLdapGroupSearchBadRequest
const VerifyLdapGroupSearchBadRequestCode int = 400

/*VerifyLdapGroupSearchBadRequest Bad request

swagger:response verifyLdapGroupSearchBadRequest
*/
type VerifyLdapGroupSearchBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewVerifyLdapGroupSearchBadRequest creates VerifyLdapGroupSearchBadRequest with default headers values
func NewVerifyLdapGroupSearchBadRequest() *VerifyLdapGroupSearchBadRequest {

	return &VerifyLdapGroupSearchBadRequest{}
}

// WithPayload adds the payload to the verify ldap group search bad request response
func (o *VerifyLdapGroupSearchBadRequest) WithPayload(payload *models.Error) *VerifyLdapGroupSearchBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the verify ldap group search bad request response
func (o *VerifyLdapGroupSearchBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *VerifyLdapGroupSearchBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// VerifyLdapGroupSearchUnauthorizedCode is the HTTP code returned for type VerifyLdapGroupSearchUnauthorized
const VerifyLdapGroupSearchUnauthorizedCode int = 401

/*VerifyLdapGroupSearchUnauthorized Incorrect credentials

swagger:response verifyLdapGroupSearchUnauthorized
*/
type VerifyLdapGroupSearchUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewVerifyLdapGroupSearchUnauthorized creates VerifyLdapGroupSearchUnauthorized with default headers values
func NewVerifyLdapGroupSearchUnauthorized() *VerifyLdapGroupSearchUnauthorized {

	return &VerifyLdapGroupSearchUnauthorized{}
}

// WithPayload adds the payload to the verify ldap group search unauthorized response
func (o *VerifyLdapGroupSearchUnauthorized) WithPayload(payload *models.Error) *VerifyLdapGroupSearchUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the verify ldap group search unauthorized response
func (o *VerifyLdapGroupSearchUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *VerifyLdapGroupSearchUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// VerifyLdapGroupSearchInternalServerErrorCode is the HTTP code returned for type VerifyLdapGroupSearchInternalServerError
const VerifyLdapGroupSearchInternalServerErrorCode int = 500

/*VerifyLdapGroupSearchInternalServerError Internal server error

swagger:response verifyLdapGroupSearchInternalServerError
*/
type VerifyLdapGroupSearchInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewVerifyLdapGroupSearchInternalServerError creates VerifyLdapGroupSearchInternalServerError with default headers values
func NewVerifyLdapGroupSearchInternalServerError() *VerifyLdapGroupSearchInternalServerError {

	return &VerifyLdapGroupSearchInternalServerError{}
}

// WithPayload adds the payload to the verify ldap group search internal server error response
func (o *VerifyLdapGroupSearchInternalServerError) WithPayload(payload *models.Error) *VerifyLdapGroupSearchInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the verify ldap group search internal server error response
func (o *VerifyLdapGroupSearchInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *VerifyLdapGroupSearchInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
