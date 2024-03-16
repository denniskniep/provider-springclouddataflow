# Spring Cloud Data Flow Provider

Spring Cloud Data Flow Provider is a [Crossplane](https://www.crossplane.io/) provider. It was build based on the [Crossplane Template](https://github.com/crossplane/provider-template). It is used to manage and configure [Spring Cloud Data Flow](https://dataflow.spring.io/). It uses the [Rest API](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources)

# How to use 
Repository and package:
```
xpkg.upbound.io/denniskniep/provider-springclouddataflow:<version>
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
  package: xpkg.upbound.io/denniskniep/provider-springclouddataflow:<version>
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
- [Stream](#stream)
- [TaskDefinition](#taskdefinition)
- [TaskSchedule](#taskschedule)

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
    version: "3.0.0"
    uri: "docker:springcloudtask/timestamp-task:3.0.0"
    bootVersion: "2"
    defaultVersion: true
  providerConfigRef:
    name: provider-spring-cloud-dataflow-config
```

## Stream
[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#spring-cloud-dataflow-streams) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources-stream-definitions)

Example:
```
apiVersion: core.springclouddataflow.crossplane.io/v1alpha1
kind: Stream
metadata:
  name: stream-1
spec:
  forProvider:
    name: "Stream01"
    description: "Test Stream"
    definition: "CHANGE | ME"
    deploy: false
  providerConfigRef:
    name: provider-spring-cloud-dataflow-config
```


## TaskDefinition 

[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#spring-cloud-dataflow-task) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources-task-definitions)

Example:
```
apiVersion: core.springclouddataflow.crossplane.io/v1alpha1
kind: TaskDefinition
metadata:
  name: task-1
spec:
  forProvider:
    name: "MyTask01"
    description: "Test Task"
    definition: "App001"
  providerConfigRef:
    name: provider-spring-cloud-dataflow-config
```

## TaskSchedule 

[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#_the_scheduler) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources-task-scheduler)

Example:
```
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
type: Opaque
stringData:
  credentialA: SecretA
  credentialB: SecretB
---
apiVersion: core.springclouddataflow.crossplane.io/v1alpha1
kind: TaskSchedule
metadata:
  name: schedule-1
spec:
  forProvider:
    scheduleName: "myschedule01"
    taskDefinitionNameRef: 
      name: "task-1"
    cronExpression: "* * * * *"
    platform: "default"
    arguments: "--myarg1=value1 --myarg2=value2"
    properties: "scheduler.kubernetes.jobAnnotations=annotation1:value1,annotation2:value2,scheduler.kubernetes.secretRefs=[my-secret]"
  providerConfigRef:
    name: provider-spring-cloud-dataflow-config
```

Reference for properties: https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#configuration-kubernetes-app-props


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


