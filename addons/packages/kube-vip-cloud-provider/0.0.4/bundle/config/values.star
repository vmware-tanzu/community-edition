load("@ytt:data", "data")
load("@ytt:assert", "assert")

# export
values = data.values

# validate
if data.values.kubevipCloudProvider.loadbalancerCIDRs == "" and data.values.kubevipCloudProvider.loadbalancerIPRanges == "":
    assert.fail("Either loadbalancerCIDRs or loadbalancerIPRanges needs to be set. They can't be empty at the same time")
end
