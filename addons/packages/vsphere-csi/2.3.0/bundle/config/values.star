load("@ytt:data", "data")
load("@ytt:assert", "assert")

def validate_vsphereCSI():
   data.values.vsphereCSI.namespace != "" or assert.fail("vsphereCSI namespace should be provided")
   data.values.vsphereCSI.clusterName != "" or assert.fail("vsphereCSI clusterName should be provided")
   data.values.vsphereCSI.server != "" or assert.fail("vsphereCSI server should be provided")
   data.values.vsphereCSI.datacenter != "" or assert.fail("vsphereCSI datacenter should be provided")
   data.values.vsphereCSI.publicNetwork != "" or assert.fail("vsphereCSI publicNetwork should be provided")
   data.values.vsphereCSI.username != "" or assert.fail("vsphereCSI username should be provided")
   data.values.vsphereCSI.password != "" or assert.fail("vsphereCSI password should be provided")
   if not data.values.vsphereCSI.insecureFlag:
     data.values.vsphereCSI.tlsThumbprint != "" or assert.fail("vsphereCSI tlsThumbprint should be provided when insecureFlag is False")
   end
end

#export
values = data.values

# validate
validate_vsphereCSI()
