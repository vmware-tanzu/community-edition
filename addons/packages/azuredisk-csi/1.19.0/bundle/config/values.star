load("@ytt:data", "data")
load("@ytt:assert", "assert")

def validate_azureCSI():
   data.values.azureCSI.namespace or assert.fail("azureCSI namespace should be provided")
end

#export
values = data.values

# validate
validate_azureCSI()
