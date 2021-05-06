#!/bin/sh

# this shell script creates the prescribed directory structure for a tanzu package.

# set this value to your package name
NAME=$1

if [ -z "$NAME" ]
then
  # this name var comes a Makefile
  # kinda hacky, should figure out a better way
  echo "usage: hack/create-channel.sh foobar"
  exit 2
fi

CHANNEL_FILENAME=addons/repos/${NAME}.yaml

cat >> "$CHANNEL_FILENAME" <<EOL
#@data/values
---

package_repository:
  #! The name of the Package Repository.
  #! example: delta-foo.example.com
  name:

  #! The imgpkgBundle for the repo image.
  #! Note: this value is not known until a imgpkgBundle has been pushed to an OCI registry
  #! example: registry.example.com/foo/delta:v1
  imgpkgBundle:

  packages:
      #! The name of the package.
      #! example: foo
    - name:

      #! The domain that the package belongs to. This is used in conjunction with the name to create a fully qualified domain name for the package.
      #! example: example.com
      domain:

      #! The version of the package.
      #! example: 0.0.1
      version:

      #! The path to the image in the OCI repository.
      #! example: registry.example.com/foo/foo@sha256:abcd1234...
      image:

      #! A short description of the package.
      #! example: The foo package provides f, o and more o functionality.
      description:
EOL
echo
echo "CHANNEL file created at ${CHANNEL_FILENAME}"
echo
