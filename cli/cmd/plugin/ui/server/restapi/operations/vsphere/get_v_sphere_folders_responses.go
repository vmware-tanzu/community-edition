// Code generated by go-swagger; DO NOT EDIT.

package vsphere

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// GetVSphereFoldersOKCode is the HTTP code returned for type GetVSphereFoldersOK
const GetVSphereFoldersOKCode int = 200

/*GetVSphereFoldersOK Successful retrieval of vSphere folders

swagger:response getVSphereFoldersOK
*/
type GetVSphereFoldersOK struct {

	/*
	  In: Body
	*/
	Payload []*models.VSphereFolder `json:"body,omitempty"`
}

// NewGetVSphereFoldersOK creates GetVSphereFoldersOK with default headers values
func NewGetVSphereFoldersOK() *GetVSphereFoldersOK {

	return &GetVSphereFoldersOK{}
}

// WithPayload adds the payload to the get v sphere folders o k response
func (o *GetVSphereFoldersOK) WithPayload(payload []*models.VSphereFolder) *GetVSphereFoldersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v sphere folders o k response
func (o *GetVSphereFoldersOK) SetPayload(payload []*models.VSphereFolder) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVSphereFoldersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.VSphereFolder, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetVSphereFoldersBadRequestCode is the HTTP code returned for type GetVSphereFoldersBadRequest
const GetVSphereFoldersBadRequestCode int = 400

/*GetVSphereFoldersBadRequest Bad request

swagger:response getVSphereFoldersBadRequest
*/
type GetVSphereFoldersBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVSphereFoldersBadRequest creates GetVSphereFoldersBadRequest with default headers values
func NewGetVSphereFoldersBadRequest() *GetVSphereFoldersBadRequest {

	return &GetVSphereFoldersBadRequest{}
}

// WithPayload adds the payload to the get v sphere folders bad request response
func (o *GetVSphereFoldersBadRequest) WithPayload(payload *models.Error) *GetVSphereFoldersBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v sphere folders bad request response
func (o *GetVSphereFoldersBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVSphereFoldersBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetVSphereFoldersUnauthorizedCode is the HTTP code returned for type GetVSphereFoldersUnauthorized
const GetVSphereFoldersUnauthorizedCode int = 401

/*GetVSphereFoldersUnauthorized Incorrect credentials

swagger:response getVSphereFoldersUnauthorized
*/
type GetVSphereFoldersUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVSphereFoldersUnauthorized creates GetVSphereFoldersUnauthorized with default headers values
func NewGetVSphereFoldersUnauthorized() *GetVSphereFoldersUnauthorized {

	return &GetVSphereFoldersUnauthorized{}
}

// WithPayload adds the payload to the get v sphere folders unauthorized response
func (o *GetVSphereFoldersUnauthorized) WithPayload(payload *models.Error) *GetVSphereFoldersUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v sphere folders unauthorized response
func (o *GetVSphereFoldersUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVSphereFoldersUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetVSphereFoldersInternalServerErrorCode is the HTTP code returned for type GetVSphereFoldersInternalServerError
const GetVSphereFoldersInternalServerErrorCode int = 500

/*GetVSphereFoldersInternalServerError Internal server error

swagger:response getVSphereFoldersInternalServerError
*/
type GetVSphereFoldersInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVSphereFoldersInternalServerError creates GetVSphereFoldersInternalServerError with default headers values
func NewGetVSphereFoldersInternalServerError() *GetVSphereFoldersInternalServerError {

	return &GetVSphereFoldersInternalServerError{}
}

// WithPayload adds the payload to the get v sphere folders internal server error response
func (o *GetVSphereFoldersInternalServerError) WithPayload(payload *models.Error) *GetVSphereFoldersInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v sphere folders internal server error response
func (o *GetVSphereFoldersInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVSphereFoldersInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
