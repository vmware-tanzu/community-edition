load("@ytt:data", "data")
load("@ytt:assert", "assert")

def validate_antrea():
  data.values.infraProvider or assert.fail("Infrastructure provider should be provided")
end

# export data.values
values = data.values

# validate antrea configuration
validate_antrea()
