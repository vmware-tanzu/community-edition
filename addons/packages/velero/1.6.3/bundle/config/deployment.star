load("@ytt:data", "data")
load("@ytt:assert", "assert")


def get_image_location():
    if data.values.image.digest:
        return '{0}/{1}@{2}'.format(data.values.image.repository, data.values.image.name, data.values.image.digest)
    end
    return '{0}/{1}:{2}'.format(data.values.image.repository, data.values.image.name, data.values.image.tag)
end

# export
deployment = data.values
