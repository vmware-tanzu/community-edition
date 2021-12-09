load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("/globals.star", "globals")
load("@ytt:regexp", "regexp")

def validate_harbor_namespace():
  values.namespace or assert.fail("harbor namespace should be provided")
end

def validate_log_level():
  values.logLevel in ['debug', 'info', 'warning', 'error', 'fatal'] or assert.fail("logLevel must be debug, info, warning, error or fatal")
end

def validate_tls_certificate():
  ca_crt = getattr(values.tlsCertificate, "ca.crt")
  tls_crt = getattr(values.tlsCertificate, "tls.crt")
  tls_key = getattr(values.tlsCertificate, "tls.key")

  if not ca_crt and not tls_crt and not tls_key:
    return
  end

  values.tlsCertificateSecretName and assert.fail("tlsCertificateSecretName should be empty when tlsCertificate is provided")

  (tls_crt and tls_key) or assert.fail("Both the tls.crt and tls.key of tls certificate should be configured when tlsCertificate is provided")
end

def validate_harbor_admin_password():
  values.harborAdminPassword or assert.fail("Harbor admin password should be provided")
end

def validate_secret_key():
  values.secretKey or assert.fail("Secret key should be provided")
  len(values.secretKey) == 16 or assert.failed("Secret key must be a string of 16 chars")
end

def validate_core():
  values.core.secret or assert.fail("The core secret should be provided")

  values.core.xsrfKey or assert.fail("The core xsrf key should be provided")
  len(values.core.xsrfKey) == 32 or assert.failed("The core xsrf key must be a string of 32 chars")
end

def validate_jobservice():
  values.jobservice.secret or assert.fail("The jobservice secret should be provided")
end

def validate_registry():
  data.values.registry.secret or assert.fail("The registry secret should be provided")

  # if data.values.persistence.imageChartStorage.type == "filesystem":
  #   if data.values.persistence.persistentVolumeClaim.registry.accessMode == "ReadWriteOnce":
  #     data.values.registry.replicas == 1 or assert.fail("The registry replicas must be 1 when the image storage is filesystem and the access mode of persistentVolumeClaim is ReadWriteOnce")
  #   end
  # end
end

def validate_database():
  values.database.password or assert.fail("The database password should be provided")
end

def validate_image_chart_storage():
  storage = data.values.persistence.imageChartStorage
  storage.type in ("filesystem", "azure", "gcs", "s3", "swift", "oss") or assert.fail("Image and chart storage type should be filesystem|azure|gcs|s3|swift|oss")
  if storage.type == "azure":
    storage.azure.accountname or assert.fail("Azure account name should be provided")
    storage.azure.accountkey or assert.fail("Azure account key should be provided")
    storage.azure.container or assert.fail("Azure container should be provided")
  end
  if storage.type == "gcs":
    storage.gcs.bucket or assert.fail("GCS bucket should be provided")
  end
  if storage.type == "s3":
    storage.s3.region or assert.fail("S3 region should be provided")
    storage.s3.bucket or assert.fail("S3 bucket should be provided")
    storage.s3.storageclass in ("STANDARD", "REDUCED_REDUNDANCY") or assert.fail("S3 storage class should be STANDARD|REDUCED_REDUNDANCY")
  end
  if storage.type == "swift":
    storage.swift.authurl or assert.fail("Swift authurl should be provided")
    storage.swift.username or assert.fail("Swift username should be provided")
    storage.swift.password or assert.fail("Swift password should be provided")
    storage.swift.container or assert.fail("Swift container should be provided")
  end
  if storage.type == "oss":
    storage.oss.accesskeyid or assert.fail("OSS access key id should be provided")
    storage.oss.accesskeysecret or assert.fail("OSS access key secret should be provided")
  end
end

def validate_ip_families():
  ip_families = data.values.network.ipFamilies
  len(ip_families) > 0 or assert.fail("ipFamilies should be IPv4|IPv6, and it can not be empty")
  for ipType in ip_families:
    ipType in ("IPv4","IPv6") or assert.fail("ipFamilies should be IPv4|IPv6")
  end
end

def validate_httpproxy_timeout():
  pattern = '^(((\d*(\.\d*)?h)|(\d*(\.\d*)?m)|(\d*(\.\d*)?s)|(\d*(\.\d*)?ms)|(\d*(\.\d*)?us)|(\d*(\.\d*)?µs)|(\d*(\.\d*)?ns))+|infinity|infinite)$'
  if data.values.contourHttpProxy.timeout:
    regexp.match(pattern, data.values.contourHttpProxy.timeout) or assert.fail("The contourHttpProxy.timeout should be #h, or #m, or #s, or infinity or infinite, where # is a positive integer")
  end
  if data.values.contourHttpProxy.idleTimeout:
    regexp.match(pattern, data.values.contourHttpProxy.idleTimeout) or assert.fail("The contourHttpProxy.idleTimeout should be #h, or #m, or #s, or infinity or infinite, where # is a positive integer")
  end
end

def validate_trivy_timeout():
  pattern = '^(((\d*(\.\d*)?h)|(\d*(\.\d*)?m)|(\d*(\.\d*)?s)|(\d*(\.\d*)?ms)|(\d*(\.\d*)?us)|(\d*(\.\d*)?µs)|(\d*(\.\d*)?ns))+)$'
  if data.values.trivy.timeout:
    regexp.match(pattern, data.values.trivy.timeout) or assert.fail("The timeout set for trivy scanner should be #h, or #m, or #s, #ms, where # is a positive integer.")
  end
end

def validate_jaeger_config():
  if data.values.trace.enabled and data.values.trace.provider == "jaeger":
    if data.values.trace.jaeger.endpoint:
      data.values.trace.jaeger.endpoint != "http://hostname:14268/api/traces" or assert.fail("The endpoint must be configured if using collector mode with jaeger provider.")
      (not data.values.trace.jaeger.agent_host and not data.values.trace.jaeger.agent_port) or assert.fail("The agent mode should not be configured if collector mode has been set when using jaeger provider for tracing.")
    else:
      (data.values.trace.jaeger.agent_host and data.values.trace.jaeger.agent_port) or assert.fail("Neither collector mode was set nor agent mode was configured when using jaeger provider for tracing")
    end
  end
end

def validate_otel_config():
  if data.values.trace.enabled and data.values.trace.provider == "otel":
    (data.values.trace.otel.endpoint and data.values.trace.otel.endpoint != "hostname:4318") or assert.fail("The endpoint must be configured when using otel collector.")
    if data.values.trace.otel.timeout:
      data.values.trace.otel.timeout > 0 or assert.fail("It should be an integer representing the max waiting time for the backend to process each spans batch, in seconds.")
    end
  end
end

def validate_harbor():
  validate_funcs = [
    validate_harbor_namespace,
    validate_log_level,
    validate_tls_certificate,
    validate_harbor_admin_password,
    validate_secret_key,
    validate_core,
    validate_jobservice,
    validate_registry,
    validate_database,
    validate_image_chart_storage,
    validate_ip_families,
    validate_httpproxy_timeout,
    validate_trivy_timeout,
    validate_jaeger_config,
    validate_otel_config,
  ]
   for validate_func in validate_funcs:
     validate_func()
   end
end

#export
values = data.values

# validate harbor data values
validate_harbor()