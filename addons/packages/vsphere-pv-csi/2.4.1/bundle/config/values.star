load("@ytt:data", "data")
load("@ytt:assert", "assert")

def validate_vspherePVCSI():
   data.values.vspherePVCSI.namespace or assert.fail("vspherePVCSI namespace should be provided")
   data.values.vspherePVCSI.supervisor_master_endpoint_hostname or assert.fail("vspherePVCSI supervisor_master_endpoint_hostname should be provided")
   data.values.vspherePVCSI.supervisor_master_port or assert.fail("vspherePVCSI supervisor_master_port should be provided")
   data.values.vspherePVCSI.tanzukubernetescluster_uid or assert.fail("vspherePVCSI tanzukubernetescluster_uid should be provided")
   data.values.vspherePVCSI.tanzukubernetescluster_name or assert.fail("vspherePVCSI tanzukubernetescluster_name should be provided")
end

#export
values = data.values

# validate
validate_vspherePVCSI()
