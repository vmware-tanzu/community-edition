load("@ytt:data", "data")
load("@ytt:struct", "struct")

def _is_pkg_enabled(name):
    return (name not in data.values.excluded_packages)
end

exclusions = struct.make(
    is_pkg_enabled=_is_pkg_enabled,
)
