// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

const (
	// DefaultConfigFile is config filename
	DefaultConfigFile string = "config.yaml"

	// TokenMinLength is 24
	TokenMinLength int = 24

	// DefaultLogLevel is 2
	DefaultLogLevel string = "2"

	// TODO(joshrosso): workaround so policy file does not need to be on user's machine.
	//                  get rid of this eventually
	StaticPolicy string = `{
    "default": [
        {
            "type": "insecureAcceptAnything"
        }
    ],
    "transports":
        {
            "docker-daemon":
                {
                    "": [{"type":"insecureAcceptAnything"}]
                }
        }
}`
)
