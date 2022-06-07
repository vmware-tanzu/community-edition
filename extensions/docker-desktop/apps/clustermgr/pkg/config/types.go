// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

// GetDefaultPorts returns the default ports exposed for a cluster.
func GetDefaultPorts() []string {
	return []string{"80", "443"}
}

// RunningResponse is a helper to return a command response for a running cluster.
func RunningResponse() Response {
	return Response{
		Status:      Running,
		Description: "Cluster is running",
	}
}

// DeletingResponse is a helper to return a command response for a cluster being deleted.
func DeletingResponse() Response {
	return Response{
		Status:      Deleting,
		Description: "Cluster is deleting",
	}
}

// DeletedResponse is a helper to return a command response for a deleted cluster.
func DeletedResponse() Response {
	return Response{
		Status:      Deleted,
		Description: "Cluster is deleted",
	}
}

// NotExistsResponse is a helper to return a command response for a cluster that does not exist.
func NotExistsResponse() Response {
	return Response{
		Status:      NotExists,
		Description: "Cluster does not exist",
	}
}

// EmptyStatsResponse is a helper to return a command response that is meant to be updated with
// current information.
func EmptyStatsResponse() Response {
	return Response{
		Status:      NotExists,
		Description: "",
		Error:       false,
	}
}

// Status is the cluster status string.
type Status string

const (
	// Unknown indicates the cluster's status is unknown.
	Unknown Status = "Unknown"
	// NotExists indicates the cluster does not exist.
	NotExists Status = "NotExists"
	// Initializing indicates the cluster is being initialized.
	Initializing Status = "Initializing"
	// Creating indicates the cluster is being created.
	Creating Status = "Creating"
	// Running indicates the cluster is running.
	Running Status = "Running"
	// Stopped indicates the cluster is stopped.
	Stopped Status = "Stopped"
	// Deleting indicates the cluster is being deleted.
	Deleting Status = "Deleting"
	// Deleted indicates the cluster is deleted.
	Deleted Status = "Deleted"
	// Error indicates there was an error getting the cluster's status.
	Error Status = "Error"
)

// Response contains the results of a cluster operation.
type Response struct {
	// Status is the status of the cluster or operation.
	Status Status `json:"status,omitempty"`
	// Description contains the text description of the operation.
	Description string `json:"description,omitempty"`
	// Error indicates if there was an error performing the operation.
	Error bool `json:"isError,omitempty"`
	// ErrorMessage is the text description of an error.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Output is the output text of the operation.
	Output string `json:"output,omitempty"`
}
