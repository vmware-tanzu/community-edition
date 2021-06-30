load("@ytt:data", "data")
load("@ytt:assert", "assert")

def validate_calico():
  data.values.infraProvider or assert.fail("Infrastructure provider should be provided")
end

#export
values = data.values

# validate
validate_calico()
