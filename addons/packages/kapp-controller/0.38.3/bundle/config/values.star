load("@ytt:data", "data")

#export
values = data.values
kappNamespace = ""
if hasattr(values.kappController, 'namespace') and values.kappController.namespace:
    kappNamespace = values.kappController.namespace
else:
    kappNamespace = values.namespace
end

def generateBashCmdForDNS(coreDNSIP):
    # This command added the coreDNS IP as the first entry of resolv.conf
    # In this way, Kapp Controller will have cluster IP access,
    # and still able to resolve enternal urls when core DNS is unavailable

    return "cp /etc/resolv.conf /etc/resolv.conf.bak; sed '1 i nameserver " + coreDNSIP + "' /etc/resolv.conf.bak > /etc/resolv.conf; rm /etc/resolv.conf.bak; cp -R /etc/* /kapp-etc; chmod g+w /kapp-etc/pki/tls/certs/"
end
