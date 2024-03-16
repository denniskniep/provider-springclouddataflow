# Spring Cloud Data Flow Provider

Spring Cloud Data Flow Provider is a [Crossplane](https://www.crossplane.io/) provider. It was build based on the [Crossplane Template](https://github.com/crossplane/provider-template). It is used to manage and configure [Spring Cloud Data Flow](https://dataflow.spring.io/). It uses the [Rest API](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources)

# How to use 
Repository and package:
```
xpkg.upbound.io/denniskniep/provider-springclouddataflow:<version>
```

Provider Credentials Structure:
```
{
  "url": "http://dataflow:9393/"
}
```

[View Example](./examples/provider/provider.yaml)

# Troubleshooting
Create a DeploymentRuntimeConfig and set the arg `--debug` on the package-runtime container

[View Example](./examples/provider/troubleshooting.yaml)

# Covered Managed Resources
Currently covered Managed Resources:
- [Application](#application)
- [Stream](#stream)
- [TaskDefinition](#taskdefinition)
- [TaskSchedule](#taskschedule)

## Application 

[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#applications) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#resources-registered-applications)

[View Example](./examples/application/application.yaml)

## Stream
[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#spring-cloud-dataflow-streams) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources-stream-definitions)

[View Example](./examples/stream/stream.yaml)


## TaskDefinition 

[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#spring-cloud-dataflow-task) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources-task-definitions)

[View Example](./examples/taskdefinition/taskdefinition.yaml)

## TaskSchedule 

[docs](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#_the_scheduler) 

[rest api](https://docs.spring.io/spring-cloud-dataflow/docs/current/reference/htmlsingle/#api-guide-resources-task-scheduler)

[View Example](./examples/taskschedule/taskschedule.yaml)

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


