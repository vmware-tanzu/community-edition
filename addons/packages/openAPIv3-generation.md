# How to generate openAPIv3 schema and use it in a package

Follow the below mentioned steps to get started on generating openAPIv3 schema and specifying it in a package:

1. Create a schema file for given data values file. In ytt, before a Data Value can be used in a template, it must be declared. This is typically done via Data Values Schema
   Check out [How to write Schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/), to explore the different annotations that can be used when writing a schema.

2. Run this command to use a data value and schema file together. You can see if you change a value in the `values.yml` to be the incorrect type, you will see an error catches this.

   Use the following template.yaml file

    ```yaml
    #@ load("@ytt:data", "data")
    ---
    data_values: #@ data.values
    ```

   Let's use secretgen-controller package as an example to generate schema and embed it in Package

   ```bash
   cd ~/community-edition/addons/packages/secretgen-controller/0.7.1/bundle/config
   ```

   You can now generate and verify data values by using the command below

   `ytt -f values.yaml -f schema.yaml -f template.yaml`

3. Generate the OpenAPI schema from the `schema.yaml` file.

   `ytt -f schema.yaml --data-values-schema-inspect -o openapi-v3 > openapi-schema.yaml`

4. Create package-template.yaml file from package.yaml by adding the following 2 lines in `spec` field:

   ```yaml
   valuesSchema:
       openAPIv3:  #@ yaml.decode(data.values.openapi)["components"]["schemas"]["dataValues"]
   ```

   For example, for generating openAPIv3 schema for secretgen-controller, the following package-template.yaml file is used:

   ```yaml
   #@ load("@ytt:data", "data")
   #@ load("@ytt:yaml", "yaml")
   ---
   apiVersion: data.packaging.carvel.dev/v1alpha1
   kind: Package
   metadata:
     name: secretgen-controller.community.tanzu.vmware.com.0.7.1
   spec:
     valuesSchema:
       openAPIv3:  #@ yaml.decode(data.values.openapi)["components"]["schemas"]["dataValues"]
     refName: secretgen-controller.community.tanzu.vmware.com
     version: 0.7.1
     releaseNotes: "secretgen-controller 0.7.1 https://github.com/vmware-tanzu/carvel-secretgen-controller"
     licenses:
       - "Apache 2.0"
     template:
       spec:
         fetch:
           - imgpkgBundle:
               image: projects.registry.vmware.com/tce/secretgen-controller@sha256:4248e36490eb888d7f8bd0b62739a5acc3f178f67d8c2abfb3a6181b814c074e
         template:
           - ytt:
               paths:
                 - config/
           - kbld:
               paths:
                 - "-"
                 - .imgpkg/images.yml
         deploy:
           - kapp: {}
   ```

5. You can use this OpenAPI schema as documentation, or use it in a Package, to document what inputs are allowed in the package.

   `ytt -f ../../package-template.yaml --data-value-file openapi=openapi-schema.yml > ../../package.yaml`

6. You can now see openAPIv3 schema specified in the Package for secretgen-controller

   ```yaml
   apiVersion: data.packaging.carvel.dev/v1alpha1
   kind: Package
   metadata:
     name: secretgen-controller.community.tanzu.vmware.com.0.7.1
   spec:
     valuesSchema:
       openAPIv3:
         type: object
         additionalProperties: false
         description: OpenAPIv3 Schema for secret gen controller
         properties:
           secretgenController:
             type: object
             additionalProperties: false
             description: Configuration for secret gen controller
             properties:
               namespace:
                 type: string
                 default: secretgen-controller
                 description: Namespace for secret gen controller
               createNamespace:
                 type: boolean
                 default: false
                 description: Whether to create namespace for secret gen controller if not present
     refName: secretgen-controller.community.tanzu.vmware.com
     version: 0.7.1
     releaseNotes: secretgen-controller 0.7.1 https://github.com/vmware-tanzu/carvel-secretgen-controller
     licenses:
     - Apache 2.0
     template:
       spec:
         fetch:
         - imgpkgBundle:
             image: projects.registry.vmware.com/tce/secretgen-controller@sha256:5a77400c7c56f9a439d47e19086256587607e988cf0529702f87d5005724d472
         template:
         - ytt:
             paths:
             - config/
         - kbld:
             paths:
             - '-'
             - .imgpkg/images.yml
         deploy:
         - kapp: {}
   ```
