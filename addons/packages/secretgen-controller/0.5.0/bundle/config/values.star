load("@ytt:data", "data")
load("@ytt:assert", "assert")

#export
values = data.values

values.secretgenController.namespace or assert.fail("namespace should be provided")


