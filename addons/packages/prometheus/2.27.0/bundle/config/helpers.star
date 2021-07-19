load("@ytt:data", "data")

def generate_prometheus_tls():
  for key in ["tls.crt", "tls.key"]:
    if getattr(data.values.ingress.tlsCertificate, key):
      return False
    end
  end
  return True
end
