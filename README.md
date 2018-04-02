# ECS Task Ports

This is a Go library for use by microservices running in an AWS Elastic Container Service environment that wish to know their port mappings to update service discovery systems.

For example a microservice might run on the port 8080 (the container port) but when run as a container using a Docker bridge network then the port is remapped to port 32481 (the host port).

This library supports two methods for determining the container's ports:
* Using the [container metadata file](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/container-metadata.html)
* Using [container agent introspection](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-agent-introspection.html).

The simplest and fastest solution (which also doesn't require any IAM permissions) is to use container metadata files. To enable this your EC2 instance must be running at least ECS container agent 1.15.0 or greater.

Either add `ECS_ENABLE_CONTAINER_METADATA=true` to the `/etc/ecs/ecs.config` file or ensure the container agent start with `ECS_ENABLE_CONTAINER_METADATA=true` in it's environment.

The second option using agent introspection takes a bit longer and requires three API calls - two to the ECS container agent, and another to the ECS API. Your task will also need the `ecs:DescribeTasks` IAM permission.

## API

```go
package main
import "github.com/alexbilbie/ecs-task-ports"

func main() {
    addresses, err := ecstaskports.Discover()
    if err != nil {
        // handle err...
    }

    for _, addr := range addresses {
        _ := os.Setenv(
            fmt.Sprintf("PORT_%d", addr.ContainerPort),
            fmt.Sprintf("%d", addr.HostPort),
        )
    }

    // serviceDiscovery.UpdateAdvertiseAddr(os.Getenv("PORT_8080"))
}
```