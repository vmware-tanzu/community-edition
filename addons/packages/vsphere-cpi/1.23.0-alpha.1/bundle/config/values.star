load("@ytt:data", "data")
load("@ytt:assert", "assert")

def validate_vsphereCPI():
   data.values.vsphereCPI.server or assert.fail("vsphereCPI server should be provided")
   data.values.vsphereCPI.datacenter or assert.fail("vsphereCPI datacenter should be provided")
   data.values.vsphereCPI.username or assert.fail("vsphereCPI username should be provided")
   data.values.vsphereCPI.password or assert.fail("vsphereCPI password should be provided")
   if not data.values.vsphereCPI.insecureFlag:
     data.values.vsphereCPI.tlsThumbprint or assert.fail("vsphereCPI tlsThumbprint should be provided when insecureFlag is False")
   end
   if data.values.vsphereCPI.ipFamily and (data.values.vsphereCPI.ipFamily not in ["ipv4", "ipv6", "ipv4,ipv6", "ipv6,ipv4"]):
     assert.fail("vsphereCPI ipFamily should be one of \"ipv4\", \"ipv6\", \"ipv4,ipv6\", or \"ipv6,ipv4\" if provided")
   end
end

def validate_vsphereParavirtualCPI():
   data.values.vsphereCPI.clusterName or assert.fail("vsphereParavirtualCPI name of cluster should be provided")
   data.values.vsphereCPI.clusterUID or assert.fail("vsphereParavirtualCPI UID of cluster should be provided")
   data.values.vsphereCPI.supervisorMasterEndpointIP or assert.fail("vsphereParavirtualCPI supervisorMasterEndpointIP of cluster should be provided")
   data.values.vsphereCPI.supervisorMasterPort or assert.fail("vsphereParavirtualCPI supervisorMasterPort of cluster should be provided")
end

def validate_nsxt_config():
   if validate_nsxt_username_password() == False and validate_nsxt_secret() == False and validate_nsxt_token() == False and validate_nsxt_cert() == False:
     assert.fail("Invalid NSX-T credentials: username/password or vmc access token or client certificates must be set")
   end
   data.values.vsphereCPI.nsxt.host or assert.fail("vsphereCPI nsxtHost should be provided")
   data.values.vsphereCPI.nsxt.routes.clusterCidr or assert.fail("vsphereCPI nsxt routes clusterCidr should be provided")
end

def validate_nsxt_token():
   if data.values.vsphereCPI.nsxt.vmcAccessToken == "" or data.values.vsphereCPI.nsxt.vmcAuthHost == "":
     return False
   end
   return True
end

def validate_nsxt_cert():
   if data.values.vsphereCPI.nsxt.clientCertKeyData == "" or data.values.vsphereCPI.nsxt.clientCertData == "":
     return False
   end
   return True
end

def validate_nsxt_secret():
   if data.values.vsphereCPI.nsxt.secretName == "" or data.values.vsphereCPI.nsxt.secretNamespace == "" or validate_nsxt_username_password() == False:
     return False
   end
   return True
end

def validate_nsxt_username_password():
   if data.values.vsphereCPI.nsxt.username == "" or data.values.vsphereCPI.nsxt.password == "":
     return False
   end
   return True
end

# export
values = data.values

# validate

def validate():
    if data.values.vsphereCPI.mode == "vsphereCPI" or not data.values.vsphereCPI.mode:
        validate_vsphereCPI()
    elif data.values.vsphereCPI.mode == "vsphereParavirtualCPI":
        validate_vsphereParavirtualCPI()
    else:
        assert.fail("vsphereCPI mode should be either vsphereCPI or vsphereParavirtualCPI")
    end
end

validate()
if data.values.vsphereCPI.nsxt.podRoutingEnabled:
validate_nsxt_config()
end
