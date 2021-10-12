load("@ytt:data", "data")
load("@ytt:assert", "assert")

#export
values = data.values
secretgenNamespace = ""
if hasattr(values.secretgenController, 'namespace') and values.secretgenController.namespace:
    secretgenNamespace = values.secretgenController.namespace
else:
    secretgenNamespace = values.namespace
end
