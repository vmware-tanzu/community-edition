load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("/globals.star", "globals")

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

  (tls_crt and tls_key) or assert.fail("Both the tls.crt and tls.key of tls certificate should be provided")
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
  ]
   for validate_func in validate_funcs:
     validate_func()
   end
end

#export
values = data.values

# validate harbor data values
validate_harbor()