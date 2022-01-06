// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// package log provides logging mechanisms
// This file defines animator options to be used with AnimateProgressWithOptions
package log

import "context"

type progressAnimatorOptions struct {
	messagef     string
	messagefArgs []string

	maxLen   int
	ctx      context.Context
	statChan chan string
}

// AnimatorOption is an option to be passed to the AnimateProgressWithOptions logging method
type AnimatorOption interface {
	apply(*progressAnimatorOptions)
}

type animatorOptionAdapter func(*progressAnimatorOptions)

func (a animatorOptionAdapter) apply(o *progressAnimatorOptions) {
	a(o)
}

// AnimatorWithMessagef sets the format string message and any format arguments to use
// Ex: AnimatorWithMessagef("Downloading to: %s", filePathLocation)
//
// If no format arguments are provided, just the message will be displayed when animating
// Ex: AnimatorWithMessagef("Creating controller")
//
// If a status channel is used, the first format argument in the template string
// _must_ be the status associated with the status channel
// Ex: AnimatorWithMessagef("controller status: %s")
func AnimatorWithMessagef(formatString string, formatArgs ...string) AnimatorOption {
	return animatorOptionAdapter(func(o *progressAnimatorOptions) {
		o.messagef = formatString
		o.messagefArgs = formatArgs
	})
}

// AnimatorWithMaxLen sets the maximum number of dots to animate
// Ex:AnimatorWithMaxLen(3)
func AnimatorWithMaxLen(l int) AnimatorOption {
	return animatorOptionAdapter(func(o *progressAnimatorOptions) {
		o.maxLen = l
	})
}

// AnimatorWithContext provides a context to the async call.
// That context can be canceled which stops the animation.
// Ex: AnimatorWithContext(myContext)
func AnimatorWithContext(ctx context.Context) AnimatorOption {
	return animatorOptionAdapter(func(o *progressAnimatorOptions) {
		o.ctx = ctx
	})
}

// AnimatorWithStatusChan sets a string status channel that will
// be used to asynchronously inspect the status of an operation
// Ex: AnimatorWithStatusChan(myStatusChan)
func AnimatorWithStatusChan(s chan string) AnimatorOption {
	return animatorOptionAdapter(func(o *progressAnimatorOptions) {
		o.statChan = s
	})
}
