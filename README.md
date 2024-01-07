# Spring Cloud Data Flow Provider

Spring Cloud Data Flow Provider is a [Crossplane](https://www.crossplane.io/) provider. It was build based on the [Crossplane Template](https://github.com/crossplane/provider-template). It is used to manage and configure [Spring Cloud Data Flow](https://dataflow.spring.io/). It uses the [Rest API](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources)

# How to use 
Repository and package:
```
xpkg.upbound.io/denniskniep/provider-springclouddataflow:v0.0.1
```

Provider Credentials:
```
{
  "url": "http://dataflow:9393/"
}
```

Example:
```
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-spring-cloud-dataflow
spec:
  package: xpkg.upbound.io/denniskniep/provider-springclouddataflow:v0.0.1
  packagePullPolicy: IfNotPresent
  revisionActivationPolicy: Automatic
---
apiVersion: v1
kind: Secret
metadata:
  name: provider-spring-cloud-dataflow-config-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "url": "http://dataflow:9393/"
    }
---
apiVersion: springclouddataflow.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: provider-spring-cloud-dataflow-config
spec: 
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: provider-spring-cloud-dataflow-config-creds
      key: credentials  
```
# Troubleshooting
Create a DeploymentRuntimeConfig and set the arg `--debug` on the package-runtime container:

```
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: debug-config
spec:
  deploymentTemplate:
    spec:
      selector: {}
      template:
        spec:
          containers:
            - name: package-runtime
              args:
                - --debug
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-spring-cloud-dataflow
spec:
  package: xpkg.upbound.io/denniskniep/provider-springclouddataflow:v0.0.1
  packagePullPolicy: IfNotPresent
  revisionActivationPolicy: Automatic
  runtimeConfigRef:
    name: debug-config
```

# Covered Managed Resources
Currently covered Managed Resources:
- [Application](#application)

## Application 

[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#applications) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#resources-registered-applications)

Example:
```
apiVersion: core.springclouddataflow.crossplane.io/v1alpha1
kind: Application
metadata:
  name: app-1
spec:
  forProvider:
    name: "App001"
    type: "task"
    version: "v1.0.0"
    uri: "docker://hello-world:v1.0.0"
    bootVersion: "2"
    defaultVersion: true
  providerConfigRef:
    name: provider-spring-cloud-dataflow-config
```

# Contribute
## Developing
1. Add new type by running the following command:
```shell
  export provider_name=SpringCloudDataFlow
  export group=core # lower case e.g. core, cache, database, storage, etc.
  export type=MyType # Camel casee.g. Bucket, Database, CacheCluster, etc.
  make provider.addtype provider=${provider_name} group=${group} kind=${type}
```
2. Replace the *core* group with your new group in apis/{provider}.go
3. Replace the *MyType* type with your new type in internal/controller/{provider}.go

4. Run `make reviewable` to run code generation, linters, and tests. (`make generate` to only run code generation)
5. Run `make build` to build the provider.

Refer to Crossplane's [CONTRIBUTING.md] file for more information on how the
Crossplane community prefers to work. The [Provider Development][provider-dev]
guide may also be of use.

[CONTRIBUTING.md]: https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md
[provider-dev]: https://github.com/crossplane/crossplane/blob/master/contributing/guide-provider-development.md

## Tests
Start SpringCloudDataFlow environment for tests
```
sudo docker-compose -f tests/docker-compose.yaml up 
```
UI: http://localhost:9393/dashboard
OpenAPI Spec: http://localhost:9393/v3/api-docs
Swagger-Ui: http://localhost:9393/swagger-ui/index.html


