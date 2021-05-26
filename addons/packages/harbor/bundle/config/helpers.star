load("@ytt:data", "data")

def generate_harbor_tls():
  for key in ["tls.crt", "tls.key"]:
    if getattr(data.values.tlsCertificate, key):
      return False
    end
  end
  return True
end

def get_storage_class(sc):
    return "" if sc == "-" else sc
end

def use_existing_claim(config):
    def _use_existing_claim(x, y, z):
        return bool(config.existingClaim)
    end
    return _use_existing_claim
end

def get_no_proxy():
    components = [
        "core",
        "jobservice",
        "database",
        "chartmuseum",
        "clair",
        "notary-server",
        "notary-signer",
        "registry",
        "portal",
        "trivy",
    ]
    items = []
    for component in components:
        items.append("harbor-{0}".format(component))
    end

    items.extend("127.0.0.1,localhost,.local,.internal".split(","))

    for item in data.values.proxy.noProxy.split(","):
        if item not in items:
            items.append(item)
        end
    end

    return ",".join(items)
end

def get_external_url():
  if data.values.port.https == 443:
    return "https://{}".format(data.values.hostname)
  else:
    return "https://{0}:{1}".format(data.values.hostname, data.values.port.https)
  end
end

def get_db_url(db):
  return "postgres://postgres:{0}@harbor-database:5432/{1}?sslmode=disable".format(data.values.database.password, db)
end

def get_notary_hostname():
  return "notary.{}".format(data.values.hostname)
end