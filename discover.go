package ecstaskports

import "os"

// PortMapping struct contains the container port (which is the original EXPOSEd port in a container) and the host port which is the port mapped to the container port
type PortMapping struct {
	ContainerPort int
	HostPort      int
}

// Discover will attempt to return an array ofPortMapping, or an error
func Discover() ([]PortMapping, error) {
	ecsMetadataFile := os.Getenv("ECS_CONTAINER_METADATA_FILE")
	if ecsMetadataFile != "" {
		return readInstanceMetadataFile(ecsMetadataFile)
	}

	return discoverWithTaskIntrospection()
}
